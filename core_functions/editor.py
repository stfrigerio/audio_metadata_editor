"""
Main editor logic with refactored functions.
"""

import os
from typing import Optional, List
from enum import Enum

from .audio import AudioFile, find_audio_files, load_audio_file
from .ui import ui
from .config import TAG_CONFIG
from .colors import print_error, print_success, print_info, print_warning

class EditAction(Enum):
    """Possible actions from editing session."""
    CONTINUE_EDITING = "continue_editing"
    NEXT_FILE_SAVE = "next_file_save"
    NEXT_FILE_SKIP = "next_file_skip"
    QUIT_BATCH = "quit_batch"
    QUIT_SESSION = "quit_session"
    BACK_TO_MAIN = "back_to_main"

class MainMenuAction(Enum):
    """Main menu actions."""
    SINGLE_FILE_MODE = "single_file"
    BATCH_PROCESSING = "batch_processing"
    QUIT = "quit"

class EditorSession:
    """Manages an editing session for a single file."""
    
    def __init__(self, audio_file: AudioFile, is_batch: bool = False, 
                 file_index: int = 0, total_files: int = 0):
        self.audio_file = audio_file
        self.is_batch = is_batch
        self.file_index = file_index
        self.total_files = total_files
        self.has_changes = False
    
    def run(self) -> EditAction:
        """Run the interactive editing session."""
        if not self.audio_file.is_loaded:
            print_error("Cannot edit tags, audio file not loaded properly.")
            return EditAction.QUIT_SESSION
        
        while True:
            self._display_current_state()
            ui.display_main_menu(self.is_batch)
            
            choice = ui.get_menu_choice()
            action = self._handle_menu_choice(choice)
            
            if action != EditAction.CONTINUE_EDITING:
                return action
    
    def _display_current_state(self):
        """Display current file state and metadata."""
        ui.display_welcome(
            self.audio_file.filename, 
            self.is_batch, 
            self.file_index, 
            self.total_files
        )
        ui.display_metadata(self.audio_file)
    
    def _handle_menu_choice(self, choice: str) -> EditAction:
        """Handle user's menu choice."""
        # Handle numeric choices for tag editing
        if choice.isdigit():
            return self._handle_tag_edit(choice)
        
        # Handle letter choices
        if choice == 'A':
            return self._handle_album_art_management()
        elif choice == 'S':
            return self._handle_save()
        elif choice == 'R':
            return self._handle_revert()
        elif choice == 'N' and self.is_batch:
            return self._handle_next_save()
        elif choice == 'K' and self.is_batch:
            return self._handle_next_skip()
        elif choice == 'Q':
            return self._handle_quit()
        else:
            ui.show_invalid_option()
            return EditAction.CONTINUE_EDITING
    
    def _handle_tag_edit(self, choice: str) -> EditAction:
        """Handle editing a specific tag."""
        try:
            choice_idx = int(choice) - 1
            tag_keys = list(TAG_CONFIG.keys())
            
            if 0 <= choice_idx < len(tag_keys):
                tag_key = tag_keys[choice_idx]
                config = TAG_CONFIG[tag_key]
                
                current_value = self.audio_file.get_tag_value(tag_key)
                new_value = ui.get_tag_edit_input(tag_key, current_value)
                
                if new_value.upper() == 'CLEAR':
                    new_value = ""
                    print_info(f"{config['name']} will be cleared.")
                elif not new_value and current_value:
                    ui.show_changes_kept_message(config['name'], current_value)
                    return EditAction.CONTINUE_EDITING
                
                if self.audio_file.set_tag_value(tag_key, new_value):
                    ui.show_field_updated_message(config['name'], new_value)
                    self.has_changes = True
            else:
                ui.show_invalid_option()
        except ValueError:
            ui.show_invalid_option()
        
        return EditAction.CONTINUE_EDITING
    
    def _handle_album_art_management(self) -> EditAction:
        """Handle album art management."""
        if self._manage_album_art():
            self.has_changes = True
        return EditAction.CONTINUE_EDITING
    
    def _manage_album_art(self) -> bool:
        """Manage album art operations. Returns True if changes were made."""
        changes_made = False
        
        while True:
            ui.display_album_art_menu(self.audio_file)
            art_choice = ui.get_user_input("Art option").strip().lower()
            
            if art_choice == 'r':
                if self.audio_file.remove_all_album_art():
                    changes_made = True
            elif art_choice == 's':
                img_path = ui.get_image_path_input()
                if img_path and self.audio_file.set_album_art(img_path):
                    changes_made = True
            elif art_choice == 'b':
                break
            else:
                ui.show_invalid_option()
        
        return changes_made
    
    def _handle_save(self) -> EditAction:
        """Handle save operation."""
        if self.audio_file.save():
            self.has_changes = False
        return EditAction.CONTINUE_EDITING
    
    def _handle_revert(self) -> EditAction:
        """Handle revert operation."""
        ui.show_revert_message()
        if self.audio_file.reload():
            self.has_changes = False
        else:
            ui.show_revert_failed_message()
        return EditAction.CONTINUE_EDITING
    
    def _handle_next_save(self) -> EditAction:
        """Handle next file with save (batch mode)."""
        if self.has_changes:
            if not self.audio_file.save():
                # Save failed, ask user what to do
                if ui.confirm_action("Save failed. Skip to next file?"):
                    return EditAction.NEXT_FILE_SKIP
                else:
                    return EditAction.CONTINUE_EDITING
        return EditAction.NEXT_FILE_SAVE
    
    def _handle_next_skip(self) -> EditAction:
        """Handle next file with skip (batch mode)."""
        if self.has_changes:
            if not ui.confirm_action("You have unsaved changes. Skip anyway?"):
                return EditAction.CONTINUE_EDITING
        return EditAction.NEXT_FILE_SKIP
    
    def _handle_quit(self) -> EditAction:
        """Handle quit operation."""
        if self.has_changes:
            if self.is_batch:
                quit_action = ui.confirm_quit_with_changes()
                if quit_action == 'save':
                    if not self.audio_file.save():
                        print_error("Save failed. Not quitting batch.")
                        return EditAction.CONTINUE_EDITING
                elif quit_action == 'cancel':
                    return EditAction.CONTINUE_EDITING
            else:
                if not ui.confirm_action("You have unsaved changes. Quit anyway?"):
                    return EditAction.CONTINUE_EDITING
        
        return EditAction.QUIT_BATCH if self.is_batch else EditAction.QUIT_SESSION

