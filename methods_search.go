package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

// ByFolder Method populates state with file paths from provided folder. Does not descend - for that use walk method.
func (f *FilesOps) SearchByFolder(folder string) *FilesOps {
	if f.e != nil {
		return nil
	}

	files, errRead := os.ReadDir(folder)
	if errRead != nil {
		return &FilesOps{
			e: errRead,
		}
	}

	var res []string

	if folder[len(folder)-1:] != "/" {
		folder = folder + "/"
	}

	for _, file := range files {
		res = append(res, folder+file.Name())
	}

	return &FilesOps{
		filePaths: res,
	}
}

func (f *FilesOps) SearchWalkFolder(folder string) *FilesOps {
	if f.e != nil {
		return nil
	}

	var res []fs.DirEntry
	var paths []string

	walkFunction := func(path string, infoCurrent os.DirEntry, errCurrent error) error {
		if errCurrent != nil {
			return errCurrent
		}

		if infoCurrent.IsDir() {
			return nil
		}

		res = append(res, infoCurrent)
		paths = append(paths, path)
		return nil
	}

	errWalk := filepath.WalkDir(folder, walkFunction)
	if errWalk != nil {
		return &FilesOps{
			e: errWalk,
		}
	}

	f.filePaths = paths

	return f
}
