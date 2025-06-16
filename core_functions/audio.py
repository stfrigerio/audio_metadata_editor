"""
Audio file operations using eyed3.
"""

import eyed3
import os
import mimetypes
from typing import Optional, List
from .config import TAG_CONFIG, DEFAULT_ID3_VERSION, SUPPORTED_EXTENSIONS
from .colors import print_error, print_info, print_success, error, success

class AudioFile:
    """Wrapper class for audio file operations."""
    
    def __init__(self, filepath: str):
        self.filepath = filepath
        self.filename = os.path.basename(filepath)
        self._audiofile = None
        self._load()
    
    def _load(self) -> bool:
        """Load the audio file using eyed3."""
        try:
            # Suppress eyed3 internal logging
            eyed3.log.setLevel("ERROR")
            
            self._audiofile = eyed3.load(self.filepath)
            if self._audiofile is None:
                print_error(f"Could not load '{self.filename}'. Is it a valid audio file?")
                return False
            
            if self._audiofile.tag is None:
                print_info("No existing ID3 tag found. Initializing a new one.")
                self._audiofile.initTag(version=DEFAULT_ID3_VERSION)
            
            return True
        except Exception as e:
            print_error(f"Error loading file '{self.filename}': {e}")
            return False
    
    @property
    def is_loaded(self) -> bool:
        """Check if audio file is properly loaded."""
        return self._audiofile is not None and self._audiofile.tag is not None
    
    @property
    def tag(self):
        """Get the tag object."""
        return self._audiofile.tag if self._audiofile else None
    
    def get_tag_value(self, tag_key: str) -> str:
        """Get formatted tag value for display."""
        if not self.is_loaded or tag_key not in TAG_CONFIG:
            return ""
        
        config = TAG_CONFIG[tag_key]
        
        if 'getter' in config:
            return config['getter'](self.tag)
        elif hasattr(self.tag, config['attr']):
            raw_value = getattr(self.tag, config['attr'])
            return str(raw_value) if raw_value is not None else ""
        
        return ""
    
    def set_tag_value(self, tag_key: str, value: str) -> bool:
        """Set tag value with validation."""
        if not self.is_loaded or tag_key not in TAG_CONFIG:
            return False
        
        config = TAG_CONFIG[tag_key]
        
        try:
            if 'setter' in config:
                config['setter'](self.tag, value)
            else:
                setattr(self.tag, config['attr'], value if value else None)
            return True
        except ValueError as e:
            print_error(str(e))
            return False
        except Exception as e:
            print_error(f"Error setting {config['name']}: {e}")
            return False
    
    def save(self) -> bool:
        """Save changes to the audio file."""
        if not self.is_loaded:
            return False
        
        try:
            self.tag.save(version=DEFAULT_ID3_VERSION, encoding='utf-8')
            print_success(f"Changes saved successfully to '{self.filename}'!")
            return True
        except Exception as e:
            print_error(f"Error saving tags to '{self.filename}': {e}")
            return False
    
    def reload(self) -> bool:
        """Reload the audio file from disk."""
        return self._load()
    
    def has_album_art(self) -> bool:
        """Check if file has album art."""
        return self.is_loaded and bool(self.tag.images)
    
    def get_album_art_info(self) -> List[dict]:
        """Get information about album art."""
        if not self.has_album_art():
            return []
        
        art_info = []
        for i, img in enumerate(self.tag.images):
            art_info.append({
                'index': i + 1,
                'type': img.picture_type,
                'description': img.description or "(no description)",
                'mime_type': img.mime_type
            })
        
        return art_info
    
    def remove_all_album_art(self) -> bool:
        """Remove all album art."""
        if not self.is_loaded:
            return False
        
        if self.tag.images:
            # Remove all images by iterating and removing each one
            for img in list(self.tag.images):
                self.tag.images.remove(img.description)
            print_success("All album art removed.")
            return True
        else:
            print_info("No album art to remove.")
            return False
    
    def set_album_art(self, image_path: str) -> bool:
        """Set album art from image file."""
        if not self.is_loaded:
            return False
        
        if not os.path.exists(image_path):
            print_error(f"Image file not found at '{image_path}'")
            return False
        
        if not os.path.isfile(image_path):
            print_error(f"'{image_path}' is not a file.")
            return False
        
        try:
            with open(image_path, "rb") as img_f:
                img_data = img_f.read()
            
            # Determine MIME type
            mime_type = mimetypes.guess_type(image_path)[0]
            if not mime_type:
                # Fallback for common types
                ext = image_path.lower()
                if ext.endswith((".jpg", ".jpeg")):
                    mime_type = "image/jpeg"
                elif ext.endswith(".png"):
                    mime_type = "image/png"
                else:
                    print_error("Could not determine image MIME type. Please ensure it's a JPEG or PNG.")
                    return False
            
            # Remove existing art first
            if self.tag.images:
                for img in list(self.tag.images):
                    self.tag.images.remove(img.description)
            
            # Set new album art
            self.tag.images.set(3, img_data, mime_type, description="Cover")
            print_success(f"Album art set from '{os.path.basename(image_path)}'.")
            return True
            
        except Exception as e:
            print_error(f"Error setting album art: {e}")
            return False

def find_audio_files(directory: str) -> List[str]:
    """Find all supported audio files in a directory."""
    if not os.path.isdir(directory):
        return []
    
    audio_files = []
    try:
        for filename in sorted(os.listdir(directory)):
            if filename.lower().endswith(SUPPORTED_EXTENSIONS):
                audio_files.append(os.path.join(directory, filename))
    except PermissionError:
        print_error(f"Permission denied accessing directory: {directory}")
    
    return audio_files

def load_audio_file(filepath: str) -> Optional[AudioFile]:
    """Load an audio file and return AudioFile wrapper."""
    audio_file = AudioFile(filepath)
    return audio_file if audio_file.is_loaded else None 