class AudioMetadataEditor:
    """Main editor class that handles single files and batch processing."""
    
    def __init__(self):
        self.current_directory = None
        self.audio_files = []
    
    def run_folder_navigator(self, root_directory: str) -> bool:
        """Run the folder navigation interface starting from root directory."""
        if not os.path.exists(root_directory):
            print_error(f"Directory not found: {root_directory}")
            return False
        
        if not os.path.isdir(root_directory):
            print_error(f"Path is not a directory: {root_directory}")
            return False
        
        current_path = os.path.abspath(root_directory)
        
        while True:
            # Get subdirectories and audio files in current path
            subdirs = self._get_subdirectories(current_path)
            audio_files = find_audio_files(current_path)
            
            # Display navigation menu
            ui.display_folder_navigation_menu(current_path, subdirs, len(audio_files))
            choice = ui.get_menu_choice()
            
            if choice.upper() == 'Q':
                return True
            elif choice.upper() == 'U' and current_path != root_directory:
                # Go up one directory (but not above root)
                parent = os.path.dirname(current_path)
                if parent != current_path and len(parent) >= len(root_directory):
                    current_path = parent
                else:
                    ui.show_invalid_option()
            elif choice.upper() == 'E' and audio_files:
                # Edit audio files in current directory
                if self.run_main_menu(current_path):
                    continue  # Return to navigation after editing
                else:
                    return False
            elif choice.isdigit():
                # Navigate to subdirectory
                dir_index = int(choice) - 1
                if 0 <= dir_index < len(subdirs):
                    current_path = subdirs[dir_index]
                else:
                    ui.show_invalid_option()
            else:
                ui.show_invalid_option()
    
    def _get_subdirectories(self, directory: str) -> list:
        """Get list of subdirectories in the given directory."""
        subdirs = []
        try:
            for item in sorted(os.listdir(directory)):
                item_path = os.path.join(directory, item)
                if os.path.isdir(item_path) and not item.startswith('.'):
                    subdirs.append(item_path)
        except PermissionError:
            print_error(f"Permission denied accessing directory: {directory}")
        
        return subdirs

    def run_main_menu(self, directory: str) -> bool:
        """Run the main menu interface."""
        self.current_directory = directory
        self.audio_files = find_audio_files(directory)
        
        if not self.audio_files:
            print_error(f"No audio files found in: {directory}")
            return False
        
        print_info(f"Found {len(self.audio_files)} audio files in: {directory}")
        
        while True:
            ui.display_main_mode_menu(len(self.audio_files))
            choice = ui.get_menu_choice()
            
            if choice == '1':
                self._run_single_file_mode()
            elif choice == '2':
                self._run_batch_processing_mode()
            elif choice.upper() == 'Q':
                return True
            else:
                ui.show_invalid_option()
    
    def _run_single_file_mode(self):
        """Run single file selection and editing mode."""
        while True:
            ui.display_file_list(self.audio_files, self.current_directory)
            choice = ui.get_menu_choice()
            
            if choice.upper() == 'B':
                break
            elif choice.upper() == 'Q':
                return
            elif choice.isdigit():
                file_index = int(choice) - 1
                if 0 <= file_index < len(self.audio_files):
                    filepath = self.audio_files[file_index]
                    self._edit_single_file(filepath)
                else:
                    ui.show_invalid_option()
            else:
                ui.show_invalid_option()
    
    def _run_batch_processing_mode(self):
        """Run batch processing mode with bulk operations."""
        while True:
            ui.display_batch_menu(len(self.audio_files))
            choice = ui.get_menu_choice()
            
            if choice == '1':
                self._batch_strip_phrases()
            elif choice == '2':
                self._batch_set_genre()
            elif choice == '3':
                self._batch_set_year()
            elif choice == '4':
                self._batch_set_artist()
            elif choice == '5':
                self._batch_set_album_artist()
            elif choice == '6':
                self._batch_set_album_art()
            elif choice.upper() == 'B':
                break
            elif choice.upper() == 'Q':
                return
            else:
                ui.show_invalid_option()
    
    def _batch_strip_phrases(self):
        """Strip phrases from both filename and song title."""
        phrase = ui.get_user_input("Enter phrase to strip from titles and filenames")
        if not phrase:
            return
        
        if not ui.confirm_action(f"Strip '{phrase}' from titles AND filenames of {len(self.audio_files)} files?"):
            return
        
        success_count = 0
        renamed_files = []
        
        for filepath in self.audio_files:
            audio_file = load_audio_file(filepath)
            if audio_file:
                # Strip from title metadata
                current_title = audio_file.get_tag_value('title') or ""
                new_title = current_title.replace(phrase, "").strip()
                
                # Strip from filename
                directory = os.path.dirname(filepath)
                filename = os.path.basename(filepath)
                name, ext = os.path.splitext(filename)
                new_name = name.replace(phrase, "").strip()
                new_filename = new_name + ext
                new_filepath = os.path.join(directory, new_filename)
                
                changes_made = False
                
                # Update metadata if changed
                if new_title != current_title:
                    audio_file.set_tag_value('title', new_title)
                    changes_made = True
                
                # Save metadata changes first
                if changes_made:
                    if not audio_file.save():
                        print_error(f"Failed to save metadata for: {filename}")
                        continue
                
                # Rename file if filename changed and new name doesn't already exist
                if new_filename != filename:
                    if os.path.exists(new_filepath):
                        print_warning(f"Cannot rename '{filename}' - '{new_filename}' already exists")
                    else:
                        try:
                            os.rename(filepath, new_filepath)
                            renamed_files.append((filename, new_filename))
                            changes_made = True
                        except OSError as e:
                            print_error(f"Failed to rename '{filename}': {e}")
                            continue
                
                if changes_made:
                    success_count += 1
                    if new_filename != filename:
                        print_info(f"Updated: {filename} → {new_filename}")
                    else:
                        print_info(f"Updated metadata: {filename}")
        
        # Update the internal file list with new paths
        if renamed_files:
            new_audio_files = []
            for old_path in self.audio_files:
                old_filename = os.path.basename(old_path)
                directory = os.path.dirname(old_path)
                
                # Check if this file was renamed
                new_filename = old_filename
                for old_name, new_name in renamed_files:
                    if old_name == old_filename:
                        new_filename = new_name
                        break
                
                new_audio_files.append(os.path.join(directory, new_filename))
            
            self.audio_files = new_audio_files
        
        print_success(f"Successfully processed {success_count}/{len(self.audio_files)} files")
        if renamed_files:
            print_info(f"Renamed {len(renamed_files)} files")
    
    def _batch_set_genre(self):
        """Set genre for all files."""
        genre = ui.get_user_input("Enter genre to set for all files")
        if not genre:
            return
        
        if not ui.confirm_action(f"Set genre '{genre}' for {len(self.audio_files)} files?"):
            return
        
        self._apply_batch_operation('genre', genre)
    
    def _batch_set_year(self):
        """Set year for all files."""
        year = ui.get_user_input("Enter year to set for all files")
        if not year:
            return
        
        if not ui.confirm_action(f"Set year '{year}' for {len(self.audio_files)} files?"):
            return
        
        self._apply_batch_operation('year', year)
    
    def _batch_set_artist(self):
        """Set artist for all files."""
        artist = ui.get_user_input("Enter artist to set for all files")
        if not artist:
            return
        
        if not ui.confirm_action(f"Set artist '{artist}' for {len(self.audio_files)} files?"):
            return
        
        self._apply_batch_operation('artist', artist)
    
    def _batch_set_album_artist(self):
        """Set album artist for all files."""
        album_artist = ui.get_user_input("Enter album artist to set for all files")
        if not album_artist:
            return
        
        if not ui.confirm_action(f"Set album artist '{album_artist}' for {len(self.audio_files)} files?"):
            return
        
        self._apply_batch_operation('album_artist', album_artist)
    
    def _batch_set_album_art(self):
        """Set the same album art image for all files."""
        img_path = ui.get_image_path_input()
        if not img_path:
            return
        
        # Validate image file exists
        if not os.path.exists(img_path):
            print_error(f"Image file not found: {img_path}")
            return
        
        if not ui.confirm_action(f"Set album art '{os.path.basename(img_path)}' for {len(self.audio_files)} files?"):
            return
        
        success_count = 0
        for filepath in self.audio_files:
            audio_file = load_audio_file(filepath)
            if audio_file:
                if audio_file.set_album_art(img_path):
                    if audio_file.save():
                        success_count += 1
                        print_info(f"Updated: {os.path.basename(filepath)}")
        
        print_success(f"Successfully updated album art for {success_count}/{len(self.audio_files)} files")
    
    def _apply_batch_operation(self, tag_key: str, value: str):
        """Apply a tag operation to all files."""
        success_count = 0
        for filepath in self.audio_files:
            audio_file = load_audio_file(filepath)
            if audio_file:
                if audio_file.set_tag_value(tag_key, value):
                    if audio_file.save():
                        success_count += 1
                        print_info(f"Updated: {os.path.basename(filepath)}")
        
        print_success(f"Successfully updated {success_count}/{len(self.audio_files)} files")
    
    def _edit_single_file(self, filepath: str):
        """Edit a single audio file."""
        audio_file = load_audio_file(filepath)
        if not audio_file:
            return
        
        session = EditorSession(audio_file, is_batch=False)
        session.run()
    
    # Legacy methods for backward compatibility
    def edit_single_file(self, filepath: str) -> bool:
        """Edit a single audio file (legacy method)."""
        if os.path.isfile(filepath):
            self._edit_single_file(filepath)
            return True
        else:
            return self.run_main_menu(filepath)
    
    def edit_batch_files(self, directory: str) -> bool:
        """Edit multiple audio files in a directory (legacy method)."""
        return self.run_main_menu(directory)

def create_editor() -> AudioMetadataEditor:
    """Factory function to create an editor instance."""
    return AudioMetadataEditor() 