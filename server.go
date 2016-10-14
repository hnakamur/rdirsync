package rdirsync

import (
	"io"
	"log"
	"os"

	"bitbucket.org/hnakamur/rdirsync/rpc"
)

type server struct{}

func NewServer() rpc.RDirSyncServer {
	return new(server)
}

func (s *server) FetchFile(req *rpc.FetchRequest, stream rpc.RDirSync_FetchFileServer) error {
	log.Printf("FetchFile start. path=%s, bufSize=%d", req.Path, req.BufSize)
	defer log.Printf("FetchFile exit.")
	buf := make([]byte, req.BufSize)
	file, err := os.Open(req.Path)
	if err != nil {
		log.Printf("failed to open file; err=%+v", err)
		return err
	}
	defer file.Close()

	for {
		n, err := io.ReadFull(file, buf)
		log.Printf("after ReadFull. n=%d, err=%+v", n, err)
		if err == io.EOF {
			break
		} else if err == io.ErrUnexpectedEOF {
			buf = buf[:n]
		} else if err != nil {
			return err
		}

		err = stream.Send(&rpc.FileChunk{Chunk: buf})
		log.Printf("after stream.Send. err=%+v", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ReadDir(req *rpc.ReadDirRequest, stream rpc.RDirSync_ReadDirServer) error {
	log.Printf("ReadDir start. path=%s, atMostCount=%d", req.Path, req.AtMostCount)
	defer log.Printf("ReadDir exit.")
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
