package rdirsync_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
	"time"
)

const (
	bufSize         = 64 * 1024
	testMaxFileSize = int64(1024 * 1024)
)

func TestSimpleCopy(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatalf("fail to create tempdir; %s", err)
	}
	defer os.RemoveAll(tempDir)

	srcFilename := "src.dat"
	srcPath := filepath.Join(tempDir, srcFilename)

	destFilename := "dest.dat"
	destPath := filepath.Join(tempDir, destFilename)

	testCases := []struct {
		origSize     int64
		writeOffset  int64
		writeLen     int64
		truncates    bool
		truncateSize int64
	}{
		{origSize: 1024 * 1024, writeOffset: 512 * 1024, writeLen: 1024 * 1024},
		{origSize: 1024 * 1024, writeOffset: 1024 * 1024, writeLen: 1024 * 1024},
		{origSize: 1024 * 1024, writeOffset: 0, writeLen: 1024 * 1024},
		{origSize: 1024 * 1024, writeOffset: 0, writeLen: 512 * 1024},
		{origSize: 1024 * 1024, truncates: true, truncateSize: 0},
		{origSize: 1024 * 1024, truncates: true, truncateSize: 512 * 1024},
		{origSize: 1024 * 1024, writeOffset: 36, writeLen: 4 * 1024, truncates: true, truncateSize: 10 * 1024},
		{origSize: 8096 * 1024, writeOffset: 8096 * 1024, writeLen: 1024},
	}

	for _, tc := range testCases {
		err = generateRandomFileWithSize(srcPath, tc.origSize)
		if err != nil {
			t.Fatalf("fail to create source file; %s", err)
		}
		err = simpleCopy(srcPath, destPath)
		if err != nil {
			t.Fatalf("failed to do simple copy; %s", err)
		}

		if !sameFileContent(t, tempDir, tempDir, destFilename, srcFilename) {
			t.Errorf("not same file content after simple copy")
		}

		func() {
			srcFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				t.Fatalf("failed to open file to write; %s", err)
			}
			defer srcFile.Close()

			if tc.writeLen > 0 {
				_, err = srcFile.Seek(tc.writeOffset, os.SEEK_SET)
				if err != nil {
					t.Fatalf("failed to seek file to write; %s", err)
				}

				err = writeRandomBytes(srcFile, tc.writeLen)
				if err != nil {
					t.Fatalf("failed to write random bytes to source file; %s", err)
				}
			}

			if tc.truncates {
				err = srcFile.Truncate(tc.truncateSize)
				if err != nil {
					t.Fatalf("failed to truncate source file; %s", err)
				}
			}
		}()

		err = simpleCopy(srcPath, destPath)
		if err != nil {
			t.Fatalf("failed to do compare and copy; %s", err)
		}

		if !sameFileContent(t, tempDir, tempDir, destFilename, srcFilename) {
			t.Errorf("not same file content after compare and copy")
		}
	}
}

func TestCompareAndCopy(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatalf("fail to create tempdir; %s", err)
	}
	defer os.RemoveAll(tempDir)

	srcFilename := "src.dat"
	srcPath := filepath.Join(tempDir, srcFilename)

	destFilename := "dest.dat"
	destPath := filepath.Join(tempDir, destFilename)

	testCases := []struct {
		origSize     int64
		writeOffset  int64
		writeLen     int64
		truncates    bool
		truncateSize int64
	}{
		{origSize: 1024 * 1024, writeOffset: 512 * 1024, writeLen: 1024 * 1024},
		{origSize: 1024 * 1024, writeOffset: 1024 * 1024, writeLen: 1024 * 1024},
		{origSize: 1024 * 1024, writeOffset: 0, writeLen: 1024 * 1024},
		{origSize: 1024 * 1024, writeOffset: 0, writeLen: 512 * 1024},
		{origSize: 1024 * 1024, truncates: true, truncateSize: 0},
		{origSize: 1024 * 1024, truncates: true, truncateSize: 512 * 1024},
		{origSize: 1024 * 1024, writeOffset: 36, writeLen: 4 * 1024, truncates: true, truncateSize: 10 * 1024},
		{origSize: 8096 * 1024, writeOffset: 8096 * 1024, writeLen: 1024},
	}

	for _, tc := range testCases {
		err = generateRandomFileWithSize(srcPath, tc.origSize)
		if err != nil {
			t.Fatalf("fail to create source file; %s", err)
		}
		err = simpleCopy(srcPath, destPath)
		if err != nil {
			t.Fatalf("failed to do simple copy; %s", err)
		}

		if !sameFileContent(t, tempDir, tempDir, destFilename, srcFilename) {
			t.Errorf("not same file content after simple copy")
		}

		func() {
			srcFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				t.Fatalf("failed to open file to write; %s", err)
			}
			defer srcFile.Close()

			if tc.writeLen > 0 {
				_, err = srcFile.Seek(tc.writeOffset, os.SEEK_SET)
				if err != nil {
					t.Fatalf("failed to seek file to write; %s", err)
				}

				err = writeRandomBytes(srcFile, tc.writeLen)
				if err != nil {
					t.Fatalf("failed to write random bytes to source file; %s", err)
				}
			}

			if tc.truncates {
				err = srcFile.Truncate(tc.truncateSize)
				if err != nil {
					t.Fatalf("failed to truncate source file; %s", err)
				}
			}
		}()

		err = compareAndCopy(srcPath, destPath)
		if err != nil {
			t.Fatalf("failed to do compare and copy; %s", err)
		}

		if !sameFileContent(t, tempDir, tempDir, destFilename, srcFilename) {
			t.Errorf("not same file content after compare and copy")
		}
	}
}

