#!/usr/bin/env python3
"""
Main entry point for running the audio metadata editor as a module.
Usage: python -m audio_metadata_editor [args]
"""

import sys
import os

# Add the current directory to the path to allow imports
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

def main():
    """Main entry point for module execution."""
    try:
        from core_functions.cli import main as cli_main
        return cli_main()
    except ImportError as e:
        print(f"Error: Could not import CLI modules: {e}")
        print("Please ensure all dependencies are available.")
        print("You may need to install: eyed3, colorama")
        return 1

if __name__ == "__main__":
    sys.exit(main()) 