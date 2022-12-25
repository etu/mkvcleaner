package main

import (
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
