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

func TestSend(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatalf("fail to create tempdir; %s", err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")

	buildSrcPath := func(names ...string) string {
		return filepath.Join(srcDir, filepath.Join(names...))
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		grpclog.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	rdirsync.RegisterNewRDirSyncServer(grpcServer)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	runSend := func() error {
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
		err = client.Send(ctx, srcDir, destDir)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	}

	testCases := []struct {
		tree          testFileTreeNode
		modifications []modificationOp
	}{
		{
			tree: testFileTreeNode{
				name: "dir1", mode: os.ModeDir | 0775,
				children: []testFileTreeNode{
					{name: "file1-1", mode: 0666, size: 1024},
					{name: "file1-2", mode: 0660, size: 0},
					{name: "file1-3", mode: 0606, size: 5000 * 1000},
				},
			},
			modifications: []modificationOp{
				truncateOp(buildSrcPath("dir1", "file1-1"), 500),
				removeFileOp(buildSrcPath("dir1", "file1-2")),
				writeRandomOp(buildSrcPath("dir1", "file1-3"), 100, 300),
			},
		},
		{
			tree: testFileTreeNode{name: "file1", mode: 0600, size: 3},
			modifications: []modificationOp{
				truncateOp(buildSrcPath("file1"), 0),
			},
		},
		{
			tree: testFileTreeNode{name: "file1", mode: 0600, size: 1024 * 1024},
			modifications: []modificationOp{
				writeRandomOp(buildSrcPath("file1"), 1024*1024, 1024),
			},
		},
		{
			tree: testFileTreeNode{name: "file1", mode: 0600, size: 16 * 1024 * 1024},
			modifications: []modificationOp{
				writeRandomOp(buildSrcPath("file1"), 16*1024*1024, 1024),
			},
		},
		{
			tree: testFileTreeNode{name: "file1", mode: 0600, size: 0},
			modifications: []modificationOp{
				writeRandomOp(buildSrcPath("file1"), 0, 2000),
			},
		},
	}

	for _, testCase := range testCases {
		err = os.MkdirAll(srcDir, 0700)
		if err != nil {
			log.Fatal(err)
		}

		err = buildFileTree(srcDir, testCase.tree)
		if err != nil {
			log.Fatal(err)
		}

		err = runSend()
		if err != nil {
			return
		}
		sameDirTreeContent(t, destDir, srcDir)

		for _, op := range testCase.modifications {
			err := op()
			if err != nil {
				t.Fatal(err)
			}
		}

		err = runSend()
		if err != nil {
			return
		}
		sameDirTreeContent(t, destDir, srcDir)

		err = os.RemoveAll(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		err = os.RemoveAll(destDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}
