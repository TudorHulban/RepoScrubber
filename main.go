package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func pathExists(path string) error {
	_, errPath := os.Stat(path)
	if os.IsNotExist(errPath) {
		return errPath
	}

	return nil
}

func getFolderFileNames(folder string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(folder)
}

func filterFilesForExtension(ext string, fileInfo []fs.FileInfo) ([]fs.FileInfo, error) {
	if len(fileInfo) == 0 {
		return nil, errors.New("no file info passed")
	}

	if ext[:1] != "." {
		ext = "." + ext
	}

	var res []fs.FileInfo

	for _, info := range fileInfo {
		if ext == filepath.Ext(info.Name()) {
			fmt.Println(info.Name())
			res = append(res, info)
		}
	}

	return res, nil
}

func FolderFilesByExtension(folder, extension string) ([]fs.FileInfo, error) {
	// TODO: input validation
	errFolderExists := pathExists(folder)
	if errFolderExists != nil {
		return nil, errFolderExists
	}

	files, errFolder := getFolderFileNames(folder)
	if errFolder != nil {
		return nil, errFolder
	}

	return filterFilesForExtension(extension, files)
}

func main() {
	folderPath := "/home/tudi/ram/TestMock"

	files, err := FolderFilesByExtension(folderPath, "go")
	fmt.Println(len(files), err)
}
