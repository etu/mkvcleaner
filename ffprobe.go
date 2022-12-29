package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type FFProbe struct {
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

func (ffprobe *FFProbe) Identify(fileName string) error {
	// Set up the command to run `ffprobe` with the desired flags.
	cmd := exec.Command(
		"ffprobe",
		"-hide_banner",
		"-show_entries", "stream",
		"-print_format", "json",
		"-v", "panic",
		"-i", fileName,
	)

	// Capture the command's output.
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If there was an error running the command, return the error.
		return fmt.Errorf("Error running ffprobe: %v", err)
	}

	// Check the command's exit code.
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("ffprobe returned a non-zero exit code: %v", cmd.ProcessState.ExitCode())
	}

	// Parse the JSON output.
	if err := json.Unmarshal(output, ffprobe); err != nil {
		return fmt.Errorf("Error parsing JSON output from ffprobe: %v", err)
	}

	return nil
}

func (ffprobe *FFProbe) GetAudioTracks(languages []string) []int {
	var audioTracks []int

	for _, stream := range ffprobe.Streams {
		for _, language := range languages {
			if stream.CodecType == "audio" && strings.Index(stream.Tags.Language, language) != -1 {
				audioTracks = append(audioTracks, stream.Index)
			}
		}
	}

	if len(audioTracks) == 0 {
		// If no tracks were found that match the languages provided,
		// append the index of all audio tracks to the audioTracks slice.
		for _, stream := range ffprobe.Streams {
			if stream.CodecType == "audio" {
				audioTracks = append(audioTracks, stream.Index)
			}
		}
	}

	return audioTracks
}
