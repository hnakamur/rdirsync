package rdirsync

import (
	"context"
	"io"
	"os"
	"sort"
	"time"

	"bitbucket.org/hnakamur/rdirsync/rpc"
	"google.golang.org/grpc"
)

type ClientFacadeConfig struct {
	BufSize          int
	MaxEntriesPerRPC int
}

type ClientFacade struct {
	client      rpc.RDirSyncClient
	bufSize     int
	atMostCount int
}

func NewClientFacade(cc *grpc.ClientConn, config *ClientFacadeConfig) *ClientFacade {
	var bufSize int
	if config != nil && config.BufSize > 0 {
		bufSize = config.BufSize
	} else {
		bufSize = 64 * 1024
	}

	var atMostCount int
	if config != nil && config.MaxEntriesPerRPC > 0 {
		atMostCount = config.MaxEntriesPerRPC
	} else {
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
		} else if err != nil {
			return err
		}
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

func (c *ClientFacade) ReadDir(ctx context.Context, remotePath string) ([]os.FileInfo, error) {
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
		infos = append(infos,
			&fileInfo{
				name:    info.Name,
				size:    info.Size,
				mode:    os.FileMode(info.Mode),
				modTime: time.Unix(info.ModTime, 0),
			})
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

type osFileInfosByName []os.FileInfo

func (a osFileInfosByName) Len() int           { return len(a) }
func (a osFileInfosByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a osFileInfosByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func sortFileInfosByName(infos []os.FileInfo) {
	sort.Sort(osFileInfosByName(infos))
}
