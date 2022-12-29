package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

	// Count all the files
	fileCount := len(fileNames)
	fileNo := 0

	for _, fileName := range fileNames {
		fileNo++

		var ffprobe FFProbe

		ffprobe.Identify(fileName)

		ffmpeg := FFMpeg{
			inputFilePath:  fileName,
			audioTracks:    ffprobe.GetAudioTracks([]string{"eng"}),
			subtitleTracks: ffprobe.GetSubtitleTracks([]string{"eng"}),
			videoTracks:    ffprobe.GetVideoTracks(),
		}

		fmt.Printf("[%d/%d] Preparing to process %s\n", fileNo, fileCount, fileName)
		fmt.Printf("[%d/%d] Command to execute: %s\n", fileNo, fileCount, ffmpeg.FormatCommandParts())
		fmt.Printf("[%d/%d] Run ffmpeg command on %s? [Y/n] ", fileNo, fileCount, fileName)

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)

		if strings.ToLower(response) == "y" || response == "" {
			err := ffmpeg.Run()
			if err == nil {
				fmt.Printf("[%d/%d] Sucessfully cleaned up: %s\n", fileNo, fileCount, fileName)

				// Remove the original file
				err := os.Remove(ffmpeg.inputFilePath)
				if err != nil {
					fmt.Printf("[%d/%d] Failed to remove the input file, bailing out by removing the output file\n", fileNo, fileCount)

					os.Remove(ffmpeg.outputFilePath)
				} else {
					err := os.Rename(ffmpeg.outputFilePath, ffmpeg.inputFilePath)

					if err != nil {
						fmt.Printf("[%d/%d] Failed to rename the temporary file to the original file name", fileNo, fileCount)
					}
				}
			} else {
				fmt.Printf("[%d/%d] Failed at cleaning up: %s\n", fileNo, fileCount, fileName)
				fmt.Printf("[%d/%d] Error: %s\n", fileNo, fileCount, err)
			}
		} else {
			fmt.Printf("[%d/%d] Skipping %s\n", fileNo, fileCount, fileName)
		}
	}
}
