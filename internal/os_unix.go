package internal

// +build !windows

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
)

func EnsureDirExists(path string, mode os.FileMode) error {
	fi, err := tryHardStat(path)
	if err == nil && fi.IsDir() {
		return nil
	} else if err != nil && !os.IsNotExist(errors.Cause(err)) {
		return err
	}

	if err == nil && !fi.IsDir() {
		err = EnsureFileNotExist(path)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(path, mode.Perm())
	if err == nil {
		return nil
	} else if os.IsPermission(err) {
		err = makeReadWritableParentDir(path)
		if err != nil {
			return err
		}
		err = os.MkdirAll(path, mode.Perm())
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func EnsureNotDir(path string) error {
	fi, err := tryHardStat(path)
	if os.IsNotExist(errors.Cause(err)) {
		return nil
	} else if err != nil {
		return err
	}

	if fi.IsDir() {
		return EnsureDirNotExist(path)
	}
	return nil
}

func tryHardStat(path string) (os.FileInfo, error) {
	fi, err := os.Stat(path)
	if os.IsPermission(err) {
		err = makeReadWritableParentDir(path)
		if err != nil {
			return nil, err
		}
		fi, err = os.Stat(path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return fi, nil
}

func EnsureDirOrFileNotExist(path string) error {
	fi, err := tryHardStat(path)
	if os.IsNotExist(errors.Cause(err)) {
		return nil
	} else if err != nil {
		return err
	}

	if fi.IsDir() {
		return EnsureDirNotExist(path)
	} else {
		return EnsureFileNotExist(path)
	}
}

func EnsureDirNotExist(path string) error {
	err := os.RemoveAll(path)
	if err == nil {
		return nil
	} else if !os.IsPermission(err) {
		return errors.WithStack(err)
	}

	err = MakeReadWritableRecursive(path)
	if err != nil {
		return err
	}

	err = os.RemoveAll(path)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func EnsureFileNotExist(path string) error {
	err := os.Remove(path)
	if err == nil {
		return nil
	} else if !os.IsPermission(err) {
		return errors.WithStack(err)
	}

	err = makeReadWritableParentDir(path)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func MakeReadWritable(path string) error {
	err := makeReadWritableParentDir(path)
	if err != nil {
		return err
	}

	return makeReadWritableOneEntry(path)
}

func MakeReadWritableRecursive(path string) error {
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
		return errors.WithStack(err)
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
			return errors.WithStack(err)
		}
	}
	return nil
}

func makeReadWritableOneEntry(path string) error {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil && !os.IsPermission(err) {
		return errors.WithStack(err)
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
			return errors.WithStack(err)
		}
	}
	return nil
}
