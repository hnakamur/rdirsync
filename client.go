package rdirsync

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hnakamur/rdirsync/rpc"
	"google.golang.org/grpc"
)

type ClientFacade struct {
	client           rpc.RDirSyncClient
	bufSize          int
	atMostCount      int
	keepDeletedFiles bool
	syncModTime      bool
}

func NewClientFacade(cc *grpc.ClientConn, option ...func(*ClientFacade)) *ClientFacade {
	c := &ClientFacade{
		client:      rpc.NewRDirSyncClient(cc),
		bufSize:     64 * 1024,
		atMostCount: 1024,
	}

	for _, opt := range option {
		opt(c)
	}
	return c
}

func SetBufSize(bufSize int) func(*ClientFacade) {
	return func(c *ClientFacade) {
		c.bufSize = bufSize
	}
}

func SetMaxEntriesPerRPC(maxEntriesPerRPC int) func(*ClientFacade) {
	return func(c *ClientFacade) {
		c.atMostCount = maxEntriesPerRPC
	}
}

func SetKeepDeletedFiles(keepDeletedFiles bool) func(*ClientFacade) {
	return func(c *ClientFacade) {
		c.keepDeletedFiles = keepDeletedFiles
	}
}

func SetSyncModTime(syncModTime bool) func(*ClientFacade) {
	return func(c *ClientFacade) {
		c.syncModTime = syncModTime
	}
}

func (c *ClientFacade) stat(ctx context.Context, remotePath string) (os.FileInfo, error) {
	info, err := c.client.Stat(ctx, &rpc.StatRequest{Path: remotePath})
	if err != nil {
		return nil, err
	}
	return newFileInfoFromRPC(info), nil
}

func (c *ClientFacade) FetchFile(ctx context.Context, remotePath, localPath string) error {
	fi, err := c.stat(ctx, remotePath)
	if err != nil {
		return err
	}
	return c.fetchFileAndChmod(ctx, remotePath, localPath, fi)
}

