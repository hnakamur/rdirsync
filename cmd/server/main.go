package main

import (
	"flag"
	"net"

	"bitbucket.org/hnakamur/rdirsync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	var enableTLS bool
	flag.BoolVar(&enableTLS, "enable-tls", false, "enable TLS")
	var certFile string
	flag.StringVar(&certFile, "cert-file", "../../ssl/server/server.crt", "TLS cert file")
	var keyFile string
	flag.StringVar(&keyFile, "key-file", "../../ssl/server/server.key", "TLS key file")
	var addr string
	flag.StringVar(&addr, "addr", ":10000", "server listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatal(err)
	}

	var opts []grpc.ServerOption
	if enableTLS {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	rdirsync.RegisterRDirSyncServer(grpcServer, rdirsync.NewServer())
	grpcServer.Serve(lis)
}
