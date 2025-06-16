"""
Configuration for the audio metadata editor.
"""

try:
    import eyed3
    import eyed3.id3
    import eyed3.core
except ImportError as e:
    print(f"Error: eyed3 library not found: {e}")
    print("Please install eyed3: pip install eyed3")
    exit(1)

from typing import Dict, Any, Callable, Optional

def set_track_num(tag, value_str: str) -> None:
    """Set track number from string format 'num' or 'num/total'."""
    if not value_str:
        tag.track_num = (None, None)
        return
    
    parts = value_str.split('/')
    track_val, total_val = None, None
    
    try:
        if len(parts) >= 1 and parts[0].strip():
            track_val = int(parts[0].strip())
        if len(parts) == 2 and parts[1].strip():
            total_val = int(parts[1].strip())
        tag.track_num = (track_val, total_val)
    except ValueError:
        raise ValueError("Invalid track format. Use 'num' or 'num/total'. Example: 5 or 5/12")

def set_year(tag, value_str: str) -> None:
    """Set year from string format 'YYYY'."""
    if not value_str:
        tag.recording_date = None
        return
    
    try:
        year_val = int(value_str)
        tag.recording_date = eyed3.core.Date(year_val)
    except ValueError:
        raise ValueError("Invalid year format. Use YYYY. Example: 2023")

def set_genre(tag, value_str: str) -> None:
    """Set genre from string."""
    if value_str:
        tag.genre = eyed3.id3.Genre(name=value_str)
    else:
        tag.genre = None

# Tag configuration mapping
TAG_CONFIG: Dict[str, Dict[str, Any]] = {
    "title": {
        "attr": "title", 
        "name": "Title",
        "description": "Song title"
    },
    "artist": {
        "attr": "artist", 
        "name": "Artist",
        "description": "Artist name"
    },
    "album": {
        "attr": "album", 
        "name": "Album",
        "description": "Album name"
    },
    "album_artist": {
        "attr": "album_artist", 
        "name": "Album Artist",
        "description": "Album artist (if different from track artist)"
    },
    "track": {
        "attr": "track_num",
        "name": "Track",
        "description": "Track number (format: num or num/total)",
        "getter": lambda tag: f"{tag.track_num[0] or ''}{f'/{tag.track_num[1]}' if tag.track_num[1] else ''}" if tag.track_num else "",
        "setter": set_track_num
    },
    "genre": {
        "attr": "genre",
        "name": "Genre",
        "description": "Music genre",
        "getter": lambda tag: tag.genre.name if tag.genre else "",
        "setter": set_genre
    },
    "year": {
        "attr": "recording_date",
        "name": "Year",
        "description": "Recording year",
        "getter": lambda tag: str(tag.recording_date.year) if tag.recording_date else "",
        "setter": set_year
    }
}

# Supported audio file extensions
SUPPORTED_EXTENSIONS = (".mp3", ".m4a", ".flac", ".ogg")

# Default ID3 version for saving
try:
    DEFAULT_ID3_VERSION = eyed3.id3.ID3_V2_3
except AttributeError:
    # Fallback if eyed3.id3 is not available
    DEFAULT_ID3_VERSION = (2, 3, 0) 