package main

import (
	"fmt"
	"os"
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

	for _, fileName := range fileNames {
		var ffprobe FFProbe

		ffprobe.Identify(fileName)

		ffmpeg := FFMpeg{
			inputFilePath:  fileName,
			audioTracks:    ffprobe.GetAudioTracks([]string{"eng"}),
			subtitleTracks: ffprobe.GetSubtitleTracks([]string{"eng"}),
			videoTracks:    ffprobe.GetVideoTracks(),
		}

		fmt.Println("ffmpeg struct:", ffmpeg)
		fmt.Println("ffmpeg command:", ffmpeg.FormatCommandParts())
	}
}
