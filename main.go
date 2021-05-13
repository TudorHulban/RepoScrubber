package main

import (
	"os"
)

func main() {
	folderPath := "/home/tudi/ram/test"

	filter := FilesOps{}
	// files := filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Rename("bak")

	// filter.ByFolder(folderPath).PrintFileNames(os.Stdout)

	filter.WalkFolder(folderPath).ByExtension("go").Copy("bak").PrintFileNames(os.Stdout)
}
