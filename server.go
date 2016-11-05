package rdirsync

import (
	"bytes"
	"io"
	"os"
	"syscall"

	"google.golang.org/grpc"

	context "golang.org/x/net/context"

	"github.com/hnakamur/rdirsync/internal"
	"github.com/hnakamur/rdirsync/internal/pb"
	"github.com/pkg/errors"
)

type server struct {
	userGroupDB *userGroupDB
}

func RegisterNewRDirSyncServer(s *grpc.Server) {
	pb.RegisterRDirSyncServer(s, newServer())
}

func newServer() pb.RDirSyncServer {
	return &server{
		userGroupDB: newUserGroupDB(),
	}
}

func (s *server) Stat(ctx context.Context, req *pb.StatRequest) (*pb.FileInfo, error) {
	fi, err := os.Stat(req.Path)
	if os.IsNotExist(err) {
		fi = nil
	} else if err != nil {
		return nil, err
	}
	return s.newFileInfoFromOS(fi, req.WantsOwnerAndGroup)
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

		infos, err := s.newFileInfosFromOS(selectDirAndRegularFiles(osFileInfos),
			req.WantsOwnerAndGroup)
		if err != nil {
			return err
		}
		err = stream.Send(&pb.FileInfos{Infos: infos})
		if err != nil {
			return err
		}
	}

	return nil
}

func selectDirAndRegularFiles(fis []os.FileInfo) []os.FileInfo {
	ret := make([]os.FileInfo, 0, len(fis))
	for _, fi := range fis {
		if fi.IsDir() || fi.Mode().IsRegular() {
			ret = append(ret, fi)
		}
	}
	return ret
}

func (s *server) Chown(ctx context.Context, req *pb.ChownRequest) (*pb.Empty, error) {
	uid, err := s.userGroupDB.LookupUser(req.Owner)
	if err != nil {
		return new(pb.Empty), err
	}
	gid, err := s.userGroupDB.LookupGroup(req.Group)
	if err != nil {
		return new(pb.Empty), err
	}
	err = os.Chown(req.Path, int(uid), int(gid))
	return new(pb.Empty), err
}

func (s *server) Chmod(ctx context.Context, req *pb.ChmodRequest) (*pb.Empty, error) {
	err := os.Chmod(req.Path, os.FileMode(req.Mode).Perm())
	return new(pb.Empty), err
}

func (s *server) Chtimes(ctx context.Context, req *pb.ChtimesRequest) (*pb.Empty, error) {
	err := os.Chtimes(req.Path,
		pb.ConvertTimeFromPB(req.Atime),
		pb.ConvertTimeFromPB(req.Mtime))
	return new(pb.Empty), err
}

func (s *server) EnsureDirExists(ctx context.Context, req *pb.EnsureDirExistsRequest) (*pb.Empty, error) {
	err := internal.EnsureDirExists(req.Path, 0700)
	return new(pb.Empty), err
}

func (s *server) EnsureNotExist(ctx context.Context, req *pb.EnsureNotExistRequest) (*pb.Empty, error) {
	err := internal.EnsureDirOrFileNotExist(req.Path)
	return new(pb.Empty), err
}

func (s *server) SendFile(stream pb.RDirSync_SendFileServer) error {
	var file *os.File
	var destPos int64
	var destEnd int64
	var destBuf []byte
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if file == nil {
			err = internal.EnsureNotDir(chunk.Path)
			if err != nil {
				return err
			}

			file, err = os.OpenFile(chunk.Path, os.O_RDWR|os.O_CREATE, 0666)
			if os.IsPermission(err) {
				err = internal.MakeReadWritable(chunk.Path)
				if err != nil {
					return err
				}
				file, err = os.OpenFile(chunk.Path, os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					return errors.WithStack(err)
				}
			} else if err != nil {
				return errors.WithStack(err)
			}
			defer file.Close()

			fi, err := file.Stat()
			if err != nil {
				return errors.WithStack(err)
			}
			destEnd = fi.Size()

			if destEnd > 0 {
				destBuf = make([]byte, len(chunk.Chunk))
			}
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

		if len(chunk.Chunk) > 0 {
			n, err := file.Write(chunk.Chunk)
			if err != nil {
				return err
			}
			destPos += int64(n)
		}
	}

	if destPos < destEnd {
		err := file.Truncate(destPos)
		if err != nil {
			return err
		}
	}

	return stream.SendAndClose(new(pb.Empty))
}

func (s *server) newFileInfosFromOS(fis []os.FileInfo, wantsOwnerAndGroup bool) ([]*pb.FileInfo, error) {
	infos := make([]*pb.FileInfo, 0, len(fis))
	for _, fi := range fis {
		pbInfo, err := s.newFileInfoFromOS(fi, wantsOwnerAndGroup)
		if err != nil {
			return nil, err
		}
		infos = append(infos, pbInfo)
	}
	return infos, nil
}

func (s *server) newFileInfoFromOS(fi os.FileInfo, wantsOwnerAndGroup bool) (*pb.FileInfo, error) {
	if fi == nil {
		return new(pb.FileInfo), nil
	}

	info := &pb.FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    int32(fi.Mode()),
		ModTime: pb.ConvertTimeToPB(fi.ModTime()),
	}
	if !wantsOwnerAndGroup {
		return info, nil
	}

	sysInfo := fi.Sys().(*syscall.Stat_t)
	owner, err := s.userGroupDB.LookupUid(sysInfo.Uid)
	if err != nil {
		return nil, err
	}
	group, err := s.userGroupDB.LookupGid(sysInfo.Gid)
	if err != nil {
		return nil, err
	}
	info.Owner = owner
	info.Group = group
	return info, nil
}