func (c *ClientFacade) fetchFileAndChmod(ctx context.Context, remotePath, localPath string, fi os.FileInfo) error {
	stream, err := c.client.FetchFile(ctx, &rpc.FetchFileRequest{
		Path:    remotePath,
		BufSize: int32(c.bufSize),
	})
	if err != nil {
		return err
	}

	file, err := os.Create(localPath)
	if os.IsPermission(err) {
		err = makeReadWritable(localPath)
		if err != nil {
			return err
		}
		file, err = os.Create(localPath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	defer file.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		_, err = file.Write(chunk.Chunk)
		if err != nil {
			return err
		}
	}

	err = file.Chmod(fi.Mode().Perm())
	if err != nil {
		return err
	}
	if c.syncModTime {
		err = os.Chtimes(localPath, time.Now(), fi.ModTime())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientFacade) chmod(ctx context.Context, remotePath string, mode os.FileMode) error {
	_, err := c.client.Chmod(ctx,
		&rpc.ChmodRequest{
			Path: remotePath,
			Mode: int32(mode.Perm())})
	return err
}

func (c *ClientFacade) chtimes(ctx context.Context, remotePath string, atime, mtime time.Time) error {
	_, err := c.client.Chtimes(ctx,
		&rpc.ChtimesRequest{
			Path:  remotePath,
			Atime: atime.Unix(),
			Mtime: mtime.Unix()})
	return err
}

func (c *ClientFacade) readDir(ctx context.Context, remotePath string) ([]os.FileInfo, error) {
	stream, err := c.client.ReadDir(ctx, &rpc.ReadDirRequest{
		Path:        remotePath,
		AtMostCount: int32(c.atMostCount),
	})
	if err != nil {
		return nil, err
	}

	var allInfos []*rpc.FileInfo
	for {
		infos, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		allInfos = append(allInfos, infos.Infos...)
	}

	infos := convertRPCFileInfosToOSFileInfos(allInfos)
	sortFileInfosByName(infos)
	return infos, nil
}

func convertRPCFileInfosToOSFileInfos(rpcFileInfos []*rpc.FileInfo) []os.FileInfo {
	infos := make([]os.FileInfo, 0, len(rpcFileInfos))
	for _, info := range rpcFileInfos {
		infos = append(infos, newFileInfoFromRPC(info))
	}
	return infos
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi fileInfo) Name() string { return fi.name }

func (fi fileInfo) Size() int64 { return fi.size }

func (fi fileInfo) Mode() os.FileMode { return fi.mode }

func (fi fileInfo) ModTime() time.Time { return fi.modTime }

func (fi fileInfo) IsDir() bool { return fi.Mode().IsDir() }

func (fi fileInfo) Sys() interface{} { return nil }

func newFileInfoFromRPC(info *rpc.FileInfo) *fileInfo {
	return &fileInfo{
		name:    info.Name,
		size:    info.Size,
		mode:    os.FileMode(info.Mode),
		modTime: time.Unix(info.ModTime, 0),
	}
}

func (c *ClientFacade) FetchDir(ctx context.Context, remotePath, localPath string) error {
	fi, err := c.stat(ctx, remotePath)
	if err != nil {
		return err
	}
	return c.fetchDirAndChmod(ctx, remotePath, localPath, fi)
}

func (c *ClientFacade) fetchDirAndChmod(ctx context.Context, remotePath, localPath string, fi os.FileInfo) error {
	remoteInfos, err := c.readDir(ctx, remotePath)
	if err != nil {
		return err
	}

	err = ensureDirExists(localPath, 0777)
	if err != nil {
		return err
	}

	localInfos, err := readLocalDir(localPath)
	if err != nil {
		return err
	}

	li := 0
	for _, rfi := range remoteInfos {
		for li < len(localInfos) && localInfos[li].Name() < rfi.Name() {
			if !c.keepDeletedFiles {
				lfi := localInfos[li]
				err = ensureNotExist(filepath.Join(localPath, lfi.Name()), lfi)
				if err != nil {
					return err
				}
			}
			li++
		}

		if li < len(localInfos) && localInfos[li].Name() == rfi.Name() {
			lfi := localInfos[li]
			li++
			if lfi.IsDir() != rfi.IsDir() {
				err = ensureNotExist(filepath.Join(localPath, lfi.Name()), lfi)
				if err != nil {
					return err
				}
			}
		}

		if rfi.IsDir() {
			err = c.fetchDirAndChmod(ctx,
				filepath.Join(remotePath, rfi.Name()),
				filepath.Join(localPath, rfi.Name()),
				rfi)
			if err != nil {
				return err
			}
		} else {
			err = c.fetchFileAndChmod(ctx,
				filepath.Join(remotePath, rfi.Name()),
				filepath.Join(localPath, rfi.Name()),
				rfi)
			if err != nil {
				return err
			}
		}
	}

	if !c.keepDeletedFiles {
		for li < len(localInfos) {
			lfi := localInfos[li]
			li++
			err = ensureNotExist(filepath.Join(localPath, lfi.Name()), lfi)
			if err != nil {
				return err
			}
		}
	}

	err = os.Chmod(localPath, fi.Mode().Perm())
	if err != nil {
		return err
	}
	if c.syncModTime {
		err = os.Chtimes(localPath, time.Now(), fi.ModTime())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientFacade) SendFile(ctx context.Context, localPath, remotePath string) error {
	fi, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	return c.sendFileAndChmod(ctx, localPath, remotePath, fi)
}

func (c *ClientFacade) sendFileAndChmod(ctx context.Context, localPath, remotePath string, fi os.FileInfo) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	stream, err := c.client.SendFile(ctx)
	if err != nil {
		return err
	}
	err = stream.Send(&rpc.SendFileRequest{Path: remotePath})
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

		err = stream.Send(&rpc.SendFileRequest{Chunk: buf})
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

func (c *ClientFacade) ensureDirExists(ctx context.Context, remotePath string) error {
	_, err := c.client.EnsureDirExists(ctx, &rpc.EnsureDirExistsRequest{Path: remotePath})
	return err
}

func (c *ClientFacade) ensureNotExist(ctx context.Context, remotePath string) error {
	_, err := c.client.EnsureNotExist(ctx, &rpc.EnsureNotExistRequest{Path: remotePath})
	return err
}

func (c *ClientFacade) SendDir(ctx context.Context, localPath, remotePath string) error {
	fi, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	return c.sendDirAndChmod(ctx, localPath, remotePath, fi)
}

func (c *ClientFacade) sendDirAndChmod(ctx context.Context, localPath, remotePath string, fi os.FileInfo) error {
	err := c.ensureDirExists(ctx, remotePath)
	if err != nil {
		return err
	}

	remoteInfos, err := c.readDir(ctx, remotePath)
	if err != nil {
		return err
	}

	localInfos, err := readLocalDir(localPath)
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

		for ri < len(remoteInfos) && remoteInfos[ri].Name() == lfi.Name() {
			rfi := remoteInfos[ri]
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
				lfi)
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
