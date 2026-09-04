package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
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

// copyFilePermissions copies the owner, group, and mode of srcPath onto
// dstPath, so a remuxed file ends up with the same permissions as the
// original file it's about to replace.
func copyFilePermissions(srcPath, dstPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", srcPath, err)
	}

	if err := os.Chmod(dstPath, info.Mode()); err != nil {
		return fmt.Errorf("failed to chmod %s: %w", dstPath, err)
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("unable to determine owner and group of %s", srcPath)
	}

	if err := os.Chown(dstPath, int(stat.Uid), int(stat.Gid)); err != nil {
		return fmt.Errorf("failed to chown %s: %w", dstPath, err)
	}

	return nil
}
