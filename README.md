# mkvtrackcleaner
Tool to bulk-remux mkv-files from tracks of unwanted languages.

It will go through all files in a given directory or work on a single
file. Then it will identify all tracks in that file and remove all
tracks marked with a language that isn't in the list of wanted
languages.

## Notes about languages
Always keep the `und` as a language. Lots of files out there with only
one audio or subtitle track got it's only track marked as `undefined`
language, so you probably always want to have `und` in your list of
wanted languages.

## Notes about audio tracks
If the script filter away all audio tracks, it will choose to not touch
them at all. Instead it will keep all the audio tracks.

The script will also disable the default-track flag from all
audio-tracks to let the player decide which track to play. This may be
based on order of the tracks or the language of the track. It differs
from player to player.

## Notes about subtitle tracks
This script may remove all subtitle tracks if there's no tracks matching
the wanted languages list. So you may end up without subtitles.

This script will also disable the default-track flag from all subtitles
to let it be up to the player to decide which one to use, based on order
or language tag of the track. It differs from player to player.

# Dependencies
- bash (at least version 4)
- mkvmerge (part of mkvtoolnix cli package)
- coreutils
- find
- grep
- sed

# Usage
Edit the `WANTED_LANGS` array at the top of the file to specify
languages that you want to keep, then run the script like this:
```bash
./mkvtrackcleaner /path/to/directory
```

or like this:
```bash
./mkvtrackcleaner /path/to/file.mkv
```

# TODO
Make some kind of config-file to avoid having to edit the file to set
languages.
