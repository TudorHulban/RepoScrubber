package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FilesOps Methods are nulifying state files on error.
type FilesOps struct {
	e         error // only one error for simplicity
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

	f.filePaths = paths

	return f
}

// ByFolder Method populates state with file paths from provided folder. Does not descend - for that use walk method.
func (f *FilesOps) ByFolder(folder string) *FilesOps {
	if f.e != nil {
		return nil
	}

	files, errRead := os.ReadDir(folder)
	if errRead != nil {
		return &FilesOps{
			e: errRead,
		}
	}

	var res []string

	if folder[len(folder)-1:] != "/" {
		folder = folder + "/"
	}

	for _, file := range files {
		res = append(res, folder+file.Name())
	}

	return &FilesOps{
		filePaths: res,
	}
}

func (f *FilesOps) ByExtension(extension string) *FilesOps {
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

func (f *FilesOps) ByFileName(fileName string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to search by file name"),
		}
	}

	var paths []string

	for _, path := range f.filePaths {
		if fileName == filepath.Base(path) {
			paths = append(paths, path)
		}
	}

	f.filePaths = paths

	return f
}

// ByContent Method would select the files containing passed pattern.
func (f *FilesOps) ByContent(pattern string) *FilesOps {
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

// Rename Method would add passed extension to state files.
func (f *FilesOps) Rename(withExtension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to move"),
		}
	}

	if withExtension[:1] != "." {
		withExtension = "." + withExtension
	}

	for _, file := range f.filePaths {
		if errMove := os.Rename(file, file+withExtension); errMove != nil {
			return &FilesOps{
				e: fmt.Errorf("error when renaming %s", file),
			}
		}
	}

	return f
}

// Copy Method should be used as a backup. Would add passed extension as extension to state files.
func (f *FilesOps) Copy(withExtension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to copy"),
		}
	}

	if withExtension[:1] != "." {
		withExtension = "." + withExtension
	}

	for _, file := range f.filePaths {
		if _, errCopy := fileCopy(file, file+withExtension); errCopy != nil {
			return &FilesOps{
				e: fmt.Errorf("error when copying %s", file),
			}
		}
	}

	return f
}

// Delete Method should delete state files.
func (f *FilesOps) Delete() *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to delete"),
		}
	}

	for _, file := range f.filePaths {
		if errRemove := os.Remove(file); errRemove != nil {
			return &FilesOps{
				e: fmt.Errorf("error when deleting %s", file),
			}
		}
	}

	return f
}

func (f *FilesOps) PrintFileName(w io.Writer) *FilesOps {
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

func (f *FilesOps) PrintFiles(w io.Writer) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to print"),
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

// Replace Method would replace old string with new string searching in state files.
func (f *FilesOps) Replace(old, new string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to consider for replace"),
		}
	}

	for _, file := range f.filePaths {
		content, errRead := os.ReadFile(file)
		if errRead != nil {
			return &FilesOps{
				e: errRead,
			}
		}

		data := strings.ReplaceAll(string(content), old, new)

		errWrite := ioutil.WriteFile(file, []byte(data), 0644)
		if errWrite != nil {
			return &FilesOps{
				e: errWrite,
			}
		}
	}

	return f
}