func simpleCopy(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buf := make([]byte, bufSize)
	for {
		n, err := io.ReadFull(srcFile, buf)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}

		_, err = destFile.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func compareAndCopy(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	destInfo, err := destFile.Stat()
	if err != nil {
		return err
	}
	destEnd := destInfo.Size()
	if srcInfo.Size() < destInfo.Size() {
		err = destFile.Truncate(srcInfo.Size())
		if err != nil {
			return err
		}
		destEnd = srcInfo.Size()
	}

	var destPos int64
	srcBuf := make([]byte, bufSize)
	destBuf := make([]byte, bufSize)
	for {
		srcN, err := io.ReadFull(srcFile, srcBuf)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}

		if destPos < destEnd {
			destN, err := io.ReadFull(destFile, destBuf)
			if err == io.EOF {
				break
			}
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			destPos += int64(destN)

			if bytes.Equal(destBuf[:destN], srcBuf[:srcN]) {
				continue
			}

			if destN > 0 {
				_, err := destFile.Seek(int64(-destN), os.SEEK_CUR)
				if err != nil {
					return err
				}
			}
		}

		_, err = destFile.Write(srcBuf[:srcN])
		if err != nil {
			return err
		}
	}
	return nil
}

func writeRandomBytes(file *os.File, size int64) error {
	_, err := io.CopyN(file, rand.Reader, size)
	return err
}

func generateRandomFileWithSize(filename string, size int64) error {
	return generateRandomFileWithSizeAndMode(filename, size, 0600)
}

func generateRandomFileWithSizeAndMode(filename string, size int64, mode os.FileMode) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.CopyN(file, rand.Reader, size)
	return err
}

func sameFileContent(t *testing.T, gotDir, wantDir, gotFilename, wantFilename string) bool {
	wantPath := filepath.Join(wantDir, wantFilename)
	wantFile, err := os.Open(wantPath)
	if err != nil {
		t.Fatalf("failed to open file %s; %s", wantPath, err)
	}
	defer wantFile.Close()

	gotPath := filepath.Join(gotDir, gotFilename)
	gotFile, err := os.Open(gotPath)
	if err != nil {
		t.Fatalf("failed to open file %s; %s", gotPath, err)
	}
	defer gotFile.Close()

	wantInfo, err := wantFile.Stat()
	if err != nil {
		t.Fatalf("failed to stat file %s; %s", wantPath, err)
	}
	gotInfo, err := gotFile.Stat()
	if err != nil {
		t.Fatalf("failed to stat file %s; %s", gotPath, err)
	}
	if gotInfo.Size() != wantInfo.Size() {
		t.Errorf("unmatch file size. got:%s, want:%s", gotInfo.Size(), wantInfo.Size())
	}

	wantBuf := make([]byte, 4096)
	gotBuf := make([]byte, 4096)
	for {
		wantN, err := io.ReadFull(wantFile, wantBuf)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			t.Fatalf("failed to read file %s; %s", wantPath, err)
		}
		gotN, err := io.ReadFull(gotFile, gotBuf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			t.Fatalf("failed to read file %s; %s", gotPath, err)
		}
		if !bytes.Equal(gotBuf[:wantN], wantBuf[:gotN]) {
			t.Errorf("unmatch file content. got:%s, want:%s", gotFilename, wantFilename)
			return false
		}
	}
	return true
}

