package rdirsync

import (
	"os"

	"github.com/pkg/errors"
)

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
