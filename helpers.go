package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func pathExists(filePath string) error {
	sourceFileStat, errPath := os.Stat(filePath)
	if os.IsNotExist(errPath) {
		return errPath
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", filePath)
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

func fileCopy(sourcePath, destinationPath string) (int64, error) {
	if errExists := pathExists(sourcePath); errExists != nil {
		return 0, errExists
	}

	source, errOpenSourcePath := os.Open(sourcePath)
	if errOpenSourcePath != nil {
		return 0, errOpenSourcePath
	}
	defer source.Close()

	destination, errOpenDestinationPath := os.Create(destinationPath)
	if errOpenDestinationPath != nil {
		return 0, errOpenDestinationPath
	}
	defer destination.Close()

	return io.Copy(destination, source)
}
