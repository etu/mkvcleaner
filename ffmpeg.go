package main

import (
	"fmt"
	"path/filepath"
)

type FFMpeg struct {
	inputFilePath  string
	outputFilePath string
	audioTracks    []int
	subtitleTracks []int
	videoTracks    []int
}

func (ffmpeg *FFMpeg) FormatCommandParts() []string {
	// Extract the directory and file name from the path
	dir, file := filepath.Split(ffmpeg.inputFilePath)

	// Join the directory and modified file name to get the new file path
	ffmpeg.outputFilePath = filepath.Join(dir, ".tmp."+file)

	// Go through audio tracks to append to ffmpeg command
	audioTracksArgs := []string{}
	for _, track := range ffmpeg.audioTracks {
		audioTracksArgs = append(audioTracksArgs, fmt.Sprintf("-map 0:%d", track))
	}

	// Go through video tracks to append to ffmpeg command
	videoTracksArgs := []string{}
	for _, track := range ffmpeg.videoTracks {
		videoTracksArgs = append(videoTracksArgs, fmt.Sprintf("-map 0:%d", track))
	}

	// Go through subtitle tracks to append to ffmpeg command
	subtitleTracksArgs := []string{}
	for _, track := range ffmpeg.subtitleTracks {
		subtitleTracksArgs = append(subtitleTracksArgs, fmt.Sprintf("-map 0:%d", track))
	}

	// Build and return the ffmpeg command as a slice of strings
	return append(
		[]string{
			"ffmpeg",
			"-i",
			ffmpeg.inputFilePath,
			"-c",
			"copy",
		},
		append(videoTracksArgs, append(audioTracksArgs, append(subtitleTracksArgs, ffmpeg.outputFilePath)...)...)...,
	)
}
