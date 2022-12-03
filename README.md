# mkvtrackcleaner
Tool to bulk-remux mkv-files from tracks of unwanted languages.

It will go through all files in a given directory or work on a single
file. Then it will identify all tracks in that file and remove all
tracks marked with a language that isn't in the list of wanted
languages.

## Notes about languages
Always keep `und` as a language. Lots of files out there with only one
audio or subtitle track got it's only track marked as `undefined`
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
Either edit the `WANTED_LANGS` array at the top of the script or create
a config-file to alter the same variable to define which languages you
want to keep, then run the script like this:
```bash
./mkvtrackcleaner /path/to/directory
```

or like this:
```bash
./mkvtrackcleaner /path/to/file.mkv
```

# Configfile
The config is a simple bash-script that the script sources, it allows
you to override some variables.

Put the configfile in `$XDG_CONFIG_DIR/mkvtrackcleaner.conf` which by
default is `$HOME/.config/mkvtrackcleaner.conf`.

So, for example if you want to extend the `WANTED_LANGS` list you could
have a file that looks like this (to add norwegian):
```bash
WANTED_LANGS+=(
    "nor"
)
```

Or you can just override the entire list and roll your own:
```bash
WANTED_LANGS=(
    "und"
    "eng"
)
```

Variables that's possible to alter includes:
- `$WANTED_LANGS` -- Array with wanted languages in three letter codes.
- `$GREP_MATCH_AUDIO` -- Regex to match audio tracks.
- `$GREP_MATCH_SUBS` -- Regex to match subtitle tracks.
- `$SED_EXTRACT_INFO` -- Regex to extract info about a track.

# Tested on
- Archlinux x86_64 2016-04-03
- bash 4.3.42
- mkvmerge v9.0.0
