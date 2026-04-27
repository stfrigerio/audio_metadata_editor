# audio-meta-edit 🎵

TUI audio metadata editor. Browse a directory tree of music files,
preview tags, and edit them in place — title, artist, album, album
artist, track, year, genre, and embedded cover art. Built on Bubble
Tea.

## Build

Needs Go ≥ 1.22.

```bash
cd go
make build
./audio-meta-edit /path/to/music
```

Or drop a symlink on your `$PATH` — see `INSTALL.md`.

## Usage

```bash
audio-meta-edit /path/to/music
```

If invoked with no argument the wrapper script defaults to the current
working directory.

## Requires

- Linux (tested on Arch, should work anywhere Go runs)
- `eyeD3`-tagged MP3s read fine; other formats pass through `dhowden/tag`
