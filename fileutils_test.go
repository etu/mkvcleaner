package main

import (
	"os"
	"path/filepath"
	"sort"
	"syscall"
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

func TestCopyFilePermissions(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")

	if err := os.WriteFile(src, []byte("source"), 0o640); err != nil {
		t.Fatalf("failed to create src file: %v", err)
	}
	if err := os.WriteFile(dst, []byte("dest"), 0o600); err != nil {
		t.Fatalf("failed to create dst file: %v", err)
	}

	if err := copyFilePermissions(src, dst); err != nil {
		t.Fatalf("copyFilePermissions returned an error: %v", err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatalf("failed to stat src: %v", err)
	}
	dstInfo, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("failed to stat dst: %v", err)
	}

	if dstInfo.Mode() != srcInfo.Mode() {
		t.Errorf("expected dst mode %v, got %v", srcInfo.Mode(), dstInfo.Mode())
	}

	srcStat := srcInfo.Sys().(*syscall.Stat_t)
	dstStat := dstInfo.Sys().(*syscall.Stat_t)

	if dstStat.Uid != srcStat.Uid || dstStat.Gid != srcStat.Gid {
		t.Errorf("expected dst owner %d:%d, got %d:%d", srcStat.Uid, srcStat.Gid, dstStat.Uid, dstStat.Gid)
	}
}

func TestCopyFilePermissions_missingSource(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "dst")

	if err := os.WriteFile(dst, []byte("dest"), 0o644); err != nil {
		t.Fatalf("failed to create dst file: %v", err)
	}

	err := copyFilePermissions(filepath.Join(dir, "does-not-exist"), dst)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
