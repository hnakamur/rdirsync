package rdirsync

import (
	"os"

	"github.com/pkg/errors"
)

func selectDirAndRegularFiles(fis []os.FileInfo) []os.FileInfo {
	ret := make([]os.FileInfo, 0, len(fis))
	for _, fi := range fis {
		if fi.IsDir() || fi.Mode().IsRegular() {
			ret = append(ret, fi)
		}
	}
	return ret
}

func ensureDirExists(path string, mode os.FileMode) error {
	lfi, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, mode.Perm())
		if err != nil {
			return errors.WithStack(err)
		}
	} else if err != nil {
		return errors.WithStack(err)
	} else if !lfi.IsDir() {
		err = os.Remove(path)
		if err != nil {
			return errors.WithStack(err)
		}
		err = os.MkdirAll(path, mode.Perm())
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func ensureNotDir(path string, fi os.FileInfo) error {
	var err error
	if fi == nil {
		fi, err = os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return errors.WithStack(err)
		}
	}
	if fi.IsDir() {
		err = os.RemoveAll(path)
		if err != nil {
			return errors.WithStack(err)
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

	infos, err := file.Readdir(0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return infos, nil
}
