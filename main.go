package main

import (
	"os"
)

func main() {
	folderPath := "/home/tudi/ram/test"

	filter := FilesOps{}

	// filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Rename("bak")

	// filter.ByFolder(folderPath).PrintFileNames(os.Stdout)

	// filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Copy("bak").ByFolder(folderPath).PrintFiles(os.Stdout).ByFileName("main.go.bak").ByExtension("bak").PrintFiles(os.Stdout)

	// filter.ByFolder(folderPath).ByFileName("main.go").Rename("bak").PrintFilePath(os.Stdout)

	filter.SearchByFolder(folderPath).FilterByExtension("bak").PrintFilePath(os.Stdout).FilesRevert("bak")
}
