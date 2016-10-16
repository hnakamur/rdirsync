package rdirsync

import (
	"os"
	"sort"
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

type osFileInfosByName []os.FileInfo

func (a osFileInfosByName) Len() int           { return len(a) }
func (a osFileInfosByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a osFileInfosByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func sortFileInfosByName(infos []os.FileInfo) {
	sort.Sort(osFileInfosByName(infos))
}

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
		err = os.RemoveAll(path)
		if !os.IsPermission(err) {
			return err
		}

		err = makeReadWritableRecursive(path)
		if err != nil {
			return err
		}

		return os.RemoveAll(path)
	} else {
		err = os.Remove(path)
		if !os.IsPermission(err) {
			return err
		}

		err = makeReadWritable(path)
		if err != nil {
			return err
		}
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
