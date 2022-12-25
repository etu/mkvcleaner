package main

import (
	"encoding/json"
	"fmt"
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

type FFProbeOutput struct {
	Streams []struct {
		Index     int    `json:"index"`
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
		Tags      struct {
			Title    string `json:"title"`
			Language string `json:"language"`
		} `json:"tags"`
		Disposition struct {
			Default  int `json:"default"`
			Dub      int `json:"dub"`
			Original int `json:"original"`
		} `json:"disposition"`
	} `json:"streams"`
}

func runFFProbe(filename string) (*FFProbeOutput, error) {
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
		// If there was an error running the command, return the error.
		return nil, fmt.Errorf("Error running ffprobe: %v", err)
	}

	// Check the command's exit code.
	if cmd.ProcessState.ExitCode() != 0 {
		return nil, fmt.Errorf("ffprobe returned a non-zero exit code: %v", cmd.ProcessState.ExitCode())
	}

	// Parse the JSON output.
	var data FFProbeOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("Error parsing JSON output from ffprobe: %v", err)
	}

	return &data, nil
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
		runFFProbe(fileName)
	}
}
