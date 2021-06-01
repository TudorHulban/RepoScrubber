package main

import (
	"fmt"
	"os"
)

func main() {
	folderPath := "/home/tudi/ram/TestMock/"

	filter := FilesOps{}

	// filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Rename("bak")
	// filter.ByFolder(folderPath).PrintFileNames(os.Stdout)
	// filter.ByFolder(folderPath).ByExtension("go").ByContent("package").Copy("bak").ByFolder(folderPath).PrintFiles(os.Stdout).ByFileName("main.go.bak").ByExtension("bak").PrintFiles(os.Stdout)
	// filter.ByFolder(folderPath).ByFileName("main.go").Rename("bak").PrintFilePath(os.Stdout)
	// filter.SearchByFolder(folderPath).FilterByExtension("bak").PrintFilePath(os.Stdout).FilesRevert("bak")

	makefile := []string{folderPath + "/" + makefileName}

	fmt.Println("mk:", makefile)

	// filter.SearchWalkFolder(folderPath).FilesCreate(makefile...).FilterByFileName(makefile...).PrintFileNames(os.Stdout).ContentAppend("xxx").PrintFilesContent(os.Stdout)

	filter.SearchWalkFolder(folderPath).
		FilesCreate(makefile...).
		ContentAppendMockTargetsMakefile(folderPath).
		PrintContent(os.Stdout)
}
