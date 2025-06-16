"""
Command-line interface for the audio metadata editor.
"""

import argparse
import os
import sys
from typing import Optional

from . import __version__
from .editor import create_editor
from .colors import print_error, print_info, print_success

def create_parser() -> argparse.ArgumentParser:
    """Create and configure the argument parser."""
    parser = argparse.ArgumentParser(
        description="Interactive CLI tool to view and edit audio file metadata using eyed3.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s song.mp3                    # Edit a single audio file
  %(prog)s /path/to/music/folder       # Edit all audio files in a folder
  %(prog)s /path/to/music/library      # Navigate through music library folders
  %(prog)s --version                   # Show version information

Supported formats: MP3, M4A, FLAC, OGG

The tool provides two main modes:

1. FOLDER NAVIGATION MODE:
   - Navigate through directory structures
   - Only shows editing options when audio files are found
   - Perfect for large music libraries with nested folders

2. DIRECT EDITING MODE:
   - Edit files in a specific directory
   - Individual file editing or batch processing
   - Full metadata editing capabilities

The tool provides an interactive interface for editing:
- Title, Artist, Album, Album Artist
- Track number (format: num or num/total)
- Genre, Year
- Album artwork (JPEG/PNG images)

In batch mode, you can process multiple files with options to save, skip, or quit.
        """
    )
    
    parser.add_argument(
        'path',
        help="Path to the audio file or folder containing audio files"
    )
    
    parser.add_argument(
        '--version',
        action='version',
        version=f'%(prog)s {__version__}'
    )
    
    parser.add_argument(
        '--no-color',
        action='store_true',
        help="Disable colored output"
    )
    
    return parser

def validate_path(path: str) -> bool:
    """Validate the provided path."""
    if not os.path.exists(path):
        print_error(f"Path not found: '{path}'")
        return False
    
    if not (os.path.isfile(path) or os.path.isdir(path)):
        print_error(f"Path is neither a file nor a directory: '{path}'")
        return False
    
    return True

def main() -> int:
    """Main CLI entry point."""
    try:
        parser = create_parser()
        args = parser.parse_args()
        
        # Handle --no-color option
        if args.no_color:
            # Disable colorama colors
            os.environ['NO_COLOR'] = '1'
        
        # Validate path
        if not validate_path(args.path):
            return 1
        
        # Create editor and process files
        editor = create_editor()
        
        if os.path.isfile(args.path):
            print_info(f"Loading audio file: {os.path.basename(args.path)}")
            success = editor.edit_single_file(args.path)
        else:
            # It's a directory - check if it has audio files
            from .audio import find_audio_files
            audio_files = find_audio_files(args.path)
            
            if audio_files:
                # Directory has audio files - ask user which mode to use
                print_info(f"Found {len(audio_files)} audio files in: {args.path}")
                print_info("Choose mode:")
                print("  1) Navigate through subdirectories")
                print("  2) Edit files in this directory directly")
                
                while True:
                    choice = input("Enter choice (1 or 2): ").strip()
                    if choice == '1':
                        success = editor.run_folder_navigator(args.path)
                        break
                    elif choice == '2':
                        success = editor.run_main_menu(args.path)
                        break
                    else:
                        print("Invalid choice. Please enter 1 or 2.")
            else:
                # No audio files in root directory - use folder navigation
                print_info(f"No audio files found in root directory: {args.path}")
                print_info("Starting folder navigation mode...")
                success = editor.run_folder_navigator(args.path)
        
        if success:
            print_info("Editor session completed.")
            return 0
        else:
            return 1
            
    except KeyboardInterrupt:
        print_info("\nEditor interrupted by user.")
        return 130  # Standard exit code for Ctrl+C
    except Exception as e:
        print_error(f"Unexpected error: {e}")
        return 1

def cli_entry_point():
    """Entry point for setuptools console scripts."""
    sys.exit(main())

if __name__ == "__main__":
    cli_entry_point() 