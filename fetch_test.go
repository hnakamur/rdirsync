package rdirsync_test

import (
	"context"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"

	"github.com/hnakamur/rdirsync"
	"github.com/hnakamur/rdirsync/internal"
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

	// NOTE: We dial each time for Fetch to avoid the following error.
	// transport: http2Server.HandleStreams failed to read frame: read tcp 127.0.0.1:37025->127.0.0.1:38422: use of closed network connection
	runFetch := func() error {
		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
		if err != nil {
			t.Error(err)
			return err
		}
		defer conn.Close()

		client, err := rdirsync.NewClient(conn,
			rdirsync.SetSyncModTime(true))
		if err != nil {
			t.Error(err)
			return err
		}

		ctx := context.Background()
		err = client.Fetch(ctx, srcDir, destDir)
		if err != nil {
			t.Errorf("%+v", err)
			return err
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
			tree: testFileTreeNode{name: "file1", mode: 0600, size: 0},
			modifications: []modificationOp{
				writeRandomOp(buildSrcPath("file1"), 0, 2000),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dirorfile1", mode: os.ModeDir | 0775,
			},
			modifications: []modificationOp{
				removeDirOp(buildSrcPath("dirorfile1")),
				writeRandomOp(buildSrcPath("dirorfile1"), 0, 0),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dirorfile1", mode: os.ModeDir | 0775,
				children: []testFileTreeNode{
					{name: "dir2", mode: os.ModeDir | 0750},
					{name: "file2-1", mode: 0660, size: 64},
				},
			},
			modifications: []modificationOp{
				removeDirOp(buildSrcPath("dirorfile1")),
				writeRandomOp(buildSrcPath("dirorfile1"), 0, 0),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dirorfile1", mode: 0644, size: 0,
			},
			modifications: []modificationOp{
				removeFileOp(buildSrcPath("dirorfile1")),
				makeDirOp(buildSrcPath("dirorfile1"), 0700),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dir1", mode: os.ModeDir | 0700,
				children: []testFileTreeNode{
					{name: "file1-1", mode: 0600, size: 64},
				},
			},
			modifications: []modificationOp{
				chmodOp(buildSrcPath("dir1"), 0500),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dir1", mode: os.ModeDir | 0700,
				children: []testFileTreeNode{
					{name: "file1-1", mode: 0400, size: 64},
				},
			},
			modifications: []modificationOp{
				chmodOp(buildSrcPath("dir1", "file1-1"), 0200),
				removeFileOp(buildSrcPath("dir1", "file1-1")),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dir1", mode: os.ModeDir | 0700,
				children: []testFileTreeNode{
					{name: "file1-1", mode: 0400, size: 64},
				},
			},
			modifications: []modificationOp{
				chmodOp(buildSrcPath("dir1", "file1-1"), 0600),
				writeRandomOp(buildSrcPath("dir1", "file1-1"), 0, 128),
				chmodOp(buildSrcPath("dir1", "file1-1"), 0400),
			},
		},
		{
			tree: testFileTreeNode{
				name: "dir1", mode: os.ModeDir | 0700,
				children: []testFileTreeNode{
					{
						name: "dir2", mode: os.ModeDir | 0700,
						children: []testFileTreeNode{
							{name: "file3-1", mode: 0400, size: 64},
						},
					},
					{name: "file1-1", mode: 0400, size: 48},
				},
			},
			modifications: []modificationOp{
				removeDirOp(buildSrcPath("dir1", "dir2")),
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

		err = runFetch()
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

		err = runFetch()
		if err != nil {
			return
		}
		sameDirTreeContent(t, destDir, srcDir)

		err = internal.EnsureDirNotExist(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		err = internal.EnsureDirNotExist(destDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type modificationOp func() error

func truncateOp(path string, size int64) modificationOp {
	return func() error {
		err := os.Truncate(path, size)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func writeRandomOp(path string, offset, length int64) modificationOp {
	return func() error {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return errors.WithStack(err)
		}
		defer file.Close()

		_, err = file.Seek(offset, os.SEEK_SET)
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = io.CopyN(file, rand.Reader, length)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func removeFileOp(path string) modificationOp {
	return func() error {
		err := os.Remove(path)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func removeDirOp(path string) modificationOp {
	return func() error {
		err := os.RemoveAll(path)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func makeDirOp(path string, mode os.FileMode) modificationOp {
	return func() error {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func chmodOp(path string, mode os.FileMode) modificationOp {
	return func() error {
		return os.Chmod(path, mode)
	}
}
