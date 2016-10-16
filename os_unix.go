package rdirsync

// +build !windows

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

func makeReadWritable(path string) error {
	err := makeReadWritableParentDir(path)
	if err != nil {
		return err
	}

	return makeReadWritableOneEntry(path)
}

func makeReadWritableRecursive(path string) error {
	err := makeReadWritableParentDir(path)
	if err != nil {
		return err
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return makeReadWritableOneEntry(path)
	})
}

func makeReadWritableParentDir(path string) error {
	dir := filepath.Dir(path)
	fi, err := os.Stat(dir)
	if err != nil && !os.IsPermission(err) {
		return err
	}
	mode := fi.Mode().Perm()

	myUid := uint32(os.Getuid())
	myGid := uint32(os.Getgid())

	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("cannot cast file info to syscall.Stat_t")
	}
	if sys.Uid == myUid {
		mode |= 0700
	} else if sys.Gid == myGid {
		mode |= 0070
	} else {
		mode |= 0007
	}
	if mode != fi.Mode().Perm() {
		err = os.Chmod(dir, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeReadWritableOneEntry(path string) error {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil && !os.IsPermission(err) {
		return err
	}
	mode := fi.Mode().Perm()

	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("cannot cast file info to syscall.Stat_t")
	}

	myUid := uint32(os.Getuid())
	myGid := uint32(os.Getgid())
	if fi.IsDir() {
		if sys.Uid == myUid {
			mode |= 0700
		} else if sys.Gid == myGid {
			mode |= 0070
		} else {
			mode |= 0007
		}
	} else {
		if sys.Uid == myUid {
			mode |= 0600
		} else if sys.Gid == myGid {
			mode |= 0060
		} else {
			mode |= 0006
		}
	}
	if mode != fi.Mode().Perm() {
		err = os.Chmod(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}
