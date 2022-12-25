package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func findFilesInDirectory(dir string) []string {
	var files []string

	// Use filepath.Walk to traverse the directory tree.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".mkv" {
			// If the file has the correct extension, add it to the list.
			files = append(files, path)
		}
		return nil
	})

	return files
}

func main() {
	var fileNames []string

	// Loop over the arguments and check if they are valid file or directory names.
	for _, arg := range os.Args[1:] {
		info, err := os.Stat(arg)
		if err == nil {
			if info.IsDir() {
				// If the argument is a directory, call the findFilesInDirectory function.
				files := findFilesInDirectory(arg)

				for _, file := range files {
					fileNames = append(fileNames, file)
				}
			} else {
				fileNames = append(fileNames, arg)
			}
		}
	}

	for _, fileName := range fileNames {
		fmt.Println("fileName:", fileName)
		fmt.Print("metaData: ")
		fmt.Println(runFFProbe(fileName))
	}
}
