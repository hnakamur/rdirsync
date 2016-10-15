package rdirsync_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileFailsWhenDirExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filename := filepath.Join(tempDir, "file.txt")
	err = os.Mkdir(filename, 0777)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(filename)
	if err == nil {
		t.Error("should got an error when creating a file where a directory exists")
	}
	defer file.Close()
}

func TestMkdirFailsWhenFileExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filename := filepath.Join(tempDir, "file.txt")

	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	err = os.Mkdir(filename, 0777)
	if err == nil {
		t.Error("should got an error when creating a directory where a file exists")
	}
}

func TestStatFailsWhenFileNotExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rdirsync_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filename := filepath.Join(tempDir, "file.txt")

	_, err = os.Stat(filename)
	if !os.IsNotExist(err) {
		t.Error("should got an is not exist error when running stat for non existent file")
	}
}
