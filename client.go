package rdirsync

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"bitbucket.org/hnakamur/rdirsync/rpc"
	"google.golang.org/grpc"
)

type ClientFacade struct {
	client      rpc.RDirSyncClient
	bufSize     int
	atMostCount int
}

func NewClientFacade(cc *grpc.ClientConn, bufSize, atMostCount int) *ClientFacade {
	if bufSize == 0 {
		bufSize = 64 * 1024
	}
	if atMostCount == 0 {
		atMostCount = 1024
	}
	return &ClientFacade{
		client:      rpc.NewRDirSyncClient(cc),
		bufSize:     bufSize,
		atMostCount: atMostCount,
	}
}

func (c *ClientFacade) FetchFileToWriter(ctx context.Context, remotePath string, w io.Writer) error {
	stream, err := c.client.FetchFile(ctx, &rpc.FetchRequest{
		Path:    remotePath,
		BufSize: int32(c.bufSize),
	})
	if err != nil {
		return err
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("FetchFileToWriter. w=%+v, chunk=%+v", w, chunk)
		_, err = w.Write(chunk.Chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientFacade) FetchFile(ctx context.Context, remotePath, localPath string) error {
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return c.FetchFileToWriter(ctx, remotePath, file)
}

func (c *ClientFacade) ReadDirToQueue(ctx context.Context, remotePath string, out chan<- os.FileInfo) error {
	stream, err := c.client.ReadDir(ctx, &rpc.ReadDirRequest{
		Path:        remotePath,
		AtMostCount: int32(c.atMostCount),
	})
	if err != nil {
		return err
	}
	for {
		infos, err := stream.Recv()
		if err == io.EOF {
			break
		}
		for _, info := range infos.Infos {
			out <- newRemoteFileInfoFromRPC(info)
		}
	}
	return nil
}

type RemoteFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi RemoteFileInfo) Name() string { return fi.name }

func (fi RemoteFileInfo) Size() int64 { return fi.size }

func (fi RemoteFileInfo) Mode() os.FileMode { return fi.mode }

func (fi RemoteFileInfo) ModTime() time.Time { return fi.modTime }

func (fi RemoteFileInfo) IsDir() bool { return fi.Mode().IsDir() }

func (fi RemoteFileInfo) Sys() interface{} { return nil }

func newRemoteFileInfoFromRPC(info *rpc.FileInfo) *RemoteFileInfo {
	return &RemoteFileInfo{
		name:    info.Name,
		size:    info.Size,
		mode:    os.FileMode(info.Mode),
		modTime: time.Unix(info.ModTime, 0),
	}
}