func sameDirTreeContent(t *testing.T, gotDir, wantDir string) bool {
	gotNames, err := filepath.Glob(filepath.Join(gotDir, "*"))
	if err != nil {
		t.Fatalf("failed to glob under dir %s; %s", gotDir, err)
	}
	wantNames, err := filepath.Glob(filepath.Join(wantDir, "*"))
	if err != nil {
		t.Fatalf("failed to glob under dir %s; %s", wantDir, err)
	}
	if len(gotNames) != len(wantNames) {
		t.Errorf("unmatch entry count. got:%d, want:%d", len(gotNames), len(wantNames))
		return false
	}
	sort.Strings(gotNames)
	sort.Strings(wantNames)
	for i := 0; i < len(gotNames); i++ {
		gotName := filepath.Base(gotNames[i])
		wantName := filepath.Base(wantNames[i])
		gotPath := filepath.Join(gotDir, gotName)
		wantPath := filepath.Join(wantDir, wantName)
		gotFileInfo, err := os.Stat(gotPath)
		if err != nil {
			t.Fatalf("cannot stat file; %s", err)
			return false
		}
		wantFileInfo, err := os.Stat(wantPath)
		if err != nil {
			t.Fatalf("cannot stat file; %s", err)
			return false
		}

		same := sameDirOrFile(t, gotDir, wantDir, gotFileInfo, wantFileInfo)
		if !same {
			return false
		}

		if gotFileInfo.IsDir() {
			same = sameDirTreeContent(t, filepath.Join(gotDir, gotName), filepath.Join(wantDir, wantName))
			if !same {
				return false
			}
		}
	}

	return false
}

func sameDirOrFile(t *testing.T, gotDir, wantDir string, gotFileInfo, wantFileInfo os.FileInfo) bool {
	if !sameFileInfo(t, gotDir, wantDir, gotFileInfo, wantFileInfo) {
		return false
	}
	if gotFileInfo.IsDir() {
		return true
	}
	return sameFileContent(t, gotDir, wantDir, gotFileInfo.Name(), wantFileInfo.Name())
}

func sameFileInfoAndContent(t *testing.T, gotDir, wantDir, gotFilename, wantFilename string) bool {
	wantPath := filepath.Join(wantDir, wantFilename)
	wantFileInfo, err := os.Stat(wantPath)
	if err != nil {
		t.Fatalf("fail to stat file %s; %s", wantPath, err)
	}
	gotPath := filepath.Join(gotDir, gotFilename)
	gotFileInfo, err := os.Stat(gotPath)
	if err != nil {
		t.Fatalf("fail to stat file %s; %s", gotPath, err)
	}
	return sameFileInfo(t, gotDir, wantDir, gotFileInfo, wantFileInfo) &&
		sameFileContent(t, gotDir, wantDir, gotFilename, wantFilename)
}

func sameFileInfo(t *testing.T, gotDir, wantDir string, gotFileInfo, wantFileInfo os.FileInfo) bool {
	same := true
	if gotFileInfo.Size() != wantFileInfo.Size() {
		t.Errorf("unmatch size. wantDir:%s; gotFileInfo:%d; wantFileInfo:%d", wantDir, gotFileInfo.Size(), wantFileInfo.Size())
		same = false
	}
	if gotFileInfo.Mode() != wantFileInfo.Mode() {
		t.Errorf("unmatch mode. wantDir:%s; gotFileInfo:%s; wantFileInfo:%s", wantDir, gotFileInfo.Mode(), wantFileInfo.Mode())
		same = false
	}
	gotModTime := gotFileInfo.ModTime()
	wantModTime := wantFileInfo.ModTime()
	if gotModTime != wantModTime {
		t.Errorf("unmatch modification time. wantDir:%s; gotFileInfo:%s; wantFileInfo:%s", wantDir, gotModTime, wantModTime)
		same = false
	}
	if gotFileInfo.IsDir() != wantFileInfo.IsDir() {
		t.Errorf("unmatch isDir. wantDir:%s; gotFileInfo:%v; wantFileInfo:%v", wantDir, gotFileInfo.IsDir(), wantFileInfo.IsDir())
		same = false
	}
	return same
}

type testFileTreeNode struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	owner   string
	group   string

	children []testFileTreeNode
}

func (n testFileTreeNode) IsDir() bool { return n.mode&os.ModeDir != 0 }

func (n testFileTreeNode) Perm() os.FileMode { return n.mode & os.ModePerm }

func buildFileTree(baseDir string, n testFileTreeNode) error {
	path := filepath.Join(baseDir, n.name)
	if n.IsDir() {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}

		for _, child := range n.children {
			err := buildFileTree(filepath.Join(baseDir, n.name), child)
			if err != nil {
				return err
			}
		}
	} else {
		err := generateRandomFileWithSizeAndMode(path, n.size, 0600)
		if err != nil {
			return err
		}
	}

	err := os.Chmod(path, n.Perm())
	if err != nil {
		return err
	}
	if n.owner != "" && n.group != "" && os.Getuid() == 0 {
		u, err := user.Lookup(n.owner)
		if err != nil {
			return err
		}
		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			return err
		}

		g, err := user.LookupGroup(n.group)
		if err != nil {
			return err
		}
		gid, err := strconv.Atoi(g.Gid)
		if err != nil {
			return err
		}
		err = os.Chown(path, uid, gid)
		if err != nil {
			return err
		}
	}
	if !n.modTime.IsZero() {
		err := os.Chtimes(path, time.Now(), n.modTime)
		if err != nil {
			return err
		}
	}
	return nil
}
