package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

func main() {
	os.Exit(run())
}

type commandOptions struct {
	isServer bool
}

func run() int {
	var options commandOptions
	flag.BoolVar(&options.isServer, "s", false, "enable server mode")
	flag.Parse()

	if options.isServer {
		err := runServer(&options, flag.Args())
		if err != nil {
			log.Printf("%+v", err)
			return 1
		}
	} else {
		err := runClient(&options, flag.Args())
		if err != nil {
			log.Printf("%+v", err)
			return 1
		}
	}
	return 0
}

func runServer(options *commandOptions, args []string) error {
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}
		line = line[:len(line)-1]

		if strings.HasPrefix(line, "sendFile\t") {
			err = processSendFile(line, r)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "sendDir\t") {
			err = processSendDir(line, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processSendFile(line string, r *bufio.Reader) error {
	args := strings.SplitN(line, "\t", 8)
	mode, err := parseFileMode(args[1])
	if err != nil {
		return errors.WithStack(err)
	}
	// atime := args[2]
	// mtime := args[3]
	// owner := args[4]
	// group := args[5]
	size, err := parseSize(args[6])
	if err != nil {
		return errors.WithStack(err)
	}
	path := args[7]

	dir := filepath.Dir(path)
	err = ensureDirExists(dir, 0700)
	if err != nil {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	_, err = io.CopyN(file, r, size)
	if err != nil {
		return errors.WithStack(err)
	}

	err = file.Chmod(mode)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func processSendDir(line string, r *bufio.Reader) error {
	args := strings.SplitN(line, "\t", 8)
	mode, err := parseFileMode(args[1])
	if err != nil {
		return errors.WithStack(err)
	}
	// atime := args[2]
	// mtime := args[3]
	// owner := args[4]
	// group := args[5]
	path := args[6]

	err = ensureDirExists(path, mode)
	if err != nil {
		return err
	}

	err = os.Chmod(path, mode)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func ensureDirExists(path string, mode os.FileMode) error {
	err := os.MkdirAll(path, mode)
	if err != nil {
		perr := err.(*os.PathError)
		if perr.Err != syscall.ENOTDIR {
			return errors.WithStack(err)
		}

		err = os.Remove(perr.Path)
		if err != nil {
			return errors.WithStack(err)
		}

		err = os.MkdirAll(path, mode)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func parseSize(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseFileMode(s string) (os.FileMode, error) {
	var mode os.FileMode
	m, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return mode, err
	}
	return os.FileMode(m), nil
}

func runClient(options *commandOptions, args []string) error {
	if len(args) != 2 {
		return errors.Errorf("source and destination path needed")
	}

	return sendDirOrFile(args[0], args[1])
}

func sendDirOrFile(srcPath, destPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return errors.WithStack(err)
	}

	return sendDirOrFileHelper(srcPath, destPath, srcInfo)
}

func sendDirOrFileHelper(srcPath, destPath string, srcInfo os.FileInfo) error {
	if srcInfo.IsDir() {
		entries, err := readDir(srcPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			err := sendDirOrFileHelper(
				filepath.Join(srcPath, entry.Name()),
				filepath.Join(destPath, entry.Name()),
				entry)
			if err != nil {
				return err
			}
		}

		mode := srcInfo.Mode().Perm()
		fmt.Printf("sendDir\t%o\tatime\tmtime\towner\tgroup\t%s\n", mode, destPath)
	} else {
		err := sendFile(srcPath, destPath, srcInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func readDir(path string) ([]os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	return file.Readdir(0)
}

func sendFile(srcPath, destPath string, srcInfo os.FileInfo) error {
	mode := srcInfo.Mode().Perm()
	size := srcInfo.Size()
	fmt.Printf("sendFile\t%o\tatime\tmtime\towner\tgroup\t%d\t%s\n", mode, size, destPath)
	file, err := os.Open(srcPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
