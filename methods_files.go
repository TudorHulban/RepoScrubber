package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Copy Method should be used as a backup. Would add passed extension as extension to state files.
func (f *FilesOps) FilesCopyToExtension(extension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to copy"),
		}
	}

	if extension[:1] != "." {
		extension = "." + extension
	}

	for _, file := range f.filePaths {
		if _, errCopy := fileCopy(file, file+extension); errCopy != nil {
			return &FilesOps{
				e: fmt.Errorf("error when copying %s", file),
			}
		}
	}

	return f
}

// Delete Method should delete files in state.
func (f *FilesOps) FilesDelete() *FilesOps {
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

	f.filePaths = []string{} // reset state after we deleted the files

	return f
}

// FilesCreate Method should create passed file paths.
// Does not change state.
func (f *FilesOps) FilesCreate(filePaths ...string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to create"),
		}
	}

	for _, path := range filePaths {
		file, errCreate := os.Create(path)
		if errCreate != nil {
			return &FilesOps{
				e: errCreate,
			}
		}

		file.Close()
	}

	return f
}

// Rename Method would add passed extension to state files.
func (f *FilesOps) FilesRename(withExtension string) *FilesOps {
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

// Revert Method would delete passed extension from state files.
func (f *FilesOps) FilesRevert(extension string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to revert"),
		}
	}

	if extension[:1] != "." {
		extension = "." + extension
	}

	for _, file := range f.filePaths {
		ix := strings.LastIndex(file, extension)

		if errMove := os.Rename(file, file[:ix]); errMove != nil {
			return &FilesOps{
				e: fmt.Errorf("error when reverting %s", file),
			}
		}
	}

	return f
}
