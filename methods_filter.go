package main

import (
	"errors"
	"fmt"
	"path/filepath"
)

// ByContent Method would select the files containing passed pattern.
func (f *FilesOps) FilterByContent(pattern string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to search by content"),
		}
	}

	var res []string

	for _, file := range f.filePaths {
		if errContains := fileContains(pattern, file); errContains != nil {
			continue
		}

		res = append(res, file)
	}

	f.filePaths = res

	return f
}

func (f *FilesOps) FilterByExtension(extension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to search by extension"),
		}
	}

	if extension[:1] != "." {
		extension = "." + extension
	}

	var paths []string

	for _, path := range f.filePaths {
		if extension == filepath.Ext(path) {
			paths = append(paths, path)
		}
	}

	f.filePaths = paths

	return f
}

// FilterByFileName Method does not reset content.
func (f *FilesOps) FilterByFileName(filePaths ...string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no passed file names to search"),
		}
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to search by file name"),
		}
	}

	var paths []string

	fmt.Println("filtering by:", filePaths)

	for _, path := range f.filePaths {
		for _, fileName := range filePaths {
			//check works for absolute path also
			if fileName == path || fileName == filepath.Base(path) {
				paths = append(paths, path)
			}
		}
	}

	f.filePaths = paths

	return f
}
