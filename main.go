package main

import (
	"os"
)

func main() {
	folderPath := "/home/tudi/ram/TestMock"

	filter := FilesOps{}
	// files := filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Rename("bak")

	filter.ByFolder(folderPath).PrintFileNames(os.Stdout)
}
