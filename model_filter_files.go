package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FilterFiles struct {
	e     error // TODO: assess maybe better to go with slice
	files []fs.FileInfo
}

func (f *FilterFiles) ByFolder(folder string) *FilterFiles {
	if f.e != nil {
		return nil
	}

	res, errRead := ioutil.ReadDir(folder)

	return &FilterFiles{
		e:     errRead,
		files: res,
	}
}

func (f *FilterFiles) ByExtension(extension string) *FilterFiles {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilterFiles{
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

	return &FilterFiles{
		e:     nil,
		files: res,
	}
}

func (f *FilterFiles) ByContent(pattern string) *FilterFiles {
	if f.e != nil {
		return nil
	}

	if len(f.files) == 0 {
		return &FilterFiles{
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

	return &FilterFiles{
		e:     nil,
		files: res,
	}
}

func pathExists(path string) error {
	_, errPath := os.Stat(path)
	if os.IsNotExist(errPath) {
		return errPath
	}

	return nil
}

func fileContains(s, filePath string) error {
	if errPath := pathExists(filePath); errPath != nil {
		return errPath
	}

	content, errRead := ioutil.ReadFile(filePath)
	if errRead != nil {
		return errRead
	}

	if strings.Contains(string(content), s) {
		return nil
	}

	return errors.New("file does not contain string")
}
