package main

import (
	"io"
)

type MockInfo struct {
	targetPath    string
	sourceRepo    string
	targetPackage string
	interfaces    []string
}

// FilesOps Methods are nulifying state files on error.
type FilesOps struct {
	writers []io.Writer // assess if more than one needed
	e       error       // only one error for simplicity

	filePaths []string
	content   []string // contains extracted rows

	spool    []string
	spooling bool

	rootFolder string // used in Makefile generation
}

const mockgenBinaryPath = "mockgen"
const mockgenPackageName = "mock"
const makefileName = "Makefile1"
