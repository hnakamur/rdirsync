package rdirsync

import (
	"io"
	"os"
	"time"

	context "golang.org/x/net/context"

	"github.com/hnakamur/rdirsync/rpc"
)

type server struct{}

func NewServer() rpc.RDirSyncServer {
	return new(server)
}

func (s *server) Stat(ctx context.Context, req *rpc.StatRequest) (*rpc.FileInfo, error) {
	fi, err := os.Stat(req.Path)
	if err != nil {
		return nil, err
	}
	info := newFileInfoFromOS(fi)
	return info, nil
}

func (s *server) FetchFile(req *rpc.FetchFileRequest, stream rpc.RDirSync_FetchFileServer) error {
	file, err := os.Open(req.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, req.BufSize)
	for {
		n, err := io.ReadFull(file, buf)
		if err == io.EOF {
			break
		}
		if err == io.ErrUnexpectedEOF {
			buf = buf[:n]
		} else if err != nil {
			return err
		}

		err = stream.Send(&rpc.FileChunk{Chunk: buf})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ReadDir(req *rpc.ReadDirRequest, stream rpc.RDirSync_ReadDirServer) error {
	file, err := os.Open(req.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		osFileInfos, err := file.Readdir(int(req.AtMostCount))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		infos := newFileInfosFromOS(selectDirAndRegularFiles(osFileInfos))
		err = stream.Send(&rpc.FileInfos{Infos: infos})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) Chmod(ctx context.Context, req *rpc.ChmodRequest) (*rpc.Empty, error) {
	err := os.Chmod(req.Path, os.FileMode(req.Mode).Perm())
	return new(rpc.Empty), err
}

func (s *server) Chtimes(ctx context.Context, req *rpc.ChtimesRequest) (*rpc.Empty, error) {
	err := os.Chtimes(req.Path, time.Unix(req.Atime, 0), time.Unix(req.Mtime, 0))
	return new(rpc.Empty), err
}

func (s *server) EnsureDirExists(ctx context.Context, req *rpc.EnsureDirExistsRequest) (*rpc.Empty, error) {
	err := ensureDirExists(req.Path, 0777)
	return new(rpc.Empty), err
}

func (s *server) EnsureNotExist(ctx context.Context, req *rpc.EnsureNotExistRequest) (*rpc.Empty, error) {
	err := ensureNotExist(req.Path, nil)
	return new(rpc.Empty), err
}

func (s *server) SendFile(stream rpc.RDirSync_SendFileServer) error {
	var file *os.File
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if file == nil {
			err = ensureNotDir(chunk.Path, nil)
			if err != nil {
				return err
			}

			file, err = os.Create(chunk.Path)
			if os.IsPermission(err) {
				err = makeReadWritable(chunk.Path)
				if err != nil {
					return err
				}
				file, err = os.Create(chunk.Path)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			defer file.Close()
		}

		if len(chunk.Chunk) > 0 {
			_, err = file.Write(chunk.Chunk)
			if err != nil {
				return err
			}
		}
	}
	return stream.SendAndClose(new(rpc.Empty))
}

func newFileInfosFromOS(fis []os.FileInfo) []*rpc.FileInfo {
	infos := make([]*rpc.FileInfo, 0, len(fis))
	for _, fi := range fis {
		infos = append(infos, newFileInfoFromOS(fi))
	}
	return infos
}

func newFileInfoFromOS(fi os.FileInfo) *rpc.FileInfo {
	return &rpc.FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    int32(fi.Mode()),
		ModTime: fi.ModTime().Unix(),
	}
}
