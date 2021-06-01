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


## Add extension to files