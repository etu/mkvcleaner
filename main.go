package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

func runFFProbe(filename string) {
	// Set up the command to run `ffprobe` with the desired flags.
	cmd := exec.Command(
		"ffprobe",
		"-hide_banner",
		"-show_entries",
		"stream",
		"-print_format",
		"json",
		"-v",
		"panic",
		"-i",
		filename,
	)

	// Capture the command's output.
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If there was an error running the command, print the error message.
		log.Fatal("Error running ffprobe: ", err)
	}

	// Check the command's exit code.
	if cmd.ProcessState.ExitCode() != 0 {
		log.Fatal("ffprobe returned a non-zero exit code: ", cmd.ProcessState.ExitCode())
	}

	// Print the command's output.
	fmt.Printf("Output from ffprobe: %s\n", string(output))
	fmt.Println("ffprobe succeeded")
}

func main() {
	// Loop over the arguments and check if they are valid file or directory names.
	for _, arg := range os.Args[1:] {
		info, err := os.Stat(arg)
		if err == nil {
			if info.IsDir() {
				// If the argument is a directory, call the findFilesInDirectory function.
				files := findFilesInDirectory(arg)

				for _, file := range files {
					// Send each file to the runFFProbe function.
					runFFProbe(file)
				}
			} else {
				// If the argument is a file, send it to the runFFProbe function.
				runFFProbe(arg)
			}
		}
	}
}
