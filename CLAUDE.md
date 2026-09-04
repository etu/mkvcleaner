# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**mkvcleaner** is a Go CLI tool that bulk-remuxes MKV files to remove audio/subtitle tracks of unwanted languages. It uses `ffprobe` to inspect track metadata and `ffmpeg` to remux files, copying only the desired tracks into a new file.

Repository: `github.com/etu/mkvcleaner`

## Commands

```bash
make build   # compiles the binary to ./mkvcleaner (go build -o mkvcleaner)
make test    # runs go test
go test      # equivalent to make test
go test -run TestName   # run a single test
```

Tests live in `ffmpeg_test.go` and cover ffmpeg command construction (`FFMpeg.FormatCommandParts()`) via table-driven tests using the standard `testing` package — no external test frameworks.

### Nix

`flake.nix` provides a default package (builds the Go module, wraps the binary with `ffmpeg` on `PATH`) and a dev shell with `ffmpeg`, `gnumake`, `delve`, `go`, `gopls`.

## Architecture

Single Go package (`package main`), no sub-packages, all source files at repo root.

| File | Purpose |
|---|---|
| `main.go` | Entry point. Parses CLI flags (`--langs`, `--automatic`), resolves file/directory arguments, orchestrates processing with interactive confirmation. |
| `ffprobe.go` | `FFProbe` struct/methods. Shells out to `ffprobe` to identify streams in an MKV file. Filtering methods: `GetVideoTracks()`, `GetAudioTracks(languages)`, `GetSubtitleTracks(languages)`, `GetTracksStatus(tracksToKeep)`, `NeedsProcessing(tracksToKeep)`. |
| `ffmpeg.go` | `FFMpeg` struct/methods. Builds and executes the `ffmpeg` remux command with `-c copy` and `-map` flags for selected tracks. Handles shell escaping of file paths via `shellescape`. |
| `fileutils.go` | `findFilesInDirectory()` — recursively walks a directory tree and returns all `.mkv` file paths. |
| `ffmpeg_test.go` | Unit tests for `FFMpeg.FormatCommandParts()`, covering path formats, special characters, and track combinations. |

### Key data flow

1. CLI args are parsed; directories are expanded to `.mkv` file lists (`fileutils.go`).
2. Each file is identified with `ffprobe` (`ffprobe.go`) to discover all streams.
3. Streams are filtered by language against the wanted languages list.
4. **Safety behavior**: if all audio tracks would be removed, all audio tracks are kept instead.
5. Subtitle tracks matching no wanted language are simply removed — no safety fallback (a file can end up with no subtitles).
6. All video tracks are always kept.
7. A table is printed showing which tracks will be kept/removed.
8. If changes are needed and the user confirms (or `--automatic` is set), `ffmpeg` remuxes to a `.tmp.`-prefixed file, then the original is atomically replaced via rename (input renamed to `.rename-tmp`, output renamed to the input's original name, then the `.rename-tmp` file is removed).

### CLI interface

```
./mkvcleaner [--langs=und,eng,swe,jap,jpn] [--automatic] <path>...
```

- `--langs`: comma-separated wanted language codes (default: `und,eng,swe,jap,jpn`). Keep `und` in the list — many single-track files mark their only track as `undefined`.
- `--automatic`: skip interactive per-file confirmation prompts.
- Positional args: files or directories (directories are scanned recursively for `.mkv` files).

## Dependencies

- Runtime: `ffprobe` and `ffmpeg` (from FFmpeg) must be on `PATH`.
- Go modules: `github.com/olekukonenko/tablewriter` (track status tables), `gopkg.in/alessio/shellescape.v1` (shell-escaping file paths for the ffmpeg command).
- The project intentionally keeps dependencies minimal — don't add new ones without good reason.

## Conventions

- No Go bindings for ffmpeg/ffprobe — the program always shells out to the CLI tools.
- Errors during processing of an individual file are logged and that file is skipped; the tool continues with remaining files.
- File safety: remuxed output is written to a `.tmp.`-prefixed file, then atomically swapped via `os.Rename`. On failure, the temp file is cleaned up. Preserve this pattern for any new file operations.
- Keep `README.md` in sync when user-visible behavior changes (new flags, changed defaults, new dependencies).
- `AGENTS.md` in the repo root contains equivalent guidance for other coding agents — keep it in sync with this file if either is updated.
