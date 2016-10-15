package main

import (
	"context"
	"flag"
	"log"

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
	var remotePath string
	flag.StringVar(&remotePath, "remote-path", "/home/hnakamur/gocode/src/bitbucket.org/hnakamur/rdirsync/rpc/rdirsync.proto", "file path to fetch")
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
	client := rdirsync.NewClientFacade(conn, nil)
	ctx := context.Background()
	switch command {
	case "fetch":
		err := client.FetchFile(ctx, remotePath, localPath)
		if err != nil {
			log.Fatalf("failed to fetch file; %s", err)
		}
	case "readdir":
		infos, err := client.ReadDir(ctx, remotePath)
		if err != nil {
			log.Fatalf("failed to read directory; %s", err)
		}
		for _, info := range infos {
			log.Printf("info=%+v", info)
		}
	default:
		log.Fatalf("Unsupported command: %s", command)
	}
}
