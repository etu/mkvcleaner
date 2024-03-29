[![Check](https://github.com/etu/mkvcleaner/actions/workflows/check.yml/badge.svg)](https://github.com/etu/mkvcleaner/actions/workflows/check.yml)
[![Update](https://github.com/etu/mkvcleaner/actions/workflows/update.yml/badge.svg)](https://github.com/etu/mkvcleaner/actions/workflows/update.yml)

# mkvcleaner
Tool to bulk-remux mkv-files from tracks of unwanted languages.

It will go through all files in a given directory or work on a single
file. Then it will identify all tracks in that file and remove all tracks
marked with a language that isn't in the list of wanted languages.

## Usage
You can specify either files or directories as arguments. The default set of
languages specified is `und,eng,swe,jap,jpn`. You can override this list
using the `--langs` flag.

```bash
./mkvcleaner [--langs=und,eng] [--automatic] path/to/directory [path/to/file.mkv] [path/to/other/directory] […]
```

By default it will prompt the user about the changes to a file to approve the
changes before it's executed. However, if the `--automatic` flag is provided
it will skip the confirmation.

## Notes about languages
Always keep `und` as a language. Lots of files out there with only one audio
or subtitle track got it's only track marked as `undefined` language, so you
probably always want to have `und` in your list of wanted languages.

## Notes about audio tracks
If the script filter away all audio tracks, it will choose to not touch them
at all. Instead it will keep all the audio tracks.

## Notes about subtitle tracks
This script may remove all subtitle tracks if there's no tracks matching the
wanted languages list. So you may end up without subtitles.

## Dependencies
It's a go program and it depends on `ffprobe` and `ffmpeg` from the `ffmpeg`
project. `ffprobe` is used to detect changes.
