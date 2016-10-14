package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"bitbucket.org/hnakamur/rdirsync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var enableTLS bool
	flag.BoolVar(&enableTLS, "enable-tls", false, "enable TLS")
	var caFile string
	flag.StringVar(&caFile, "ca-file", "../../ssl/ca/cacert.pem", "The file containning the CA root cert file")
	var serverHostOverride string
	flag.StringVar(&serverHostOverride, "server-host-override", "grpc.example.com", "The server name use to verify the hostname returned by TLS handshake")
	var serverAddr string
	flag.StringVar(&serverAddr, "server-addr", "127.0.0.1:10000", "server listen address")
	var command string
	flag.StringVar(&command, "command", "fetch", "operation: one of fetch, readdir")
	var path string
	flag.StringVar(&path, "path", "/home/hnakamur/gocode/src/bitbucket.org/hnakamur/rdirsync/rdirsync.proto", "file path to fetch")
	var localPath string
	flag.StringVar(&localPath, "local-path", "rdirsync.proto", "file path to save")
	var atMostCount int
	flag.IntVar(&atMostCount, "at-most-count", 16, "at most file info count")
	flag.Parse()

	var opts []grpc.DialOption
	if enableTLS {
		var sn string
		if serverHostOverride != "" {
			sn = serverHostOverride
		}
		var creds credentials.TransportCredentials
		if caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(caFile, sn)
			if err != nil {
				log.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := rdirsync.NewRDirSyncClient(conn)
	ctx := context.Background()
	switch command {
	case "fetch":
		stream, err := client.FetchFile(ctx, &rdirsync.FetchRequest{
			Path:    path,
			BufSize: 64,
		})
		if err != nil {
			log.Fatalf("failed to fetch file; %s", err)
		}
		file, err := os.Create(localPath)
		if err != nil {
			log.Fatalf("failed to create file; %s", err)
		}
		defer file.Close()
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			log.Printf("len(chunk.Chunk)=%d", len(chunk.Chunk))
			_, err = file.Write(chunk.Chunk)
			if err != nil {
				log.Fatalf("failed to write file; %s", err)
			}
		}
	case "readdir":
		stream, err := client.ReadDir(ctx, &rdirsync.ReadDirRequest{
			Path:        path,
			AtMostCount: int32(atMostCount),
		})
		if err != nil {
			log.Fatalf("failed to read directory; %s", err)
		}
		for {
			infos, err := stream.Recv()
			if err == io.EOF {
				break
			}
			for _, info := range infos.Infos {
				log.Printf("info=%+v", info)
			}
		}
	default:
		log.Fatalf("Unsupported command: %s", command)
	}
}
