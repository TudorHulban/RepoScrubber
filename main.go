package main

import (
	"fmt"
)

func main() {
	folderPath := "/home/tudi/ram/TestMock"

	filter := FilterFiles{}
	files := filter.ByFolder(folderPath).ByExtension("go").ByContent("package")

	fmt.Println(len(files.files), files.e)
}
