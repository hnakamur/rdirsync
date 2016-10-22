package rdirsync_test

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/hnakamur/rdirsync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func TestFetch(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatalf("fail to create tempdir; %s", err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")

	tree := testFileTreeNode{
		name: "dir1", mode: os.ModeDir | 0775,
		children: []testFileTreeNode{
			{name: "file1-1", mode: 0666, size: 1024},
			{name: "file1-2", mode: 0660, size: 0},
			{name: "file1-3", mode: 0606, size: 5000 * 1000},
		},
	}
	err = buildFileTree(srcDir, tree)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		grpclog.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	rdirsync.RegisterNewRDirSyncServer(grpcServer)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client, err := rdirsync.NewClient(conn,
		rdirsync.SetSyncModTime(true))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err = client.Fetch(ctx, srcDir, destDir)
	if err != nil {
		t.Fatal(err)
	}
	sameDirTreeContent(t, destDir, srcDir)
}
