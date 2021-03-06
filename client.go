package rdirsync

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/sync/errgroup"

	"github.com/hnakamur/rdirsync/internal"
	"github.com/hnakamur/rdirsync/internal/pb"
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
	fileWorkerCount   int
	userGroupDB       *userGroupDB
}

type ClientOption func(*Client) error

func NewClient(cc *grpc.ClientConn, option ...ClientOption) (*Client, error) {
	c := &Client{
		client:          pb.NewRDirSyncClient(cc),
		bufSize:         64 * 1024,
		atMostCount:     1024,
		fileWorkerCount: 2 * runtime.NumCPU(),
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

func SetFileWorkerCount(fileWorkerCount int) ClientOption {
	return func(c *Client) error {
		if fileWorkerCount <= 0 {
			return errors.New("file worker count must be greater than zero")
		}
		c.fileWorkerCount = fileWorkerCount
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
		return nil, errors.WithStack(err)
	}

	return c.newFileInfoFromOS(fi)
}

func (c *Client) Fetch(ctx context.Context, remotePath, localPath string) error {
	rfi, err := c.statRemote(ctx, remotePath)
	if err != nil {
		return err
	}

	if rfi == nil {
		return errors.Errorf("remote file or directory %q not found", remotePath)
	}

	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(errors.Cause(err)) {
		lfi = nil
	} else if err != nil {
		return err
	}
	if rfi.IsDir() {
		return c.fetchDirAndChmod(ctx, remotePath, localPath, rfi, lfi)
	} else {
		return c.fetchFileAndChmod(ctx, remotePath, localPath, rfi, lfi)
	}
}

func (c *Client) FetchFile(ctx context.Context, remotePath, localPath string) error {
	rfi, err := c.statRemote(ctx, remotePath)
	if err != nil {
		return err
	}

	if rfi.IsDir() {
		return errors.Errorf("expected a remote file but is a directory %q", remotePath)
	}

	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(errors.Cause(err)) {
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

	if lfi != nil && lfi.IsDir() {
		err := c.ensureLocalNotExist(localPath, lfi)
		if err != nil {
			return err
		}
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
		err = internal.MakeReadWritable(localPath)
		if err != nil {
			return err
		}
		file, err = os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return errors.WithStack(err)
		}
	} else if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	var destEnd int64
	if rfi != nil && lfi != nil {
		destEnd = lfi.Size()
		if rfi.Size() < destEnd {
			err = file.Truncate(rfi.Size())
			if err != nil {
				return errors.WithStack(err)
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
			return errors.WithStack(err)
		}

		if destPos < destEnd {
			destN, err := io.ReadFull(file, destBuf)
			if err == io.EOF {
				break
			}
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return errors.WithStack(err)
			}
			destPos += int64(destN)

			if bytes.Equal(destBuf[:destN], chunk.Chunk) {
				continue
			}

			if destN > 0 {
				_, err := file.Seek(int64(-destN), os.SEEK_CUR)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		_, err = file.Write(chunk.Chunk)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if c.syncOwnerAndGroup {
		err = os.Chown(localPath, int(rfi.Uid()), int(rfi.Gid()))
		if err != nil {
			return errors.WithStack(err)
		}
	}
	err = file.Chmod(rfi.Mode().Perm())
	if err != nil {
		return errors.WithStack(err)
	}
	if c.syncModTime {
		err = os.Chtimes(localPath, time.Now(), rfi.ModTime())
		if err != nil {
			return errors.WithStack(err)
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

func (c *Client) changeAttributes(ctx context.Context, remotePath string, changesOwner bool, uid, gid uint32, changesMode bool, mode os.FileMode, changesTime bool, atime, mtime time.Time) error {
	var owner string
	var group string
	var atimeNano int64
	var mtimeNano int64
	var err error
	if changesOwner {
		owner, err = c.userGroupDB.LookupUid(uid)
		if err != nil {
			return err
		}
		group, err = c.userGroupDB.LookupGid(gid)
		if err != nil {
			return err
		}
	}
	if changesTime {
		atimeNano = pb.ConvertTimeToPB(atime)
		mtimeNano = pb.ConvertTimeToPB(mtime)
	}
	_, err = c.client.ChangeAttributes(ctx, &pb.ChangeAttributesRequest{
		Path:         remotePath,
		ChangesOwner: changesOwner,
		ChangesMode:  changesMode,
		ChangesTime:  changesTime,
		Owner:        owner,
		Group:        group,
		Mode:         int32(mode.Perm()),
		Atime:        atimeNano,
		Mtime:        mtimeNano,
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
			return nil, errors.WithStack(err)
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

	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(errors.Cause(err)) {
		lfi = nil
	} else if err != nil {
		return err
	}
	return c.fetchDirAndChmod(ctx, remotePath, localPath, rfi, lfi)
}

func (c *Client) fetchDirAndChmod(ctx context.Context, remotePath, localPath string, rfi, lfi *fileInfo) error {
	g, ctx2 := errgroup.WithContext(ctx)
	fileWorks := make(chan fileWork)

	var walk func(ctx context.Context, remotePath, localPath string, rfi, lfi *fileInfo, treeNode *postProcessDirTreeNode) error
	walk = func(ctx context.Context, remotePath, localPath string, rfi, lfi *fileInfo, treeNode *postProcessDirTreeNode) error {
		remoteInfos, err := c.readRemoteDir(ctx, remotePath)
		if err != nil {
			return err
		}

		err = c.ensureLocalDirExists(localPath, 0700, lfi)
		if err != nil {
			return err
		}

		localInfos, err := c.readLocalDir(localPath)
		if err != nil {
			return err
		}

		li := 0
		for _, rfi := range remoteInfos {
			for li < len(localInfos) && localInfos[li].Name() < rfi.Name() {
				if !c.keepDeletedFiles {
					lfi := localInfos[li]
					work := deleteWork(
						filepath.Join(localPath, lfi.Name()),
						lfi,
					)
					select {
					case fileWorks <- work:
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
			}

			if rfi.IsDir() {
				childNode := &postProcessDirTreeNode{
					path:    filepath.Join(localPath, rfi.Name()),
					mode:    rfi.Mode(),
					modTime: rfi.ModTime(),
					uid:     rfi.Uid(),
					gid:     rfi.Gid(),
				}
				treeNode.children = append(treeNode.children, childNode)
				err = walk(ctx,
					filepath.Join(remotePath, rfi.Name()),
					filepath.Join(localPath, rfi.Name()),
					rfi,
					lfi,
					childNode)
				if err != nil {
					return err
				}
			} else {
				work := fetchFileWork(
					filepath.Join(remotePath, rfi.Name()),
					filepath.Join(localPath, rfi.Name()),
					rfi,
					lfi,
				)
				select {
				case fileWorks <- work:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		if !c.keepDeletedFiles {
			for li < len(localInfos) {
				lfi := localInfos[li]
				li++
				work := deleteWork(
					filepath.Join(localPath, lfi.Name()),
					lfi,
				)
				select {
				case fileWorks <- work:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	}

	treeRoot := &postProcessDirTreeNode{
		path:    localPath,
		mode:    rfi.Mode(),
		modTime: rfi.ModTime(),
		uid:     rfi.Uid(),
		gid:     rfi.Gid(),
	}
	g.Go(func() error {
		defer close(fileWorks)
		return walk(ctx2, remotePath, localPath, rfi, lfi, treeRoot)
	})

	for i := 0; i < c.fileWorkerCount; i++ {
		g.Go(func() error {
			for w := range fileWorks {
				err := w(c, ctx)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	var postWalk func(ctx context.Context, n *postProcessDirTreeNode) error
	postWalk = func(ctx context.Context, n *postProcessDirTreeNode) error {
		for _, child := range n.children {
			err := postWalk(ctx, child)
			if err != nil {
				return err
			}
		}

		if c.syncOwnerAndGroup {
			err := os.Chown(n.path, int(n.uid), int(n.gid))
			if err != nil {
				return errors.WithStack(err)
			}
		}
		err := os.Chmod(n.path, n.mode.Perm())
		if err != nil {
			return errors.WithStack(err)
		}
		if c.syncModTime {
			err = os.Chtimes(n.path, time.Now(), n.modTime)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}

	return postWalk(ctx, treeRoot)
}

type fileWork func(c *Client, ctx context.Context) error

func fetchFileWork(remotePath, localPath string, rfi, lfi *fileInfo) fileWork {
	return func(c *Client, ctx context.Context) error {
		return c.fetchFileAndChmod(ctx,
			remotePath,
			localPath,
			rfi,
			lfi)
	}
}

func deleteWork(localPath string, lfi *fileInfo) fileWork {
	return func(c *Client, ctx context.Context) error {
		return c.ensureLocalNotExist(localPath, lfi)
	}
}

type postProcessDirTreeNode struct {
	path     string
	mode     os.FileMode
	modTime  time.Time
	uid      uint32
	gid      uint32
	children []*postProcessDirTreeNode
}

func (c *Client) Send(ctx context.Context, localPath, remotePath string) error {
	lfi, err := c.statLocal(localPath)
	if os.IsNotExist(errors.Cause(err)) {
		lfi = nil
	} else if err != nil {
		return errors.WithStack(err)
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
	if os.IsNotExist(errors.Cause(err)) {
		lfi = nil
	} else if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
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

	return c.changeAttributes(ctx, remotePath,
		c.syncOwnerAndGroup, lfi.Uid(), lfi.Gid(),
		true, lfi.Mode(),
		c.syncModTime, time.Now(), lfi.ModTime())
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
		return errors.Errorf("expected a local directory but is a file %q", localPath)
	}

	return c.sendDirAndChmod(ctx, localPath, remotePath, lfi)
}

func (c *Client) sendDirAndChmod(ctx context.Context, localPath, remotePath string, lfi *fileInfo) error {
	g, ctx2 := errgroup.WithContext(ctx)
	fileWorks := make(chan fileWork)

	var walk func(ctx context.Context, localPath, remotePath string, lfi *fileInfo, treeNode *postProcessDirTreeNode) error
	walk = func(ctx context.Context, localPath, remotePath string, lfi *fileInfo, treeNode *postProcessDirTreeNode) error {
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
					work := ensureRemoteNotExistWork(filepath.Join(remotePath, rfi.Name()))
					select {
					case fileWorks <- work:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				ri++
			}

			var rfi *fileInfo
			for ri < len(remoteInfos) && remoteInfos[ri].Name() == lfi.Name() {
				rfi = remoteInfos[ri]
				ri++
			}

			if lfi.IsDir() {
				childNode := &postProcessDirTreeNode{
					path:    filepath.Join(remotePath, lfi.Name()),
					mode:    lfi.Mode(),
					modTime: lfi.ModTime(),
					uid:     lfi.Uid(),
					gid:     lfi.Gid(),
				}
				treeNode.children = append(treeNode.children, childNode)
				err = walk(ctx,
					filepath.Join(localPath, lfi.Name()),
					filepath.Join(remotePath, lfi.Name()),
					lfi,
					childNode)
				if err != nil {
					return err
				}
			} else {
				work := sendFileWork(
					filepath.Join(localPath, lfi.Name()),
					filepath.Join(remotePath, lfi.Name()),
					lfi,
					rfi)
				select {
				case fileWorks <- work:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		if !c.keepDeletedFiles {
			for ri < len(remoteInfos) {
				rfi := remoteInfos[ri]
				ri++
				work := ensureRemoteNotExistWork(filepath.Join(remotePath, rfi.Name()))
				select {
				case fileWorks <- work:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		return nil
	}

	treeRoot := &postProcessDirTreeNode{
		path:    remotePath,
		mode:    lfi.Mode(),
		modTime: lfi.ModTime(),
		uid:     lfi.Uid(),
		gid:     lfi.Gid(),
	}
	g.Go(func() error {
		defer close(fileWorks)
		return walk(ctx2, localPath, remotePath, lfi, treeRoot)
	})

	for i := 0; i < c.fileWorkerCount; i++ {
		g.Go(func() error {
			for w := range fileWorks {
				err := w(c, ctx)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	var postWalk func(ctx context.Context, n *postProcessDirTreeNode) error
	postWalk = func(ctx context.Context, n *postProcessDirTreeNode) error {
		for _, child := range n.children {
			err := postWalk(ctx, child)
			if err != nil {
				return err
			}
		}

		return c.changeAttributes(ctx, n.path,
			c.syncOwnerAndGroup, n.uid, n.gid,
			true, n.mode.Perm(),
			c.syncModTime, time.Now(), n.modTime)
	}

	return postWalk(ctx, treeRoot)
}

func sendFileWork(localPath, remotePath string, lfi, rfi *fileInfo) fileWork {
	return func(c *Client, ctx context.Context) error {
		return c.sendFileAndChmod(ctx,
			localPath,
			remotePath,
			lfi,
			rfi)
	}
}

func ensureRemoteNotExistWork(remotePath string) fileWork {
	return func(c *Client, ctx context.Context) error {
		return c.ensureNotExist(ctx, remotePath)
	}
}

func (c *Client) isUpdateNeeded(src, dest *fileInfo) bool {
	return !c.updateOnly || dest == nil || dest.Size() != src.Size() || dest.ModTime().Before(src.ModTime())
}

func (c *Client) ensureLocalDirExists(path string, mode os.FileMode, fi *fileInfo) error {
	if fi != nil {
		if fi.IsDir() {
			return nil
		} else {
			err := internal.EnsureFileNotExist(path)
			if err != nil {
				return err
			}
		}
	}

	err := os.MkdirAll(path, mode.Perm())
	if err == nil {
		return nil
	} else if !os.IsPermission(err) {
		return errors.WithStack(err)
	}

	err = internal.MakeReadWritableRecursive(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, mode.Perm())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Client) ensureLocalNotExist(path string, fi *fileInfo) error {
	if fi == nil {
		return nil
	}

	if fi.IsDir() {
		return internal.EnsureDirNotExist(path)
	} else {
		return internal.EnsureFileNotExist(path)
	}
}
