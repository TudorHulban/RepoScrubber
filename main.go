package main

import (
	"fmt"
	"os"
)

func main() {
	folderPath := "/home/tudi/ram/test"

	filter := FilesOps{}

	// filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Rename("bak")

	// filter.ByFolder(folderPath).PrintFileNames(os.Stdout)

	f := filter.ByFolder(folderPath).ByExtension("go").ByContent("package").PrintFilePath(os.Stdout)
	fmt.Println(f.e)
}
