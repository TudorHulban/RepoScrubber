package main

import (
	"errors"
	"io"
	"io/ioutil"
)

// PrintContent Method prints the state content.
func (f *FilesOps) PrintContent(w io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.content) == 0 {
		return &FilesOps{
			e: errors.New("no content to print"),
		}
	}

	for _, content := range f.content {
		w.Write([]byte("\n" + content))
	}

	w.Write([]byte("\n"))

	return f
}

// PrintFileNames Method prints files state.
func (f *FilesOps) PrintFileNames(w io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no file names to print"),
		}
	}

	for _, file := range f.filePaths {
		w.Write([]byte(file + "\n"))
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

// PrintFiles Method prints state files to passed writer or the error.
func (f *FilesOps) PrintFilesContent(w io.Writer) *FilesOps {
	if f.e != nil {
		w.Write([]byte(f.e.Error()))
		return nil
	}

	if len(f.filePaths) == 0 {
		err := errors.New("no files to print its content\n")

		w.Write([]byte(err.Error()))

		return &FilesOps{
			e: err,
		}
	}

	for _, file := range f.filePaths {
		content, errRead := ioutil.ReadFile(file)
		if errRead != nil {
			return &FilesOps{
				e: errRead,
			}
		}

		w.Write(content)
		w.Write([]byte("\n" + "end of file " + file))
		w.Write([]byte("\n"))
	}

	return f
}
