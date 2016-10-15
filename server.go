package rdirsync

import (
	"io"
	"os"

	"bitbucket.org/hnakamur/rdirsync/rpc"
)

type server struct{}

func NewServer() rpc.RDirSyncServer {
	return new(server)
}

func (s *server) FetchFile(req *rpc.FetchFileRequest, stream rpc.RDirSync_FetchFileServer) error {
	buf := make([]byte, req.BufSize)
	file, err := os.Open(req.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		n, err := io.ReadFull(file, buf)
		if err == io.EOF {
			break
		} else if err == io.ErrUnexpectedEOF {
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
		} else if err != nil {
			return err
		}

		infos := newFileInfosFromOS(osFileInfos)
		err = stream.Send(&rpc.FileInfos{Infos: infos})
		if err != nil {
			return err
		}
	}

	return nil
}

func newFileInfosFromOS(fis []os.FileInfo) []*rpc.FileInfo {
	infos := make([]*rpc.FileInfo, 0, len(fis))
	for _, fi := range fis {
		if fi.IsDir() || fi.Mode().IsRegular() {
			infos = append(infos, newFileInfoFromOS(fi))
		}
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
