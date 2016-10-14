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
