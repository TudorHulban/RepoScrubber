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

	filter.ByFolder(folderPath).ByExtension("bak").ByFileName("main.go.bak").Replace("package main", "package mainbak").PrintFiles(os.Stdout)
}
