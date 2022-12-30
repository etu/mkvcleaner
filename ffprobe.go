package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
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

type FFProbeTrack struct {
	Index     string
	Language  string
	CodecName string
	CodecType string
	Keep      string
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

func (ffprobe *FFProbe) GetVideoTracks() []int {
	var videoTracks []int

	// Go through the streams and append the index of all video tracks
	// to the videoTracks slice.
	for _, stream := range ffprobe.Streams {
		if stream.CodecType == "video" {
			videoTracks = append(videoTracks, stream.Index)
		}
	}

	return videoTracks
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

func (ffprobe *FFProbe) GetSubtitleTracks(languages []string) []int {
	var subtitleTracks []int

	for _, stream := range ffprobe.Streams {
		for _, language := range languages {
			if stream.CodecType == "subtitle" && strings.Index(stream.Tags.Language, language) != -1 {
				subtitleTracks = append(subtitleTracks, stream.Index)
			}
		}
	}

	return subtitleTracks
}

func (ffprobe *FFProbe) GetTracksStatus(tracksToKeep []int) []FFProbeTrack {
	var tracks []FFProbeTrack

	for _, track := range ffprobe.Streams {
		keep := "❌"

		for _, trackToKeep := range tracksToKeep {
			if track.Index == trackToKeep {
				keep = "✅"
				break
			}
		}

		tracks = append(tracks, FFProbeTrack{
			Index:     strconv.Itoa(track.Index),
			CodecName: track.CodecName,
			CodecType: track.CodecType,
			Language:  track.Tags.Language,
			Keep:      keep,
		})
	}

	return tracks
}

func (ffprobe *FFProbe) NeedsProcessing(tracksToKeep []int) bool {
	var tracksToDelete []int

	for _, track := range ffprobe.Streams {
		shouldBeKept := false

		for _, trackToKeep := range tracksToKeep {
			if track.Index == trackToKeep {
				shouldBeKept = true
			}
		}

		if !shouldBeKept {
			tracksToDelete = append(tracksToDelete, track.Index)
		}
	}

	return len(tracksToDelete) > 0
}
