package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FilesOps Methods are nulifying state on error.
type FilesOps struct {
	e     error // TODO: assess maybe better to go with slice
	files []fs.FileInfo
}

func (f *FilesOps) ByFolder(folder string) *FilesOps {
	if f.e != nil {
		return nil
	}

	res, errRead := ioutil.ReadDir(folder)

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

	var res []fs.FileInfo

	for _, info := range f.files {
		if extension == filepath.Ext(info.Name()) {
			fmt.Println(info.Name())
			res = append(res, info)
		}
	}

	f.files = res

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

	var res []fs.FileInfo

	for _, info := range f.files {
		if errContains := fileContains(pattern, info.Name()); errContains != nil {
			continue
		}

		fmt.Println(info.Name())
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

func (f *FilesOps) PrintFileNames(w io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilesOps{
			e: errors.New("no files to print"),
		}
	}

	for _, info := range f.files {
		w.Write([]byte(info.Name() + "\n"))
	}

	return f
}
