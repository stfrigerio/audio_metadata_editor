# Audio Metadata Editor

An interactive CLI tool for editing audio file metadata with a colorful interface.

## Features

✨ **Two Operation Modes**: Choose between individual file editing or batch processing  
🎵 **Multiple Formats**: Supports MP3, M4A, FLAC, and OGG files  
🏷️ **Complete Metadata**: Edit title, artist, album, genre, year, track numbers, and more  
🖼️ **Album Art Management**: Add, remove, or replace album artwork (JPEG/PNG)  
📁 **Advanced Batch Processing**: Apply changes to ALL files at once  
🧹 **Smart Cleanup**: Strip unwanted phrases from both filenames and metadata  
⚡ **Fast & Reliable**: Built on the robust eyed3 library  
🎨 **Colored Output**: Beautiful CLI with syntax highlighting and colors  

## Installation Options

### Option 1: Standalone Executable (No Dependencies)
Download the pre-built executable for your platform:
- **Linux**: Download `audio-meta-edit` from [Releases](https://github.com/stfrigerio/audio-metadata-editor/dist)
- **Windows**: Download `audio-meta-edit.exe` from [Releases](https://github.com/stfrigerio/audio-metadata-editor/dist)

```bash
# Make executable (Linux/Mac)
chmod +x audio-meta-edit

# Run directly
./audio-meta-edit /path/to/music/folder
```

### Option 2: Install from PyPI
```bash
pip install audio-metadata-editor
```

### Option 3: Install from Source
```bash
git clone https://github.com/stfrigerio/audio-metadata-editor.git
cd audio-metadata-editor
pip install -e .
```

### Option 4: Build Your Own Executable
```bash
git clone https://github.com/stfrigerio/audio-metadata-editor.git
cd audio-metadata-editor
make build-exe
# Executable will be in dist/audio-meta-edit
```

## Usage

### Command Line
```bash
# Using standalone executable
./audio-meta-edit /path/to/music/folder

# Using installed package
audio-meta-edit /path/to/music/folder
ame /path/to/music/folder  # Short alias

# Options
audio-meta-edit --no-color /path/to/folder  # Disable colors
audio-meta-edit --help                      # Show help
```

## Interface Overview

### Main Menu
When you run the tool on a directory, you'll see:

```
Audio Metadata Editor
Directory contains 31 audio files

Choose Mode:
  1) Edit individual files
  2) Apply changes to all
  Q) Quit
```

### Mode 1: Individual File Editing
Select and edit files one by one with full control:

```
Select File to Edit - /path/to/music
Audio Files (31 found)
  1) Track 01.mp3
  2) Track 02.mp3
  ...
  B) Back to main menu
  Q) Quit
```

Then edit individual metadata fields:
```
──────────────────────────────
Editing: Track 01.mp3
──────────────────────────────

● Current Metadata
  1. Title: Song Title [Official Video]
  2. Artist: Artist Name
  3. Album: Album Name
  4. Album Artist: (empty)
  5. Track: 1/12
  6. Genre: Rock
  7. Year: 2023
  A. Album Art: Present (1 image)

● Options
  1) Edit Title
  2) Edit Artist
  3) Edit Album
  4) Edit Album Artist
  5) Edit Track
  6) Edit Genre
  7) Edit Year
  A) Manage Album Art
  S) Save changes
  Q) Quit
```

### Mode 2: Batch Processing
Apply operations to ALL files at once:

```
Batch Processing Mode
Operations will be applied to ALL 31 files

Batch Operations:
  1) Strip phrases from titles - Remove unwanted text from filenames AND metadata
  2) Set genre for all - Apply same genre to all tracks
  3) Set year for all - Apply same year to all tracks
  4) Set artist for all - Apply same artist to all tracks
  5) Set album artist for all - Apply same album artist to all tracks
  6) Set album art for all - Apply same cover image to all tracks
  B) Back to main menu
  Q) Quit
```

## Key Features

### Smart Phrase Stripping
Remove unwanted text from both filenames and metadata simultaneously:
- Input: `Song Name [Official Video].mp3` → Output: `Song Name.mp3`
- Cleans up downloaded music with suffixes like `[Official Video]`, `- YouTube`, `(Official Audio)`

### Batch Album Art
Set the same cover image for an entire album:
```bash
Enter path to image file: cover.jpg
Set album art 'cover.jpg' for 31 files? (y/N): y
```

### Bulk Metadata Operations
Apply consistent metadata across all files:
- Set the same genre for all tracks in an album
- Update year for entire discography
- Set album artist for compilation albums

## Supported Audio Formats

- **MP3**: ID3v1, ID3v2.3, ID3v2.4 tags
- **M4A**: iTunes-style metadata
- **FLAC**: Vorbis comments
- **OGG**: Vorbis comments

## Dependencies (for non-executable versions)

- **eyed3**: Audio metadata manipulation
- **colorama**: Cross-platform colored terminal output

## Building

```bash
# Development setup
make dev-setup
make install-dev

# Build standalone executable
make build-exe

# Run tests
make test

# Clean build artifacts
make clean
```

