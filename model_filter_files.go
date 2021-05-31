package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type MockInfo struct {
	targetPath    string
	sourceRepo    string
	targetPackage string
	interfaces    []string
}

// FilesOps Methods are nulifying state files on error.
type FilesOps struct {
	spooling  bool
	e         error // only one error for simplicity
	filePaths []string
	content   []string // contains extracted rows
	spool     []string
}

const mockgenBinaryPath = "mockgen"
const mockgenPackageName = "mock"
const makefileName = "Makefile1"

func (f *FilesOps) SearchWalkFolder(folder string) *FilesOps {
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
func (f *FilesOps) SearchByFolder(folder string) *FilesOps {
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
func (f *FilesOps) FilterByFileName(fileName string) *FilesOps {
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

// SelectBy Method would append to state content rows containing passed pattern.
// Only files containing pattern would remain in state.
func (f *FilesOps) ContentExtractByPattern(pattern string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to consider for selection"),
		}
	}

	var paths []string

	for _, path := range f.filePaths {
		fileHandler, errOpen := os.Open(path)
		if errOpen != nil {
			return &FilesOps{
				e: errOpen,
			}
		}
		defer fileHandler.Close()

		scanner := bufio.NewScanner(fileHandler)

		for scanner.Scan() {
			lineContent := scanner.Text()

			if strings.Contains(lineContent, pattern) {
				f.content = append(f.content, lineContent)
				paths = append(paths, path)
			}
		}
	}

	f.filePaths = paths

	return f
}

// ContentAppendMockTargetsMakefile Method would append to Makefile target to re-create mock files.
func (f *FilesOps) ContentAppendMockTargetsMakefile() *FilesOps {
	if f.e != nil {
		return nil
	}

	f.FilterByContent("// Code generated by MockGen. DO NOT EDIT.")
	f.ContentExtractByPattern("// Source:")

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to consider for Makefile mock generation"),
		}
	}

	makefileTarget := []string{".PHONY: gen-mocks", "gen-mocks:"}

	for i := 0; i < len(f.filePaths); i++ {
		content := strings.Split(f.content[i], " ")
		e := createMockTargetEntry(f.filePaths[i], mockgenPackageName, content[2], content[4][:len(content[4])-2])
		makefileTarget = append(makefileTarget, e)
	}

	f.SearchByFolder(f.filePaths[0][:strings.LastIndex(f.filePaths[0], "/")-1])
	f.FilterByFileName(makefileName).ContentAppend(strings.Join(makefileTarget, "\n"))

	return f
}

func createMockTargetEntry(destinationFile, packageName, repo, interfaceName string) string {
	return strings.Join([]string{"@" + mockgenBinaryPath, "-destination=" + destinationFile, "-package=" + packageName, repo, interfaceName}, " ")
}

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

// Delete Method should delete state files.
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
func (f *FilesOps) ContentReplace(old, new string) *FilesOps {
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

// Append Method would append text to state files.
func (f *FilesOps) ContentAppend(text string) *FilesOps {
	if f.e != nil {
		return nil
	}

	if len(f.filePaths) == 0 {
		return &FilesOps{
			e: errors.New("no files to consider for append"),
		}
	}

	for _, file := range f.filePaths {
		content, errRead := os.ReadFile(file)
		if errRead != nil {
			return &FilesOps{
				e: errRead,
			}
		}

		data := string(content) + "\n" + text

		errWrite := ioutil.WriteFile(file, []byte(data), 0644)
		if errWrite != nil {
			return &FilesOps{
				e: errWrite,
			}
		}
	}

	return f
}

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

func (f *FilesOps) SpoolOn() *FilesOps {
	if f.e != nil {
		return nil
	}

	f.spooling = true

	return f
}

func (f *FilesOps) SpoolOff() *FilesOps {
	if f.e != nil {
		return nil
	}

	f.spooling = false

	return f
}
