"""
User interface components for the CLI editor.
"""

from typing import Dict, List, Optional, Tuple
from .config import TAG_CONFIG
from .colors import (
    print_header, print_section, print_success, print_error, 
    print_warning, print_info, highlight, dim, colorize, Colors
)
from .audio import AudioFile
import os

class MenuOption:
    """Represents a menu option."""
    
    def __init__(self, key: str, description: str, action=None):
        self.key = key
        self.description = description
        self.action = action

class EditorUI:
    """Handles user interface for the editor."""
    
    def __init__(self):
        self.tag_keys = list(TAG_CONFIG.keys())
    
    def display_welcome(self, filename: str, is_batch: bool = False, 
                       file_index: int = 0, total_files: int = 0):
        """Display welcome message and file info."""
        if is_batch:
            batch_info = f"[File {file_index + 1}/{total_files}]"
            print_header(f"{batch_info} Editing: {filename}")
        else:
            print_header(f"Editing: {filename}")
    
    def display_metadata(self, audio_file: AudioFile):
        """Display current metadata in a formatted way."""
        print_section("Current Metadata")
        
        # Display tag information
        for i, tag_key in enumerate(self.tag_keys):
            config = TAG_CONFIG[tag_key]
            value = audio_file.get_tag_value(tag_key)
            
            # Format the display with colors
            index_str = colorize(f"{i+1}.", Colors.CYAN, Colors.BRIGHT)
            name_str = colorize(config['name'], Colors.WHITE, Colors.BRIGHT)
            value_str = highlight(value) if value else dim("(empty)")
            
            print(f"  {index_str} {name_str}: {value_str}")
        
        # Display album art information
        self._display_album_art_status(audio_file)
    
    def _display_album_art_status(self, audio_file: AudioFile):
        """Display album art status."""
        art_index = colorize("A.", Colors.CYAN, Colors.BRIGHT)
        art_label = colorize("Album Art", Colors.WHITE, Colors.BRIGHT)
        
        if audio_file.has_album_art():
            art_info = audio_file.get_album_art_info()
            count = len(art_info)
            status = highlight(f"Present ({count} image{'s' if count != 1 else ''})")
            print(f"  {art_index} {art_label}: {status}")
            
            # Show details for each image
            for info in art_info:
                details = dim(f"    - Image {info['index']}: Type {info['type']}, "
                            f"Desc: {info['description']}, MIME: {info['mime_type']}")
                print(details)
        else:
            status = dim("Not Present")
            print(f"  {art_index} {art_label}: {status}")
    
    def display_main_menu(self, is_batch: bool = False):
        """Display the main editing menu."""
        print_section("Options")
        
        # Tag editing options
        for i, tag_key in enumerate(self.tag_keys):
            config = TAG_CONFIG[tag_key]
            option_num = colorize(f"{i+1})", Colors.YELLOW)
            description = f"Edit {config['name']}"
            if 'description' in config:
                description += dim(f" - {config['description']}")
            print(f"  {option_num} {description}")
        
        # Other options
        options = [
            ("A)", "Manage Album Art"),
            ("S)", "Save changes for this file"),
            ("R)", "Revert changes for this file (reload from disk)")
        ]
        
        if is_batch:
            options.extend([
                ("N)", "Done with this file, Save & Next"),
                ("K)", "Done with this file, Skip (don't save) & Next"),
                ("Q)", "Quit ALL batch processing")
            ])
        else:
            options.append(("Q)", "Quit editor"))
        
        for key, desc in options:
            key_colored = colorize(key, Colors.YELLOW)
            print(f"  {key_colored} {desc}")
    
    def display_album_art_menu(self, audio_file: AudioFile):
        """Display album art management menu."""
        print_section("Album Art Management")
        
        if audio_file.has_album_art():
            art_info = audio_file.get_album_art_info()
            count = len(art_info)
            print(f"Current: {highlight(f'{count} image(s) present')}")
            
            for info in art_info:
                details = f"  Image {info['index']}: Type {info['type']}, " \
                         f"Description '{info['description']}', MIME {info['mime_type']}"
                print(dim(details))
            
            print(f"  {colorize('r)', Colors.YELLOW)} Remove ALL existing album art")
        else:
            print(f"Current: {dim('No album art present')}")
        
        print(f"  {colorize('s)', Colors.YELLOW)} Set/Replace album art (from image file)")
        print(f"  {colorize('b)', Colors.YELLOW)} Back to main edit menu")
    
    def get_user_input(self, prompt: str, allow_empty: bool = True) -> str:
        """Get user input with colored prompt."""
        colored_prompt = colorize(f"{prompt}: ", Colors.CYAN, Colors.BRIGHT)
        try:
            response = input(colored_prompt).strip()
            if not allow_empty and not response:
                print_warning("Input cannot be empty.")
                return self.get_user_input(prompt, allow_empty)
            return response
        except (KeyboardInterrupt, EOFError):
            print()  # New line after Ctrl+C
            return ""
    
    def get_menu_choice(self) -> str:
        """Get menu choice from user."""
        return self.get_user_input("Enter option").upper()
    
    def get_tag_edit_input(self, tag_key: str, current_value: str) -> str:
        """Get input for editing a specific tag."""
        config = TAG_CONFIG[tag_key]
        prompt = f"Enter new {config['name']}"
        
        # Show current value and options
        current_display = highlight(current_value) if current_value else dim("(empty)")
        print(f"Current value: {current_display}")
        
        if 'description' in config:
            print(dim(f"Help: {config['description']}"))
        
        print(dim("Enter new value, press Enter to keep current, or type 'CLEAR' to empty"))
        
        return self.get_user_input(prompt)
    
    def get_image_path_input(self) -> str:
        """Get image file path from user."""
        return self.get_user_input("Enter path to image file (e.g., cover.jpg)")
    
    def confirm_action(self, message: str, default: bool = False) -> bool:
        """Get yes/no confirmation from user."""
        suffix = "(y/N)" if not default else "(Y/n)"
        response = self.get_user_input(f"{message} {suffix}").lower()
        
        if not response:
            return default
        
        return response in ('y', 'yes')
    
    def confirm_quit_with_changes(self) -> str:
        """Get confirmation for quitting with unsaved changes."""
        options = "(y/N/c=cancel quit)" 
        response = self.get_user_input(f"Save changes to current file before quitting batch? {options}").lower()
        
        if response in ('y', 'yes'):
            return 'save'
        elif response in ('c', 'cancel'):
            return 'cancel'
        else:
            return 'discard'
    
    def show_processing_info(self, directory: str, file_count: int):
        """Show batch processing information."""
        print_info(f"Processing folder: {directory}")
        if file_count == 0:
            print_warning("No compatible audio files found in the folder.")
        else:
            print_success(f"Found {file_count} audio file(s) to process.")
    
    def show_batch_complete(self):
        """Show batch processing completion message."""
        print_success("Batch processing finished.")
    
    def show_batch_quit(self):
        """Show batch processing quit message."""
        print_info("Batch processing quit by user.")
    
    def show_skip_file_error(self, filename: str):
        """Show message when skipping a file due to error."""
        print_warning(f"Skipping {filename} due to load error.")
    
    def show_invalid_option(self):
        """Show invalid option message."""
        print_warning("Invalid option. Please try again.")
    
    def show_no_changes_message(self):
        """Show message when no changes were made."""
        print_info("No changes made.")
    
    def show_changes_kept_message(self, field_name: str, value: str):
        """Show message when field value is kept unchanged."""
        print_info(f"{field_name} kept as '{value}'.")
    
    def show_field_updated_message(self, field_name: str, value: str):
        """Show message when field is updated."""
        display_value = value if value else "(cleared)"
        print_success(f"{field_name} set to '{display_value}'.")
    
    def show_revert_message(self):
        """Show message when reverting changes."""
        print_info("Reverting changes for this file...")
        print_success("Tags reverted to original state from file.")
    
    def show_revert_failed_message(self):
        """Show revert failed message."""
        print_error("Failed to revert changes. File may be corrupted or inaccessible.")

    def display_main_mode_menu(self, file_count: int):
        """Display the main mode selection menu."""
        print_header("Audio Metadata Editor")
        print_info(f"Directory contains {file_count} audio files")
        print_section("Choose Mode")
        
        options = [
            ("1)", "Edit individual files"),
            ("2)", "Apply changes to all"),
            ("Q)", "Quit")
        ]
        
        for key, desc in options:
            key_colored = colorize(key, Colors.YELLOW)
            print(f"  {key_colored} {desc}")
    
    def display_file_list(self, audio_files: List[str], directory: str):
        """Display list of audio files for selection."""
        print_header(f"Select File to Edit - {directory}")
        print_section(f"Audio Files ({len(audio_files)} found)")
        
        for i, filepath in enumerate(audio_files):
            filename = os.path.basename(filepath)
            index_str = colorize(f"{i+1})", Colors.CYAN, Colors.BRIGHT)
            print(f"  {index_str} {filename}")
        
        print()
        options = [
            ("B)", "Back to main menu"),
            ("Q)", "Quit")
        ]
        
        for key, desc in options:
            key_colored = colorize(key, Colors.YELLOW)
            print(f"  {key_colored} {desc}")
    
    def display_batch_menu(self, file_count: int):
        """Display batch processing menu."""
        print_header("Batch Processing Mode")
        print_info(f"Operations will be applied to ALL {file_count} files")
        print_section("Batch Operations")
        
        options = [
            ("1)", "Strip phrases from titles"),
            ("2)", "Set genre for all"),
            ("3)", "Set year for all"),
            ("4)", "Set artist for all"),
            ("5)", "Set album artist for all"),
            ("6)", "Set album art for all"),
            ("B)", "Back to main menu"),
            ("Q)", "Quit")
        ]
        
        for key, desc in options:
            key_colored = colorize(key, Colors.YELLOW)
            print(f"  {key_colored} {desc}")

    def display_folder_navigation_menu(self, current_path: str, subdirs: List[str], audio_file_count: int):
        """Display folder navigation menu."""
        print_header("Folder Navigation")
        
        # Show current path
        print_info(f"Current: {current_path}")
        
        # Show audio files status
        if audio_file_count > 0:
            audio_status = highlight(f"{audio_file_count} audio files found")
            print(f"Audio files: {audio_status}")
        else:
            audio_status = dim("No audio files in this directory")
            print(f"Audio files: {audio_status}")
        
        print_section("Navigation Options")
        
        # Show subdirectories
        if subdirs:
            print(dim("Subdirectories:"))
            for i, subdir in enumerate(subdirs):
                dir_name = os.path.basename(subdir)
                index_str = colorize(f"{i+1})", Colors.CYAN, Colors.BRIGHT)
                print(f"  {index_str} {dir_name}")
            print()
        else:
            print(dim("No subdirectories found"))
            print()
        
        # Show available actions
        options = []
        
        # Add edit option if audio files are present
        if audio_file_count > 0:
            options.append(("E)", f"Edit audio files in this directory ({audio_file_count} files)"))
        
        # Add navigation options
        options.extend([
            ("U)", "Go up one directory"),
            ("Q)", "Quit")
        ])
        
        for key, desc in options:
            key_colored = colorize(key, Colors.YELLOW)
            print(f"  {key_colored} {desc}")

# Global UI instance
ui = EditorUI() 