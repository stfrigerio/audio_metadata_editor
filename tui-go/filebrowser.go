package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// openFileBrowser opens the file browser starting from home directory
func (m *model) openFileBrowser() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/"
	}

	m.showFileBrowser = true
	m.browserDir = homeDir
	m.browserCursor = 0
	m.pendingAction = "add-cover"
	m.refreshBrowserContents()
}

// refreshBrowserContents reads the current directory and populates browser lists
func (m *model) refreshBrowserContents() {
	entries, err := os.ReadDir(m.browserDir)
	if err != nil {
		m.errorMessage = "Cannot read directory: " + err.Error()
		return
	}

	m.browserDirs = make([]string, 0)
	m.browserFiles = make([]string, 0)

	// Add parent directory option
	if m.browserDir != "/" {
		m.browserDirs = append(m.browserDirs, "..")
	}

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			m.browserDirs = append(m.browserDirs, entry.Name())
		} else {
			// Only show image files
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				m.browserFiles = append(m.browserFiles, entry.Name())
			}
		}
	}

	sort.Strings(m.browserDirs)
	sort.Strings(m.browserFiles)

	// Reset cursor if out of bounds
	totalItems := len(m.browserDirs) + len(m.browserFiles)
	if m.browserCursor >= totalItems {
		m.browserCursor = totalItems - 1
	}
	if m.browserCursor < 0 {
		m.browserCursor = 0
	}
}

// handleBrowserKey handles keyboard input in file browser mode
func (m *model) handleBrowserKey(key string) {
	switch key {
	case "esc":
		m.showFileBrowser = false
		m.browserDir = ""
		m.browserFiles = nil
		m.browserDirs = nil
		m.browserCursor = 0

	case "up":
		m.browserCursor--
		if m.browserCursor < 0 {
			m.browserCursor = 0
		}

	case "down":
		totalItems := len(m.browserDirs) + len(m.browserFiles)
		m.browserCursor++
		if m.browserCursor >= totalItems {
			m.browserCursor = totalItems - 1
		}

	case "enter":
		m.handleBrowserSelect()
	}
}

// handleBrowserSelect handles selection in file browser
func (m *model) handleBrowserSelect() {
	totalDirs := len(m.browserDirs)

	if m.browserCursor < totalDirs {
		// Selected a directory
		dirName := m.browserDirs[m.browserCursor]

		if dirName == ".." {
			// Go to parent directory
			m.browserDir = filepath.Dir(m.browserDir)
		} else {
			// Enter directory
			m.browserDir = filepath.Join(m.browserDir, dirName)
		}

		m.browserCursor = 0
		m.refreshBrowserContents()
	} else {
		// Selected a file
		fileIndex := m.browserCursor - totalDirs
		if fileIndex < len(m.browserFiles) {
			selectedFile := filepath.Join(m.browserDir, m.browserFiles[fileIndex])

			// Execute the action with the selected file
			selectedFiles := make([]string, 0, len(m.selectedFiles))
			for file := range m.selectedFiles {
				selectedFiles = append(selectedFiles, file)
			}

			errors := AddCoverImageToFiles(selectedFiles, selectedFile)

			// Close browser
			m.showFileBrowser = false
			m.browserDir = ""
			m.browserFiles = nil
			m.browserDirs = nil
			m.browserCursor = 0

			if len(errors) > 0 {
				m.errorMessage = errors[0].Error()
			} else {
				m.errorMessage = ""
			}

			// Clear selection
			m.selectedFiles = make(map[string]bool)
			if len(m.getCurrentNode().Files) == 0 {
				m.inFilePanel = false
			}
		}
	}
}
