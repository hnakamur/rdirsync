package rdirsync

import (
	"io"
	"log"
	"os"
)

type server struct{}

func NewServer() *server {
	return new(server)
}

func (s *server) FetchFile(req *FetchRequest, stream RDirSync_FetchFileServer) error {
	log.Printf("FetchFile start. path=%s, bufSize=%d", req.Path, req.BufSize)
	defer log.Printf("FetchFile exit.")
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

		err = stream.Send(&FileChunk{Chunk: buf})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ReadDir(req *ReadDirRequest, stream RDirSync_ReadDirServer) error {
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
		err = stream.Send(&FileInfos{Infos: infos})
		if err != nil {
			return err
		}
	}

	return nil
}

func newFileInfosFromOS(fis []os.FileInfo) []*FileInfo {
	infos := make([]*FileInfo, 0, len(fis))
	for _, fi := range fis {
		infos = append(infos, newFileInfoFromOS(fi))
	}
	return infos
}

func newFileInfoFromOS(fi os.FileInfo) *FileInfo {
	return &FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    int32(fi.Mode()),
		ModTime: fi.ModTime().Unix(),
		IsDir:   fi.IsDir(),
	}
}
