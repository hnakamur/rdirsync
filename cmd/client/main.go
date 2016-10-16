package main

import (
	"context"
	"flag"
	"log"

	"github.com/hnakamur/rdirsync"

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
	flag.StringVar(&command, "command", "fetch", "operation: one of fetch, fetchdir, send, senddir")
	var remotePath string
	flag.StringVar(&remotePath, "remote-path", "/home/hnakamur/gocode/src/github.com/hnakamur/rdirsync/rpc/rdirsync.proto", "file path to fetch")
	var localPath string
	flag.StringVar(&localPath, "local-path", "rdirsync.proto", "file path to save")
	var atMostCount int
	flag.IntVar(&atMostCount, "at-most-count", 16, "at most file info count per rpc")
	var keepDeletedFiles bool
	flag.BoolVar(&keepDeletedFiles, "keep-deleted-files", false, "wether or not keep deleted files")
	var syncModTime bool
	flag.BoolVar(&syncModTime, "sync-mod-time", false, "sync modification time")
	var updateOnly bool
	flag.BoolVar(&updateOnly, "update-only", false, "skip update files whose size is the same and mod time is equal or newer")
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
	client := rdirsync.NewClient(conn,
		rdirsync.SetMaxEntriesPerReadDirRPC(atMostCount),
		rdirsync.SetKeepDeletedFiles(keepDeletedFiles),
		rdirsync.SetSyncModTime(syncModTime),
		rdirsync.SetUpdateOnly(updateOnly),
	)
	ctx := context.Background()
	switch command {
	case "fetch":
		err := client.FetchFile(ctx, remotePath, localPath)
		if err != nil {
			log.Fatalf("failed to fetch file; %s", err)
		}
	case "fetchdir":
		err := client.FetchDir(ctx, remotePath, localPath)
		if err != nil {
			log.Fatalf("failed to fetch directory; %s", err)
		}
	case "send":
		err := client.SendFile(ctx, localPath, remotePath)
		if err != nil {
			log.Fatalf("failed to send file; %s", err)
		}
	case "senddir":
		err := client.SendDir(ctx, localPath, remotePath)
		if err != nil {
			log.Fatalf("failed to send directory; %s", err)
		}
	default:
		log.Fatalf("Unsupported command: %s", command)
	}
}
