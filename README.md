# RepoScrubber
Created for easier handling files and replacement in files.

## Create reference
Common for all below:
```go
o := FilesOps{}
```

## Create files
Create one file:
```go
o.FilesCreate("/home/tudi/ram/created_file.txt")
```

Create two or more files:
```go
o.FilesCreate("/home/tudi/ram/created_file1.txt", "/home/tudi/ram/created_file2.txt")
```

## Search for files
In order to apply operations to files the file names need to be selected.<br/>

### Search by folder
Does not go down in the folder structure:
```go
o.SearchByFolder("/home/tudi/ram").PrintFileNames([]io.Writer{os.Stdout}...)
```

### Walk folder
Goes down in the folder structure:
```go
o.SearchWalkFolder("/home/tudi/ram").PrintFileNames([]io.Writer{os.Stdout}...)
```

## Select files
From the search we can filter files.

### Filter by extension
```go
o.SearchByFolder("/home/tudi/ram").FilterByExtension("txt").PrintFileNames([]io.Writer{os.Stdout}...)
```

### Filter by file name
No wildcards.
```go
o.SearchByFolder("/home/tudi/ram").FilterByFileName("/home/tudi/ram/created_file1.txt").
	PrintFileNames([]io.Writer{os.Stdout}...)
```

### Filter by file contents
Pass pattern, for the case at hand the files are empty thus the empty pattern.
```go
o.SearchByFolder("/home/tudi/ram").FilterByContent("").PrintFileNames([]io.Writer{os.Stdout}...)
```

## File Management
### Rename files
Would add extension to previously selected files.
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames([]io.Writer{os.Stdout}...).
		FilesRename("bak").
		SearchByFolder("/home/tudi/ram").
		PrintFileNames()
```

### Revert rename
Would revert or take out passed extension from selected files
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames([]io.Writer{os.Stdout}...).
		FilesRevert("bak").
		SearchByFolder("/home/tudi/ram").
		PrintFileNames()
```

### Backup files
Would create backup files with passed extension.
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames([]io.Writer{os.Stdout}...).
		FilesCopyToExtension("bak").
		SearchByFolder("/home/tudi/ram").
		PrintFileNames()
```

### Delete files
Would delete all previously selected files.
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames([]io.Writer{os.Stdout}...).
		FilterByExtension("bak").
		FilesDelete().
		PrintFileNames()
```

## Showing information
The writers can be set only once.
### Printing selected file names
```go
o.SearchByFolder("/home/tudi/ram").PrintFileNames([]io.Writer{os.Stdout}...)
```

### Printing selected file paths
```go
o.SearchByFolder("/home/tudi/ram").PrintFilePath([]io.Writer{os.Stdout}...)
```

### Printing selected file contents
```go
o.SearchByFolder("/home/tudi/ram").PrintFilesContent([]io.Writer{os.Stdout}...)
```

### Printing contents
Prints extracted from files information.
```go
o.SearchByFolder("/home/tudi/ram").PrintContent([]io.Writer{os.Stdout}...)
```

## Content Management
### Append content to previously selected files
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames([]io.Writer{os.Stdout}...).
		FileAppend("y").
		ContentAdd("Content:").
		ContentAddByPattern("y").
		PrintContent()
```
### Extract file lines containing specified content
```go
o.SearchByFolder("/home/tudi/ram").
		FilesCreate("/home/tudi/ram/created_file1.txt", "/home/tudi/ram/created_file2.txt").
		PrintFileNames([]io.Writer{os.Stdout}...).
		FilesAppend("y").
		ContentAdd("Content:").
		ContentAddByPattern("y").
		PrintContent()
``` 
