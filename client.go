package rdirsync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hnakamur/rdirsync/pb"
	"google.golang.org/grpc"
)

type Client struct {
	client            pb.RDirSyncClient
	bufSize           int
	atMostCount       int
	keepDeletedFiles  bool
	syncModTime       bool
	updateOnly        bool
	syncOwnerAndGroup bool
	userGroupDB       *userGroupDB
}

type ClientOption func(*Client) error

func NewClient(cc *grpc.ClientConn, option ...ClientOption) (*Client, error) {
	c := &Client{
		client:      pb.NewRDirSyncClient(cc),
		bufSize:     64 * 1024,
		atMostCount: 1024,
	}

	for _, opt := range option {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func SetBufSize(bufSize int) ClientOption {
	return func(c *Client) error {
		if bufSize < 0 {
			return errors.New("buffer size must be positive")
		}
		c.bufSize = bufSize
		return nil
	}
}

func SetMaxEntriesPerReadDirRPC(maxEntriesPerRPC int) ClientOption {
	return func(c *Client) error {
		if maxEntriesPerRPC < 0 {
			return errors.New("max entries per RPC must be positive")
		}
		c.atMostCount = maxEntriesPerRPC
		return nil
	}
}

func SetKeepDeletedFiles(keepDeletedFiles bool) ClientOption {
	return func(c *Client) error {
		c.keepDeletedFiles = keepDeletedFiles
		return nil
	}
}

func SetSyncModTime(syncModTime bool) ClientOption {
	return func(c *Client) error {
		c.syncModTime = syncModTime
		return nil
	}
}

func SetUpdateOnly(updateOnly bool) ClientOption {
	return func(c *Client) error {
		c.updateOnly = updateOnly
		return nil
	}
}

func SetSyncOwnerAndGroup(syncOwnerAndGroup bool) ClientOption {
	return func(c *Client) error {
		if syncOwnerAndGroup {
			if os.Getuid() != 0 {
				return errors.New("must be run by root user to sync owner and group")
			}
			c.userGroupDB = newUserGroupDB()
		}
		c.syncOwnerAndGroup = syncOwnerAndGroup
		return nil
	}
}

func (c *Client) statRemote(ctx context.Context, remotePath string) (*fileInfo, error) {
	info, err := c.client.Stat(ctx, &pb.StatRequest{
		Path:               remotePath,
		WantsOwnerAndGroup: c.syncOwnerAndGroup,
	})
	if err != nil {
		return nil, err
	}
	return c.newFileInfoFromRPC(info)
}

func (c *Client) statLocal(localPath string) (*fileInfo, error) {
	fi, err := os.Stat(localPath)
	if err != nil {
		return nil, err
	}

	return c.newFileInfoFromOS(fi)
}

func (c *Client) Fetch(ctx context.Context, remotePath, localPath string) error {
	rfi, err := c.statRemote(ctx, remotePath)
	if err != nil {
		return err
	}

	if rfi == nil {
		return fmt.Errorf("remote file or directory %q not found", remotePath)
	}

	if rfi.IsDir() {
		return c.fetchDirAndChmod(ctx, remotePath, localPath, rfi)
	} else {
		lfi, err := c.statLocal(localPath)
		if os.IsNotExist(err) {
			lfi = nil
		} else if err != nil {
			return err
		}
		return c.fetchFileAndChmod(ctx, remotePath, localPath, rfi, lfi)
	}
}

func (c *Client) FetchFile(ctx context.Context, remotePath, localPath string) error {
	rfi, err := c.statRemote(ctx, remotePath)
	if err != nil {
		return err
	}

	if rfi.IsDir() {
		return fmt.Errorf("expected a remote file but is a directory %q", remotePath)
	}

	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(err) {
		lfi = nil
	} else if err != nil {
		return err
	}
	return c.fetchFileAndChmod(ctx, remotePath, localPath, rfi, lfi)
}

func (c *Client) fetchFileAndChmod(ctx context.Context, remotePath, localPath string, rfi, lfi *fileInfo) error {
	if !c.isUpdateNeeded(rfi, lfi) {
		return nil
	}

	stream, err := c.client.FetchFile(ctx, &pb.FetchFileRequest{
		Path:    remotePath,
		BufSize: int32(c.bufSize),
	})
	if err != nil {
		return err
	}

	file, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
	if os.IsPermission(err) {
		err = makeReadWritable(localPath)
		if err != nil {
			return err
		}
		file, err = os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	defer file.Close()

	var destEnd int64
	if rfi != nil && lfi != nil {
		destEnd = lfi.Size()
		if rfi.Size() < destEnd {
			err = file.Truncate(rfi.Size())
			if err != nil {
				return err
			}
			destEnd = rfi.Size()
		}
	}

	var destPos int64
	var destBuf []byte
	if destPos < destEnd {
		destBuf = make([]byte, c.bufSize)
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if destPos < destEnd {
			destN, err := io.ReadFull(file, destBuf)
			if err == io.EOF {
				break
			}
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			destPos += int64(destN)

			if bytes.Equal(destBuf[:destN], chunk.Chunk) {
				continue
			}

			if destN > 0 {
				_, err := file.Seek(int64(-destN), os.SEEK_CUR)
				if err != nil {
					return err
				}
			}
		}

		_, err = file.Write(chunk.Chunk)
		if err != nil {
			return err
		}
	}

	if c.syncOwnerAndGroup {
		err = os.Chown(localPath, int(rfi.Uid()), int(rfi.Gid()))
		if err != nil {
			return err
		}
	}
	err = file.Chmod(rfi.Mode().Perm())
	if err != nil {
		return err
	}
	if c.syncModTime {
		err = os.Chtimes(localPath, time.Now(), rfi.ModTime())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) chown(ctx context.Context, remotePath string, uid, gid uint32) error {
	owner, err := c.userGroupDB.LookupUid(uid)
	if err != nil {
		return err
	}
	group, err := c.userGroupDB.LookupGid(gid)
	if err != nil {
		return err
	}
	_, err = c.client.Chown(ctx, &pb.ChownRequest{
		Path:  remotePath,
		Owner: owner,
		Group: group,
	})
	return err
}

func (c *Client) chmod(ctx context.Context, remotePath string, mode os.FileMode) error {
	_, err := c.client.Chmod(ctx, &pb.ChmodRequest{
		Path: remotePath,
		Mode: int32(mode.Perm()),
	})
	return err
}

func (c *Client) chtimes(ctx context.Context, remotePath string, atime, mtime time.Time) error {
	_, err := c.client.Chtimes(ctx, &pb.ChtimesRequest{
		Path:  remotePath,
		Atime: pb.ConvertTimeToPB(atime),
		Mtime: pb.ConvertTimeToPB(mtime),
	})
	return err
}

func (c *Client) readRemoteDir(ctx context.Context, remotePath string) ([]*fileInfo, error) {
	stream, err := c.client.ReadDir(ctx, &pb.ReadDirRequest{
		Path:               remotePath,
		AtMostCount:        int32(c.atMostCount),
		WantsOwnerAndGroup: c.syncOwnerAndGroup,
	})
	if err != nil {
		return nil, err
	}

	var allInfos []*pb.FileInfo
	for {
		infos, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		allInfos = append(allInfos, infos.Infos...)
	}

	infos, err := c.newFileInfosFromRPC(allInfos)
	if err != nil {
		return nil, err
	}
	sortFileInfosByName(infos)
	return infos, nil
}

func (c *Client) newFileInfosFromRPC(rpcFileInfos []*pb.FileInfo) ([]*fileInfo, error) {
	infos := make([]*fileInfo, 0, len(rpcFileInfos))
	for _, fi := range rpcFileInfos {
		info, err := c.newFileInfoFromRPC(fi)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	uid     uint32
	gid     uint32
}

func (fi fileInfo) Name() string { return fi.name }

func (fi fileInfo) Size() int64 { return fi.size }

func (fi fileInfo) Mode() os.FileMode { return fi.mode }

func (fi fileInfo) ModTime() time.Time { return fi.modTime }

func (fi fileInfo) IsDir() bool { return fi.Mode().IsDir() }

func (fi fileInfo) Sys() interface{} { return fi }

func (fi fileInfo) Uid() uint32 { return fi.uid }

func (fi fileInfo) Gid() uint32 { return fi.gid }

type fileInfosByName []*fileInfo

func (a fileInfosByName) Len() int      { return len(a) }
func (a fileInfosByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a fileInfosByName) Less(i, j int) bool {
	return strings.ToLower(a[i].Name()) < strings.ToLower(a[j].Name())
}

func sortFileInfosByName(infos []*fileInfo) {
	sort.Sort(fileInfosByName(infos))
}

func (c *Client) newFileInfoFromRPC(info *pb.FileInfo) (*fileInfo, error) {
	if info.Name == "" {
		return nil, nil
	}

	fi := &fileInfo{
		name:    info.Name,
		size:    info.Size,
		mode:    os.FileMode(info.Mode),
		modTime: pb.ConvertTimeFromPB(info.ModTime),
	}
	if !c.syncOwnerAndGroup {
		return fi, nil
	}

	uid, err := c.userGroupDB.LookupUser(info.Owner)
	if err != nil {
		return nil, err
	}
	gid, err := c.userGroupDB.LookupGroup(info.Group)
	if err != nil {
		return nil, err
	}
	fi.uid = uid
	fi.gid = gid
	return fi, nil
}

func (c *Client) readLocalDir(localPath string) ([]*fileInfo, error) {
	osInfos, err := readDir(localPath)
	if err != nil {
		return nil, err
	}
	infos, err := c.newFileInfosFromOS(osInfos)
	if err != nil {
		return nil, err
	}
	sortFileInfosByName(infos)
	return infos, nil
}

func (c *Client) newFileInfosFromOS(osFileInfos []os.FileInfo) ([]*fileInfo, error) {
	infos := make([]*fileInfo, 0, len(osFileInfos))
	for _, fi := range osFileInfos {
		info, err := c.newFileInfoFromOS(fi)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (c *Client) newFileInfoFromOS(info os.FileInfo) (*fileInfo, error) {
	fi := &fileInfo{
		name:    info.Name(),
		size:    info.Size(),
		mode:    info.Mode(),
		modTime: info.ModTime(),
	}
	if !c.syncOwnerAndGroup {
		return fi, nil
	}

	sys, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.New("unsupported file info sys type")
	}

	fi.uid = sys.Uid
	fi.gid = sys.Gid
	return fi, nil
}

func (c *Client) FetchDir(ctx context.Context, remotePath, localPath string) error {
	rfi, err := c.statRemote(ctx, remotePath)
	if err != nil {
		return err
	}

	if rfi.IsDir() {
		return fmt.Errorf("expected a remote directory but is a file %q", remotePath)
	}

	return c.fetchDirAndChmod(ctx, remotePath, localPath, rfi)
}

func (c *Client) fetchDirAndChmod(ctx context.Context, remotePath, localPath string, fi *fileInfo) error {
	g, ctx := errgroup.WithContext(ctx)
	fetchFileWorks := make(chan fetchFileWork)
	deleteWorks := make(chan deleteWork)
	postProcessDirWorks := make(chan fetchPostProcessDirWork)

	var walk func(ctx context.Context, remotePath, localPath string, fi *fileInfo) error
	walk = func(ctx context.Context, remotePath, localPath string, fi *fileInfo) error {
		log.Printf("walk start for localPath=%s", localPath)
		defer log.Printf("walk finish for localPath=%s", localPath)

		remoteInfos, err := c.readRemoteDir(ctx, remotePath)
		if err != nil {
			return err
		}

		err = ensureDirExists(localPath, 0777)
		if err != nil {
			return err
		}

		localInfos, err := c.readLocalDir(localPath)
		if err != nil {
			return err
		}

		li := 0
		dirWg := new(sync.WaitGroup)
		dirWg.Add(1)
		log.Printf("walk start loop for remoteInfos count=%d, localPath=%s", len(remoteInfos), localPath)
		for _, rfi := range remoteInfos {
			for li < len(localInfos) && localInfos[li].Name() < rfi.Name() {
				if !c.keepDeletedFiles {
					lfi := localInfos[li]
					log.Printf("walk sending deleteWork, localPath=%s, lfi.Name=%s", localPath, lfi.Name())
					dirWg.Add(1)
					work := deleteWork{
						localPath: filepath.Join(localPath, lfi.Name()),
						lfi:       lfi,
						dirWg:     dirWg,
					}
					select {
					case deleteWorks <- work:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				li++
			}

			var lfi *fileInfo
			if li < len(localInfos) && localInfos[li].Name() == rfi.Name() {
				lfi = localInfos[li]
				li++
				if lfi.IsDir() != rfi.IsDir() {
					err = ensureNotExist(filepath.Join(localPath, lfi.Name()), lfi)
					if err != nil {
						return err
					}
				}
			}

			if rfi.IsDir() {
				err = walk(ctx,
					filepath.Join(remotePath, rfi.Name()),
					filepath.Join(localPath, rfi.Name()),
					rfi)
				if err != nil {
					return err
				}
			} else {
				log.Printf("walk sending fetchFileWork, localPath=%s, rfi.Name=%s", localPath, rfi.Name())
				dirWg.Add(1)
				work := fetchFileWork{
					remotePath: filepath.Join(remotePath, rfi.Name()),
					localPath:  filepath.Join(localPath, rfi.Name()),
					rfi:        rfi,
					lfi:        lfi,
					dirWg:      dirWg,
				}
				select {
				case fetchFileWorks <- work:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		if !c.keepDeletedFiles {
			for li < len(localInfos) {
				lfi := localInfos[li]
				li++
				log.Printf("walk sending deleteWork, localPath=%s, lfi.Name=%s", localPath, lfi.Name())
				dirWg.Add(1)
				work := deleteWork{
					localPath: filepath.Join(localPath, lfi.Name()),
					lfi:       lfi,
					dirWg:     dirWg,
				}
				select {
				case deleteWorks <- work:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		log.Printf("walk before sending to postProcessDirWorks. localPath=%s", localPath)
		work := fetchPostProcessDirWork{
			dirWg:     dirWg,
			localPath: localPath,
			rfi:       fi,
		}
		select {
		case postProcessDirWorks <- work:
		case <-ctx.Done():
			return ctx.Err()
		}
		dirWg.Done()
		return nil
	}
	g.Go(func() error {
		defer log.Printf("start walking")
		defer close(fetchFileWorks)
		defer close(deleteWorks)
		defer close(postProcessDirWorks)
		defer log.Printf("finish walking")
		return walk(ctx, remotePath, localPath, fi)
	})

	const numFetchWorkers = 4
	for i := 0; i < numFetchWorkers; i++ {
		g.Go(func() error {
			for w := range fetchFileWorks {
				return func() error {
					log.Printf("fetchFileWorker start for localPath=%s", w.localPath)
					defer log.Printf("fetchFileWorker finish for localPath=%s", w.localPath)
					defer w.dirWg.Done()

					err := c.fetchFileAndChmod(ctx,
						w.remotePath,
						w.localPath,
						w.rfi,
						w.lfi)
					if err != nil {
						return err
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}
					return nil
				}()
			}
			return nil
		})
	}

	if !c.keepDeletedFiles {
		const numDeleteWorkers = 2
		for i := 0; i < numDeleteWorkers; i++ {
			g.Go(func() error {
				for w := range deleteWorks {
					return func() error {
						log.Printf("deleteWorker start for localPath=%s", w.localPath)
						defer log.Printf("deleteWorker finish for localPath=%s", w.localPath)
						defer w.dirWg.Done()

						err := ensureNotExist(w.localPath, w.lfi)
						if err != nil {
							return err
						}

						select {
						case <-ctx.Done():
							return ctx.Err()
						default:
						}
						return nil
					}()
				}
				return nil
			})
		}
	}

	g.Go(func() error {
		for w := range postProcessDirWorks {
			log.Printf("postProcessDirWorker start for localPath=%s", w.localPath)
			w.dirWg.Wait()
			log.Printf("postProcessDirWorker dirWg.Wait done for localPath=%s", w.localPath)
			if c.syncOwnerAndGroup {
				err := os.Chown(w.localPath, int(w.rfi.Uid()), int(w.rfi.Gid()))
				if err != nil {
					return err
				}
			}
			err := os.Chmod(w.localPath, w.rfi.Mode().Perm())
			if err != nil {
				return err
			}
			if c.syncModTime {
				err = os.Chtimes(w.localPath, time.Now(), w.rfi.ModTime())
				if err != nil {
					return err
				}
			}
			log.Printf("postProcessDirWorker finish for localPath=%s\n", w.localPath)
		}
		return nil
	})

	return g.Wait()
}

type fetchFileWork struct {
	remotePath string
	localPath  string
	rfi        *fileInfo
	lfi        *fileInfo
	dirWg      *sync.WaitGroup
}

type deleteWork struct {
	localPath string
	lfi       *fileInfo
	dirWg     *sync.WaitGroup
}

type fetchPostProcessDirWork struct {
	dirWg     *sync.WaitGroup
	localPath string
	rfi       *fileInfo
}

func (c *Client) Send(ctx context.Context, localPath, remotePath string) error {
	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(err) {
		lfi = nil
	} else if err != nil {
		return err
	}

	if lfi.IsDir() {
		return c.sendDirAndChmod(ctx, localPath, remotePath, lfi)
	} else {
		rfi, err := c.statRemote(ctx, remotePath)
		if err != nil {
			return err
		}

		return c.sendFileAndChmod(ctx, localPath, remotePath, lfi, rfi)
	}
}

func (c *Client) SendFile(ctx context.Context, localPath, remotePath string) error {
	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(err) {
		lfi = nil
	} else if err != nil {
		return err
	}

	if lfi.IsDir() {
		return fmt.Errorf("expected a local file but is a directory %q", localPath)
	}

	rfi, err := c.statRemote(ctx, remotePath)
	if err != nil {
		return err
	}

	return c.sendFileAndChmod(ctx, localPath, remotePath, lfi, rfi)
}

func (c *Client) sendFileAndChmod(ctx context.Context, localPath, remotePath string, lfi, rfi *fileInfo) error {
	if !c.isUpdateNeeded(lfi, rfi) {
		return nil
	}

	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	stream, err := c.client.SendFile(ctx)
	if err != nil {
		return err
	}
	err = stream.Send(&pb.SendFileRequest{Path: remotePath})
	if err != nil {
		return err
	}

	buf := make([]byte, c.bufSize)
	for {
		var n int
		n, err = io.ReadFull(file, buf)
		if err == io.EOF {
			err = nil
			break
		}
		if err == io.ErrUnexpectedEOF {
			buf = buf[:n]
		} else if err != nil {
			break
		}

		err = stream.Send(&pb.SendFileRequest{Chunk: buf})
		if err != nil {
			break
		}
	}
	_, err2 := stream.CloseAndRecv()
	if err != nil {
		return err
	} else if err2 != nil {
		return err2
	}

	if c.syncOwnerAndGroup {
		err = c.chown(ctx, localPath, lfi.Uid(), lfi.Gid())
		if err != nil {
			return err
		}
	}
	err = c.chmod(ctx, remotePath, lfi.Mode())
	if err != nil {
		return err
	}
	if c.syncModTime {
		err = c.chtimes(ctx, remotePath, time.Now(), lfi.ModTime())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ensureDirExists(ctx context.Context, remotePath string) error {
	_, err := c.client.EnsureDirExists(ctx, &pb.EnsureDirExistsRequest{Path: remotePath})
	return err
}

func (c *Client) ensureNotExist(ctx context.Context, remotePath string) error {
	_, err := c.client.EnsureNotExist(ctx, &pb.EnsureNotExistRequest{Path: remotePath})
	return err
}

func (c *Client) SendDir(ctx context.Context, localPath, remotePath string) error {
	lfi, err := c.statLocal(localPath)
	if err != nil {
		return err
	}

	if !lfi.IsDir() {
		return fmt.Errorf("expected a local directory but is a file %q", localPath)
	}

	return c.sendDirAndChmod(ctx, localPath, remotePath, lfi)
}

func (c *Client) sendDirAndChmod(ctx context.Context, localPath, remotePath string, fi *fileInfo) error {
	err := c.ensureDirExists(ctx, remotePath)
	if err != nil {
		return err
	}

	remoteInfos, err := c.readRemoteDir(ctx, remotePath)
	if err != nil {
		return err
	}

	localInfos, err := c.readLocalDir(localPath)
	if err != nil {
		return err
	}

	ri := 0
	for _, lfi := range localInfos {
		for ri < len(remoteInfos) && remoteInfos[ri].Name() < lfi.Name() {
			if !c.keepDeletedFiles {
				rfi := remoteInfos[ri]
				err = c.ensureNotExist(ctx, filepath.Join(remotePath, rfi.Name()))
				if err != nil {
					return err
				}
			}
			ri++
		}

		var rfi *fileInfo
		for ri < len(remoteInfos) && remoteInfos[ri].Name() == lfi.Name() {
			rfi = remoteInfos[ri]
			ri++
			if rfi.IsDir() != lfi.IsDir() {
				err = c.ensureNotExist(ctx, filepath.Join(remotePath, rfi.Name()))
				if err != nil {
					return err
				}
			}
		}

		if lfi.IsDir() {
			err = c.sendDirAndChmod(ctx,
				filepath.Join(localPath, lfi.Name()),
				filepath.Join(remotePath, lfi.Name()),
				lfi)
			if err != nil {
				return err
			}
		} else {
			err = c.sendFileAndChmod(ctx,
				filepath.Join(localPath, lfi.Name()),
				filepath.Join(remotePath, lfi.Name()),
				lfi,
				rfi)
			if err != nil {
				return err
			}
		}
	}

	if !c.keepDeletedFiles {
		for ri < len(remoteInfos) {
			rfi := remoteInfos[ri]
			ri++
			err = c.ensureNotExist(ctx, filepath.Join(remotePath, rfi.Name()))
			if err != nil {
				return err
			}
		}
	}

	if c.syncOwnerAndGroup {
		err = c.chown(ctx, localPath, fi.Uid(), fi.Gid())
		if err != nil {
			return err
		}
	}
	err = c.chmod(ctx, remotePath, fi.Mode())
	if err != nil {
		return err
	}
	if c.syncModTime {
		err = c.chtimes(ctx, remotePath, time.Now(), fi.ModTime())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) isUpdateNeeded(src, dest *fileInfo) bool {
	return !c.updateOnly || dest == nil || dest.Size() != src.Size() || dest.ModTime().Before(src.ModTime())
}
