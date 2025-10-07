package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	taglib "go.senan.xyz/taglib"
)

// StripTextFromFiles removes specified text from both metadata title and filename
func StripTextFromFiles(files []string, textToRemove string) []error {
	var errors []error

	for _, filePath := range files {
		if err := stripTextFromFile(filePath, textToRemove); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", filepath.Base(filePath), err))
		}
	}

	return errors
}

func stripTextFromFile(filePath string, textToRemove string) error {
	// Read current metadata
	tags, err := taglib.ReadTags(filePath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	// Get current title and strip text
	currentTitles := tags[taglib.Title]
	if len(currentTitles) == 0 {
		currentTitles = []string{""}
	}
	currentTitle := currentTitles[0]
	newTitle := strings.ReplaceAll(currentTitle, textToRemove, "")
	newTitle = strings.TrimSpace(newTitle)

	// Update title in metadata
	tags[taglib.Title] = []string{newTitle}
	if err := taglib.WriteTags(filePath, tags, 0); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	// Rename file if needed
	dir := filepath.Dir(filePath)
	oldFilename := filepath.Base(filePath)
	ext := filepath.Ext(oldFilename)
	nameWithoutExt := strings.TrimSuffix(oldFilename, ext)

	newFilename := strings.ReplaceAll(nameWithoutExt, textToRemove, "")
	newFilename = strings.TrimSpace(newFilename) + ext

	if oldFilename != newFilename {
		newPath := filepath.Join(dir, newFilename)
		if err := os.Rename(filePath, newPath); err != nil {
			return fmt.Errorf("failed to rename file: %w", err)
		}
	}

	return nil
}

// EditTitleForFiles updates the title metadata for selected files
func EditTitleForFiles(files []string, newTitle string) []error {
	var errors []error

	for _, filePath := range files {
		if err := editTitleForFile(filePath, newTitle); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", filepath.Base(filePath), err))
		}
	}

	return errors
}

func editTitleForFile(filePath string, newTitle string) error {
	tags, err := taglib.ReadTags(filePath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	tags[taglib.Title] = []string{newTitle}
	if err := taglib.WriteTags(filePath, tags, 0); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// EditArtistForFiles updates the artist metadata for selected files
func EditArtistForFiles(files []string, newArtist string) []error {
	var errors []error

	for _, filePath := range files {
		if err := editArtistForFile(filePath, newArtist); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", filepath.Base(filePath), err))
		}
	}

	return errors
}

func editArtistForFile(filePath string, newArtist string) error {
	tags, err := taglib.ReadTags(filePath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	tags[taglib.Artist] = []string{newArtist}
	if err := taglib.WriteTags(filePath, tags, 0); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// EditAlbumForFiles updates the album metadata for selected files
func EditAlbumForFiles(files []string, newAlbum string) []error {
	var errors []error

	for _, filePath := range files {
		if err := editAlbumForFile(filePath, newAlbum); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", filepath.Base(filePath), err))
		}
	}

	return errors
}

func editAlbumForFile(filePath string, newAlbum string) error {
	tags, err := taglib.ReadTags(filePath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	tags[taglib.Album] = []string{newAlbum}
	if err := taglib.WriteTags(filePath, tags, 0); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// EditYearForFiles updates the year metadata for selected files
func EditYearForFiles(files []string, newYear string) []error {
	var errors []error

	for _, filePath := range files {
		if err := editYearForFile(filePath, newYear); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", filepath.Base(filePath), err))
		}
	}

	return errors
}

func editYearForFile(filePath string, newYear string) error {
	tags, err := taglib.ReadTags(filePath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	tags[taglib.Date] = []string{newYear}
	if err := taglib.WriteTags(filePath, tags, 0); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// AddCoverImageToFiles adds a cover image to selected files
func AddCoverImageToFiles(files []string, imagePath string) []error {
	var errors []error

	// Validate image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return []error{fmt.Errorf("image file does not exist: %s", imagePath)}
	}

	// Read image data
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return []error{fmt.Errorf("failed to read image file: %w", err)}
	}

	// Determine MIME type based on extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	default:
		return []error{fmt.Errorf("unsupported image format: %s (use .jpg or .png)", ext)}
	}

	for _, filePath := range files {
		if err := addCoverImageToFile(filePath, imageData, mimeType); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", filepath.Base(filePath), err))
		}
	}

	return errors
}

func addCoverImageToFile(filePath string, imageData []byte, mimeType string) error {
	// Use WriteImage to add cover art
	if err := taglib.WriteImage(filePath, imageData); err != nil {
		return fmt.Errorf("failed to save cover image: %w", err)
	}

	return nil
}
