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
o.SearchByFolder("/home/tudi/ram").PrintFileNames(os.Stdout)
```

### Walk folder
Goes down in the folder structure:
```go
o.SearchWalkFolder("/home/tudi/ram").PrintFileNames(os.Stdout)
```

## Select files
From the search we can filter files.

### Filter by extension
```go
o.SearchByFolder("/home/tudi/ram").FilterByExtension("txt").PrintFileNames(os.Stdout)
```

### Filter by file name
No wildcards.
```go
o.SearchByFolder("/home/tudi/ram").FilterByFileName("/home/tudi/ram/created_file1.txt").PrintFileNames(os.Stdout)
```

### Filter by file contents
Pass pattern, for the case at hand the files are empty thus the empty pattern.
```go
o.SearchByFolder("/home/tudi/ram").FilterByContent("").PrintFileNames(os.Stdout)
```

## File Management
### Rename files
Would add extension to previously selected files.
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout).
		FilesRename("bak").
		SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout)
```

### Revert rename
Would revert or take out passed extension from selected files
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout).
		FilesRevert("bak").
		SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout)
```

### Backup files
Would create backup files with passed extension.
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout).
		FilesCopyToExtension("bak").
		SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout)
```

### Delete files
Would delete all previously selected files.
```go
o.SearchByFolder("/home/tudi/ram").
		PrintFileNames(os.Stdout).
		FilterByExtension("bak").
		FilesDelete().
		PrintFileNames(os.Stdout)
```