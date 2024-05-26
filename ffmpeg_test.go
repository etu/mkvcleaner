package main

import (
	"reflect"
	"testing"
)

// TestFormatCommandParts tests the FormatCommandParts method of the FFMpeg struct
func TestFormatCommandParts(t *testing.T) {
	tests := []struct {
		name           string
		inputFilePath  string
		audioTracks    []int
		subtitleTracks []int
		videoTracks    []int
		expectedCmd    []string
	}{
		{
			name:           "single audio, video, and subtitle track, full path",
			inputFilePath:  "/path/to/input/file.mkv",
			audioTracks:    []int{1},
			subtitleTracks: []int{2},
			videoTracks:    []int{0},
			expectedCmd: []string{
				"ffmpeg",
				"-i",
				"/path/to/input/file.mkv",
				"-c",
				"copy",
				"-map 0:0",
				"-map 0:1",
				"-map 0:2",
				"/path/to/input/.tmp.file.mkv",
			},
		},
		{
			name:           "multiple audio and subtitle tracks, full path",
			inputFilePath:  "/path/to/input/file.mkv",
			audioTracks:    []int{1, 2},
			subtitleTracks: []int{3, 4},
			videoTracks:    []int{0},
			expectedCmd: []string{
				"ffmpeg",
				"-i",
				"/path/to/input/file.mkv",
				"-c",
				"copy",
				"-map 0:0",
				"-map 0:1",
				"-map 0:2",
				"-map 0:3",
				"-map 0:4",
				"/path/to/input/.tmp.file.mkv",
			},
		},
		{
			name:           "no audio, subtitle, or video tracks, full path",
			inputFilePath:  "/path/to/input/file.mkv",
			audioTracks:    []int{},
			subtitleTracks: []int{},
			videoTracks:    []int{},
			expectedCmd: []string{
				"ffmpeg",
				"-i",
				"/path/to/input/file.mkv",
				"-c",
				"copy",
				"/path/to/input/.tmp.file.mkv",
			},
		},
		{
			name:           "no audio, subtitle, or video tracks, relative path",
			inputFilePath:  "input/file.mkv",
			audioTracks:    []int{},
			subtitleTracks: []int{},
			videoTracks:    []int{},
			expectedCmd: []string{
				"ffmpeg",
				"-i",
				"input/file.mkv",
				"-c",
				"copy",
				"input/.tmp.file.mkv",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ffmpeg := &FFMpeg{
				inputFilePath:  tt.inputFilePath,
				audioTracks:    tt.audioTracks,
				subtitleTracks: tt.subtitleTracks,
				videoTracks:    tt.videoTracks,
			}

			cmdParts := ffmpeg.FormatCommandParts()

			if !reflect.DeepEqual(cmdParts, tt.expectedCmd) {
				t.Errorf("expected command parts %v, got %v", tt.expectedCmd, cmdParts)
			}
		})
	}
}
