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
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		perr := err.(*os.PathError)
		if perr.Err != syscall.ENOTDIR {
			return errors.WithStack(err)
		}

		err = os.Remove(perr.Path)
		if err != nil {
			return errors.WithStack(err)
		}

		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return errors.WithStack(err)
		}
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

	srcPath := args[0]
	destPath := args[1]

	file, err := os.Open(srcPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	mode := fi.Mode().Perm()
	size := fi.Size()
	fmt.Printf("sendFile\t%o\tatime\tmtime\towner\tgroup\t%d\t%s\n", mode, size, destPath)
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
