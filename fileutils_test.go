package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFindFilesInDirectory(t *testing.T) {
	root := t.TempDir()

	paths := []string{
		filepath.Join(root, "movie.mkv"),
		filepath.Join(root, "notes.txt"),
		filepath.Join(root, "subdir", "episode1.mkv"),
		filepath.Join(root, "subdir", "episode2.mkv"),
		filepath.Join(root, "subdir", "nested", "extra.MKV"),
		filepath.Join(root, "subdir", "poster.jpg"),
	}

	for _, p := range paths {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("failed to create dir for %s: %v", p, err)
		}
		if err := os.WriteFile(p, []byte("data"), 0o644); err != nil {
			t.Fatalf("failed to create file %s: %v", p, err)
		}
	}

	got := findFilesInDirectory(root)
	sort.Strings(got)

	want := []string{
		filepath.Join(root, "movie.mkv"),
		filepath.Join(root, "subdir", "episode1.mkv"),
		filepath.Join(root, "subdir", "episode2.mkv"),
	}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("expected %v, got %v", want, got)
			break
		}
	}
}

func TestFindFilesInDirectory_empty(t *testing.T) {
	root := t.TempDir()

	got := findFilesInDirectory(root)

	if len(got) != 0 {
		t.Errorf("expected no files, got %v", got)
	}
}

func TestFindFilesInDirectory_nonexistent(t *testing.T) {
	got := findFilesInDirectory(filepath.Join(t.TempDir(), "does-not-exist"))

	if len(got) != 0 {
		t.Errorf("expected no files, got %v", got)
	}
}
