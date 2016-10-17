package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hnakamur/rdirsync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var name = "client"

const globalUsage = `Usage: %s <subcommand> [options] srcDirOrFile destDirOrFile

This is an example client command to synchronize directories and files between
the localhost and a remote server.

subcommands:
    send       send a local directory or file to the remote server.
    fetch      fetch a remote directory or file to the localhost.

`

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Printf(globalUsage, name)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return 1
	}

	switch args[0] {
	case "send":
		return handleSendCommand(args[1:])
	case "fetch":
		return handleFetchCommand(args[1:])
	default:
		flag.Usage()
		return 1
	}
}

type options struct {
	// grpc dial options
	serverAddr         string
	serverHostOverride string
	enableTLS          bool
	caFile             string

	// rdirsync options
	bufSize                 int
	maxEntriesPerReadDirRPC int
	keepDeletedFiles        bool
	syncModTime             bool
	updateOnly              bool
}

func parseOptions(subcommand, usage string, args []string) (*flag.FlagSet, *options) {
	fs := flag.NewFlagSet(subcommand, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(usage, name)
		fs.PrintDefaults()
	}

	var opts options
	fs.StringVar(&opts.serverAddr, "server-addr", "127.0.0.1:10000", "server address to connect")
	fs.StringVar(&opts.serverHostOverride, "server-host-override", "grpc.example.com", "The server name use to verify the hostname returned by TLS handshake")
	fs.BoolVar(&opts.enableTLS, "enable-tls", false, "enable TLS")
	fs.StringVar(&opts.caFile, "ca-file", "cacert.pem", "The file containning the CA root cert file")
	fs.IntVar(&opts.bufSize, "buf-size", 64*1024, "buffer size for reading a file")
	fs.IntVar(&opts.maxEntriesPerReadDirRPC, "at-most-count", 1024, "at most file info count per readdir rpc")
	fs.BoolVar(&opts.keepDeletedFiles, "keep-deleted-files", false, "wether or not keep deleted files")
	fs.BoolVar(&opts.syncModTime, "sync-mod-time", false, "sync modification time")
	fs.BoolVar(&opts.updateOnly, "update-only", false, "skip update files whose size is the same and mod time is equal or newer")
	fs.Parse(args)
	return fs, &opts
}

func (o *options) buildGrpcOptions() []grpc.DialOption {
	var opts []grpc.DialOption
	if o.enableTLS {
		var sn string
		if o.serverHostOverride != "" {
			sn = o.serverHostOverride
		}
		var creds credentials.TransportCredentials
		if o.caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(o.caFile, sn)
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
	return opts
}

func (o *options) buildRDirSyncClientOptions() []rdirsync.ClientOptionFunc {
	return []rdirsync.ClientOptionFunc{
		rdirsync.SetBufSize(o.bufSize),
		rdirsync.SetMaxEntriesPerReadDirRPC(o.maxEntriesPerReadDirRPC),
		rdirsync.SetKeepDeletedFiles(o.keepDeletedFiles),
		rdirsync.SetSyncModTime(o.syncModTime),
		rdirsync.SetUpdateOnly(o.updateOnly),
	}
}

const fetchUsage = `Usage: %s fetch [options] remoteDestDirOrFile localSrcDirOrFile

options:
`

func handleFetchCommand(args []string) int {
	fs, opts := parseOptions("fetch", fetchUsage, args)
	if fs.NArg() != 2 {
		fs.Usage()
		return 1
	}

	remotePath := fs.Arg(0)
	localPath := fs.Arg(1)

	conn, err := grpc.Dial(opts.serverAddr, opts.buildGrpcOptions()...)
	if err != nil {
		fmt.Printf("fail to connect to server: %+v", err)
		return 1
	}
	defer conn.Close()

	client, err := rdirsync.NewClient(conn, opts.buildRDirSyncClientOptions()...)
	if err != nil {
		fmt.Printf("fail to create client: %+v", err)
		return 1
	}

	ctx := context.Background()
	err = client.Fetch(ctx, remotePath, localPath)
	if err != nil {
		fmt.Printf("failed to fetch directory of file; %s", err)
		return 1
	}

	return 0
}

const sendUsage = `Usage: %s send [options] localSrcDirOrFile remoteDestDirOrFile

options:
`

func handleSendCommand(args []string) int {
	fs, opts := parseOptions("send", sendUsage, args)
	if fs.NArg() != 2 {
		fs.Usage()
		return 1
	}

	localPath := fs.Arg(0)
	remotePath := fs.Arg(1)

	conn, err := grpc.Dial(opts.serverAddr, opts.buildGrpcOptions()...)
	if err != nil {
		fmt.Printf("fail to connect to server: %+v", err)
		return 1
	}
	defer conn.Close()

	client, err := rdirsync.NewClient(conn, opts.buildRDirSyncClientOptions()...)
	if err != nil {
		fmt.Printf("fail to create client: %+v", err)
		return 1
	}

	ctx := context.Background()
	err = client.Send(ctx, localPath, remotePath)
	if err != nil {
		fmt.Printf("failed to fetch directory of file; %s", err)
		return 1
	}

	return 0
}
