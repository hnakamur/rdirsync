package rdirsync

import "os"

func ensureDirExists(path string, mode os.FileMode) error {
	lfi, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, mode.Perm())
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !lfi.IsDir() {
		err = os.Remove(path)
		if err != nil {
			return err
		}
		err = os.MkdirAll(path, mode.Perm())
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureNotExist(path string, fi os.FileInfo) error {
	var err error
	if fi == nil {
		fi, err = os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return err
		}
	}
	if fi.IsDir() {
		return os.RemoveAll(path)
	} else {
		return os.Remove(path)
	}
}

func ensureNotDir(path string, fi os.FileInfo) error {
	var err error
	if fi == nil {
		fi, err = os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return err
		}
	}
	if fi.IsDir() {
		return os.RemoveAll(path)
	}
	return nil
}

func readLocalDir(path string) ([]os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	infos, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}
	sortFileInfosByName(infos)
	return infos, nil
}
