package main

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
)

// PrintContent Method prints the state content.
func (f *FilesOps) PrintContent(w ...io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	// TODO: code repeat. to refactor.
	if len(w) == 0 {
		if len(f.writers) == 0 {
			return &FilesOps{
				e: errors.New("no writers passed and no previous defined writers to write"),
			}
		}

		w = f.writers
	} else {
		f.writers = w
	}

	if len(f.content) == 0 {
		return &FilesOps{
			e: errors.New("no content to print"),
		}
	}

	for _, content := range f.content {
		for _, writer := range w {
			writer.Write([]byte("\n" + content))
		}

	}

	for _, writer := range w {
		writer.Write([]byte("\n"))
	}

	f.writers = w

	return f
}

// PrintFilePath Method prints full path for files in state.
func (f *FilesOps) PrintFilePath(w ...io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(w) == 0 {
		if len(f.writers) == 0 {
			return &FilesOps{
				e: errors.New("no writers passed and no previous defined writers to write"),
			}
		}

		w = f.writers
	} else {
		f.writers = w
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no file names to print"),
		}
	}

	for _, file := range f.filePaths {
		for _, writer := range w {
			writer.Write([]byte(file + "\n"))
		}
	}

	return f
}

// PrintFileNames Method prints name only for files in state.
func (f *FilesOps) PrintFileNames(w ...io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(w) == 0 {
		if len(f.writers) == 0 {
			return &FilesOps{
				e: errors.New("no writers passed and no previous defined writers to write"),
			}
		}

		w = f.writers
	} else {
		f.writers = w
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files names to print"),
		}
	}

	for _, path := range f.filePaths {
		for _, writer := range w {
			writer.Write([]byte(filepath.Base(path) + "\n"))
		}
	}

	return f
}

// PrintFiles Method prints state files to passed writer or the error.
func (f *FilesOps) PrintFilesContent(w ...io.Writer) *FilesOps {
	if len(w) == 0 {
		if len(f.writers) == 0 {
			return &FilesOps{
				e: errors.New("no writers passed and no previous defined writers to write"),
			}
		}

		w = f.writers
	} else {
		f.writers = w
	}

	if f.e != nil {
		for _, writer := range w {
			writer.Write([]byte(f.e.Error()))
		}

		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to print its content\n"),
		}
	}

	for _, file := range f.filePaths {
		content, errRead := ioutil.ReadFile(file)
		if errRead != nil {
			return &FilesOps{
				e: errRead,
			}
		}

		for _, writer := range w {
			writer.Write(content)
			writer.Write([]byte("\n" + "end of file " + file))
			writer.Write([]byte("\n"))
		}
	}

	return f
}
