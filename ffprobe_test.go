package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

// newFFProbe builds an FFProbe from a JSON string, mirroring how
// FFProbe.Identify() populates it from ffprobe's output.
func newFFProbe(t *testing.T, jsonStr string) FFProbe {
	t.Helper()

	var ffprobe FFProbe
	if err := json.Unmarshal([]byte(jsonStr), &ffprobe); err != nil {
		t.Fatalf("failed to unmarshal test fixture: %v", err)
	}

	return ffprobe
}

const sampleStreamsJSON = `{
	"streams": [
		{"index": 0, "codec_name": "h264", "codec_type": "video", "tags": {"language": "und"}},
		{"index": 1, "codec_name": "aac", "codec_type": "audio", "tags": {"language": "eng"}},
		{"index": 2, "codec_name": "aac", "codec_type": "audio", "tags": {"language": "swe"}},
		{"index": 3, "codec_name": "ac3", "codec_type": "audio", "tags": {"language": "jpn"}},
		{"index": 4, "codec_name": "subrip", "codec_type": "subtitle", "tags": {"language": "eng"}},
		{"index": 5, "codec_name": "subrip", "codec_type": "subtitle", "tags": {"language": "ger"}}
	]
}`

func TestGetVideoTracks(t *testing.T) {
	ffprobe := newFFProbe(t, sampleStreamsJSON)

	got := ffprobe.GetVideoTracks()
	want := []int{0}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestGetVideoTracks_multiple(t *testing.T) {
	jsonStr := `{
		"streams": [
			{"index": 0, "codec_name": "h264", "codec_type": "video", "tags": {"language": "und"}},
			{"index": 1, "codec_name": "mjpeg", "codec_type": "video", "tags": {"language": "und"}},
			{"index": 2, "codec_name": "aac", "codec_type": "audio", "tags": {"language": "eng"}}
		]
	}`
	ffprobe := newFFProbe(t, jsonStr)

	got := ffprobe.GetVideoTracks()
	want := []int{0, 1}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestGetVideoTracks_none(t *testing.T) {
	jsonStr := `{"streams": [{"index": 0, "codec_name": "aac", "codec_type": "audio", "tags": {"language": "eng"}}]}`
	ffprobe := newFFProbe(t, jsonStr)

	got := ffprobe.GetVideoTracks()

	if len(got) != 0 {
		t.Errorf("expected no video tracks, got %v", got)
	}
}

func TestGetAudioTracks(t *testing.T) {
	tests := []struct {
		name      string
		languages []string
		want      []int
	}{
		{
			name:      "single matching language",
			languages: []string{"eng"},
			want:      []int{1},
		},
		{
			name:      "multiple matching languages",
			languages: []string{"eng", "swe"},
			want:      []int{1, 2},
		},
		{
			name:      "no matching language falls back to keeping all audio tracks",
			languages: []string{"fra"},
			want:      []int{1, 2, 3},
		},
		{
			name:      "empty wanted list falls back to keeping all audio tracks",
			languages: []string{},
			want:      []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ffprobe := newFFProbe(t, sampleStreamsJSON)

			got := ffprobe.GetAudioTracks(tt.languages)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestGetSubtitleTracks(t *testing.T) {
	tests := []struct {
		name      string
		languages []string
		want      []int
	}{
		{
			name:      "single matching language",
			languages: []string{"eng"},
			want:      []int{4},
		},
		{
			name:      "multiple matching languages",
			languages: []string{"eng", "ger"},
			want:      []int{4, 5},
		},
		{
			name:      "no matching language removes all subtitle tracks",
			languages: []string{"fra"},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ffprobe := newFFProbe(t, sampleStreamsJSON)

			got := ffprobe.GetSubtitleTracks(tt.languages)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestGetTracksStatus(t *testing.T) {
	ffprobe := newFFProbe(t, sampleStreamsJSON)

	got := ffprobe.GetTracksStatus([]int{0, 1, 4})

	want := []FFProbeTrack{
		{Index: "0", CodecName: "h264", CodecType: "video", Language: "und", Keep: "✅"},
		{Index: "1", CodecName: "aac", CodecType: "audio", Language: "eng", Keep: "✅"},
		{Index: "2", CodecName: "aac", CodecType: "audio", Language: "swe", Keep: "❌"},
		{Index: "3", CodecName: "ac3", CodecType: "audio", Language: "jpn", Keep: "❌"},
		{Index: "4", CodecName: "subrip", CodecType: "subtitle", Language: "eng", Keep: "✅"},
		{Index: "5", CodecName: "subrip", CodecType: "subtitle", Language: "ger", Keep: "❌"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestNeedsProcessing(t *testing.T) {
	tests := []struct {
		name         string
		tracksToKeep []int
		want         bool
	}{
		{
			name:         "keeping every track means no processing is needed",
			tracksToKeep: []int{0, 1, 2, 3, 4, 5},
			want:         false,
		},
		{
			name:         "dropping a track means processing is needed",
			tracksToKeep: []int{0, 1, 2, 3, 4},
			want:         true,
		},
		{
			name:         "keeping no tracks still means processing is needed",
			tracksToKeep: []int{},
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ffprobe := newFFProbe(t, sampleStreamsJSON)

			got := ffprobe.NeedsProcessing(tt.tracksToKeep)

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
