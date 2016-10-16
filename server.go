package rdirsync

import (
	"io"
	"os"
	"time"

	"google.golang.org/grpc"

	context "golang.org/x/net/context"

	"github.com/hnakamur/rdirsync/pb"
)

type server struct{}

func RegisterNewRDirSyncServer(s *grpc.Server) {
	pb.RegisterRDirSyncServer(s, newServer())
}

func newServer() pb.RDirSyncServer {
	return new(server)
}

func (s *server) Stat(ctx context.Context, req *pb.StatRequest) (*pb.FileInfo, error) {
	fi, err := os.Stat(req.Path)
	if err != nil {
		return nil, err
	}
	info := newFileInfoFromOS(fi)
	return info, nil
}

func (s *server) FetchFile(req *pb.FetchFileRequest, stream pb.RDirSync_FetchFileServer) error {
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

		err = stream.Send(&pb.FileChunk{Chunk: buf})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ReadDir(req *pb.ReadDirRequest, stream pb.RDirSync_ReadDirServer) error {
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
		err = stream.Send(&pb.FileInfos{Infos: infos})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) Chmod(ctx context.Context, req *pb.ChmodRequest) (*pb.Empty, error) {
	err := os.Chmod(req.Path, os.FileMode(req.Mode).Perm())
	return new(pb.Empty), err
}

func (s *server) Chtimes(ctx context.Context, req *pb.ChtimesRequest) (*pb.Empty, error) {
	err := os.Chtimes(req.Path, time.Unix(req.Atime, 0), time.Unix(req.Mtime, 0))
	return new(pb.Empty), err
}

func (s *server) EnsureDirExists(ctx context.Context, req *pb.EnsureDirExistsRequest) (*pb.Empty, error) {
	err := ensureDirExists(req.Path, 0777)
	return new(pb.Empty), err
}

func (s *server) EnsureNotExist(ctx context.Context, req *pb.EnsureNotExistRequest) (*pb.Empty, error) {
	err := ensureNotExist(req.Path, nil)
	return new(pb.Empty), err
}

func (s *server) SendFile(stream pb.RDirSync_SendFileServer) error {
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
	return stream.SendAndClose(new(pb.Empty))
}

func newFileInfosFromOS(fis []os.FileInfo) []*pb.FileInfo {
	infos := make([]*pb.FileInfo, 0, len(fis))
	for _, fi := range fis {
		infos = append(infos, newFileInfoFromOS(fi))
	}
	return infos
}

func newFileInfoFromOS(fi os.FileInfo) *pb.FileInfo {
	return &pb.FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    int32(fi.Mode()),
		ModTime: fi.ModTime().Unix(),
	}
}
