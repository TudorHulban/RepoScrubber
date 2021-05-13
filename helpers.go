package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

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
