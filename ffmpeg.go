package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/alessio/shellescape.v1"
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

	// Set output path
	ffmpeg.outputFilePath = filepath.Join(dir, ".tmp."+file)

	// Escape paths
	escapedInputFilePath := shellescape.Quote(ffmpeg.inputFilePath)
	escapedOutputFilePath := shellescape.Quote(dir + ".tmp." + file)

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
			escapedInputFilePath,
			"-c",
			"copy",
		},
		append(videoTracksArgs, append(audioTracksArgs, append(subtitleTracksArgs, escapedOutputFilePath)...)...)...,
	)
}

func (ffmpeg *FFMpeg) Run() error {
	// Get the formatted command parts
	cmdParts := ffmpeg.FormatCommandParts()

	// Run the ffmpeg command
	cmd := exec.Command("sh", "-c", strings.Join(cmdParts, " "))

	return cmd.Run()
}
