package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FilesOps Methods are nulifying state on error.
type FilesOps struct {
	e         error // TODO: assess maybe better to go with slice
	files     []fs.DirEntry
	filePaths []string
}

func (f *FilesOps) WalkFolder(folder string) *FilesOps {
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

	f.files = res
	f.filePaths = paths

	return f
}

func (f *FilesOps) ByFolder(folder string) *FilesOps {
	if f.e != nil {
		return nil
	}

	res, errRead := os.ReadDir(folder)

	return &FilesOps{
		e:     errRead,
		files: res,
	}
}

func (f *FilesOps) ByExtension(extension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilesOps{
			e: errors.New("no files to search by extension"),
		}
	}

	if extension[:1] != "." {
		extension = "." + extension
	}

	var res []fs.DirEntry

	for _, info := range f.files {
		if extension == filepath.Ext(info.Name()) {
			res = append(res, info)
		}
	}

	f.files = res

	var paths []string

	for _, path := range f.filePaths {
		if extension == filepath.Ext(path) {
			paths = append(paths, path)
		}
	}

	f.filePaths = paths

	return f
}

func (f *FilesOps) ByContent(pattern string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilesOps{
			e: errors.New("no files to search by content"),
		}
	}

	var res []fs.DirEntry

	for _, info := range f.files {
		if errContains := fileContains(pattern, info.Name()); errContains != nil {
			continue
		}

		res = append(res, info)
	}

	f.files = res

	return f
}

func (f *FilesOps) Rename(withExtension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilesOps{
			e: errors.New("no files to move"),
		}
	}

	if withExtension[:1] != "." {
		withExtension = "." + withExtension
	}

	for _, info := range f.files {
		if errMove := os.Rename(info.Name(), info.Name()+withExtension); errMove != nil {
			return &FilesOps{
				e: fmt.Errorf("error when renaming %s", info.Name()),
			}
		}
	}

	return f
}

// Copy Method should be used as a backup.
func (f *FilesOps) Copy(withExtension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilesOps{
			e: errors.New("no files to copy"),
		}
	}

	if withExtension[:1] != "." {
		withExtension = "." + withExtension
	}

	for _, info := range f.files {
		if _, errCopy := fileCopy(info.Name(), info.Name()+withExtension); errCopy != nil {
			return &FilesOps{
				e: fmt.Errorf("error when copying %s", info.Name()),
			}
		}
	}

	return f
}

func (f *FilesOps) PrintFileName(w io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilesOps{
			e: errors.New("no file names to print"),
		}
	}

	for _, info := range f.files {
		w.Write([]byte(info.Name() + "\n"))
	}

	return f
}

func (f *FilesOps) PrintFilePath(w io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files paths to print"),
		}
	}

	for _, info := range f.filePaths {
		w.Write([]byte(info + "\n"))
	}

	return f
}
