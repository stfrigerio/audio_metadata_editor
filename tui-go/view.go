package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	taglib "go.senan.xyz/taglib"
)

// View renders the current state with bordered panels
func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Show preview overlay if active
	if m.showPreview {
		return m.renderPreview()
	}

	// Show file browser overlay if active
	if m.showFileBrowser {
		return m.renderFileBrowser()
	}

	// Show input prompt overlay if active
	if m.showInputPrompt {
		return m.renderInputPrompt()
	}

	var s strings.Builder

	// Calculate panel dimensions
	leftWidth := m.width/2 - 1
	rightWidth := m.width - leftWidth - 2
	panelHeight := m.height - 2 // Leave room for footer

	// Get current node
	currentNode := m.getCurrentNode()

	// Get panel content
	treeContent := m.renderTreeContent()
	filesContent := m.renderFilesContent(currentNode)

	// If edit menu is shown, create smaller panels and separate edit menu
	if m.showEditMenu {
		menuHeight := 8
		treeHeight := panelHeight - menuHeight // No spacing needed, borders touch

		// Determine focus
		treeFocused := false
		filesFocused := false

		// Create smaller tree panel (left top)
		treePanel := renderBorderedPanel("Files", treeContent, leftWidth, treeHeight, treeFocused)

		// Create edit menu panel (left bottom)
		menuContent := m.renderEditMenu()
		menuPanel := renderBorderedPanel("Edit Options", menuContent, leftWidth, menuHeight, true)

		// Create full-height Details panel (right)
		detailsPanel := renderBorderedPanel("Details", filesContent, rightWidth, panelHeight, filesFocused)

		// Combine vertically stacked left panels with right panel
		leftColumn := append(treePanel, menuPanel...)

		// Combine with right panel side by side
		s.WriteString(combinePanelsSideBySide(leftColumn, detailsPanel))

	} else {
		// Determine focus
		treeFocused := !m.inFilePanel
		filesFocused := m.inFilePanel

		// Create bordered panels
		leftPanel := renderBorderedPanel("Files", treeContent, leftWidth, panelHeight, treeFocused)
		rightPanel := renderBorderedPanel("Details", filesContent, rightWidth, panelHeight, filesFocused)

		// Combine panels
		s.WriteString(combinePanelsSideBySide(leftPanel, rightPanel))
	}

	// Footer
	s.WriteString("\033[90m")
	s.WriteString(strings.Repeat("─", m.width))
	s.WriteString("\033[0m\n")

	// Show error message or footer
	if m.errorMessage != "" {
		s.WriteString("\033[31m")
		s.WriteString(m.errorMessage)
		s.WriteString("\033[0m \033[90m[esc to clear]\033[0m")
	} else {
		s.WriteString(m.renderFooter())
	}

	return s.String()
}

// renderTreeContent generates tree content without borders
func (m model) renderTreeContent() []string {
	lines := make([]string, 0)

	for i, node := range m.flattenedTree {
		var line strings.Builder

		// Cursor
		if !m.inFilePanel && i == m.cursor {
			line.WriteString("\033[36m> \033[0m")
		} else {
			line.WriteString("  ")
		}

		// Indentation
		line.WriteString(strings.Repeat("  ", node.Depth))

		// Expand/collapse indicator
		if len(node.Children) > 0 {
			if node.Expanded {
				line.WriteString("\033[33m▼ \033[0m")
			} else {
				line.WriteString("\033[33m▶ \033[0m")
			}
		} else {
			line.WriteString("  ")
		}

		// Folder name
		if !m.inFilePanel && i == m.cursor {
			line.WriteString("\033[1;36m")
		} else {
			line.WriteString("\033[34m")
		}
		line.WriteString(node.Name)
		if node.Depth > 0 {
			line.WriteString("/")
		}
		line.WriteString("\033[0m")

		// File count
		if len(node.Files) > 0 {
			line.WriteString(fmt.Sprintf(" \033[90m(%d)\033[0m", len(node.Files)))
		}

		lines = append(lines, line.String())
	}

	return lines
}

// renderFilesContent generates file list content without borders
func (m model) renderFilesContent(node *TreeNode) []string {
	lines := make([]string, 0)

	if node == nil {
		lines = append(lines, "\033[90mNo folder selected\033[0m")
		return lines
	}

	if len(node.Files) == 0 {
		lines = append(lines, "\033[90mNo audio files\033[0m")
		return lines
	}

	for i, file := range node.Files {
		var line strings.Builder

		// Cursor
		if m.inFilePanel && i == m.fileCursor {
			line.WriteString("\033[36m> \033[0m")
		} else {
			line.WriteString("  ")
		}

		// Checkbox
		if m.selectedFiles[file] {
			line.WriteString("\033[32m[X]\033[0m ")
		} else {
			line.WriteString("\033[90m[ ]\033[0m ")
		}

		// Filename
		if m.inFilePanel && i == m.fileCursor {
			line.WriteString("\033[1m")
		}
		line.WriteString(filepath.Base(file))
		if m.inFilePanel && i == m.fileCursor {
			line.WriteString("\033[0m")
		}

		lines = append(lines, line.String())
	}

	// Selection count
	selectedCount := 0
	for _, file := range node.Files {
		if m.selectedFiles[file] {
			selectedCount++
		}
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("\033[90mSelected: \033[32m%d\033[90m/\033[0m%d", selectedCount, len(node.Files)))

	return lines
}
func (m model) renderEditMenu() []string {
	lines := make([]string, 0)

	options := []string{
		"Strip text from titles",
		"Add cover image",
		"Edit Title",
		"Edit Artist",
		"Edit Album",
		"Edit Year",
	}

	for i, opt := range options {
		var line strings.Builder

		// Cursor
		if i == m.editCursor {
			line.WriteString("\033[36m> \033[0m")
		} else {
			line.WriteString("  ")
		}

		// Option
		if i == m.editCursor {
			line.WriteString("\033[1m")
		}
		line.WriteString(opt)
		if i == m.editCursor {
			line.WriteString("\033[0m")
		}

		lines = append(lines, line.String())
	}

	return lines
}

// renderFooter generates the footer with keyboard shortcuts
func (m model) renderFooter() string {
	if m.showEditMenu {
		return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]select [\033[33mesc\033[90m]close\033[0m"
	} else if m.inFilePanel {
		hasSelected := len(m.selectedFiles) > 0
		if hasSelected {
			return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]toggle [\033[33ma\033[90m]toggle all [\033[33me\033[90m]edit [\033[33mp\033[90m]preview [\033[33mtab\033[90m]tree [\033[33mq\033[90m]quit\033[0m"
		} else {
			return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]toggle [\033[33ma\033[90m]toggle all [\033[33mp\033[90m]preview [\033[33mtab\033[90m]tree [\033[33mq\033[90m]quit\033[0m"
		}
	} else {
		return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]open/toggle [\033[33mp\033[90m]preview [\033[33mtab\033[90m]files [\033[33mq\033[90m]quit\033[0m"
	}
}

// renderInputPrompt renders the input prompt with live preview or current values
func (m model) renderInputPrompt() string {
	// Check if this is a strip action (needs highlighting) or edit action (shows current values)
	isStripAction := m.pendingAction == "strip"

	if isStripAction {
		return m.renderStripPrompt()
	} else {
		return m.renderEditPrompt()
	}
}

// renderStripPrompt renders the strip text prompt with highlighting
func (m model) renderStripPrompt() string {
	// Split screen: left = input prompt, right = preview
	leftWidth := m.width/2 - 1
	rightWidth := m.width - leftWidth - 2
	panelHeight := m.height - 2 // Leave room for footer

	// Build prompt box content
	promptContent := make([]string, 0)
	promptContent = append(promptContent, "")
	promptContent = append(promptContent, m.inputPrompt)
	promptContent = append(promptContent, "")

	// Input field
	inputVisible := m.inputValue + " "
	promptContent = append(promptContent, "\033[7m"+inputVisible+"\033[0m")

	// Build preview content
	previewContent := m.renderStripPreview()

	// Create bordered panels
	promptPanel := renderBorderedPanel("Input", promptContent, leftWidth, panelHeight, true)
	previewPanel := renderBorderedPanel("Preview", previewContent, rightWidth, panelHeight, false)

	// Combine panels side by side
	var s strings.Builder
	s.WriteString(combinePanelsSideBySide(promptPanel, previewPanel))

	// Footer with help text
	s.WriteString("\033[90m")
	s.WriteString(strings.Repeat("─", m.width))
	s.WriteString("\033[0m\n")
	s.WriteString("\033[90m[\033[33menter\033[90m]apply [\033[33mesc\033[90m]cancel\033[0m")

	return s.String()
}

// renderStripPreview shows files with highlighted text that will be removed
func (m model) renderStripPreview() []string {
	lines := make([]string, 0)

	if len(m.previewFiles) == 0 {
		lines = append(lines, "\033[90mNo files selected\033[0m")
		return lines
	}

	for _, file := range m.previewFiles {
		filename := filepath.Base(file)
		nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

		// Highlight matches if input is not empty
		if m.inputValue != "" && strings.Contains(nameWithoutExt, m.inputValue) {
			// Replace matched text with highlighted version
			highlighted := strings.ReplaceAll(nameWithoutExt, m.inputValue, "\033[41;1m"+m.inputValue+"\033[0m")
			lines = append(lines, highlighted+filepath.Ext(filename))
		} else {
			lines = append(lines, "\033[90m"+filename+"\033[0m")
		}
	}

	if m.inputValue == "" {
		lines = append(lines, "")
		lines = append(lines, "\033[90mType text to see what will be removed\033[0m")
	}

	return lines
}

// renderEditPrompt renders a centered modal for edit operations showing current values
func (m model) renderEditPrompt() string {
	var s strings.Builder

	// Center modal
	modalWidth := 80
	modalHeight := 15
	startY := (m.height - modalHeight) / 2
	startX := (m.width - modalWidth) / 2

	// Read current metadata from first selected file for preview
	var currentValue string
	if len(m.previewFiles) > 0 {
		tags, err := taglib.ReadTags(m.previewFiles[0])
		if err == nil {
			switch m.pendingAction {
			case "edit-title":
				if titles := tags[taglib.Title]; len(titles) > 0 {
					currentValue = titles[0]
				}
			case "edit-artist":
				if artists := tags[taglib.Artist]; len(artists) > 0 {
					currentValue = artists[0]
				}
			case "edit-album":
				if albums := tags[taglib.Album]; len(albums) > 0 {
					currentValue = albums[0]
				}
			case "edit-year":
				if years := tags[taglib.Date]; len(years) > 0 {
					currentValue = years[0]
				}
			}
		}
	}

	// Build modal content
	for i := 0; i < m.height; i++ {
		if i >= startY && i < startY+modalHeight {
			s.WriteString(strings.Repeat(" ", startX))

			if i == startY {
				// Top border
				s.WriteString("\033[1;36m╭" + strings.Repeat("─", modalWidth-2) + "╮\033[0m\n")
			} else if i == startY+modalHeight-1 {
				// Bottom border
				s.WriteString("\033[1;36m╰" + strings.Repeat("─", modalWidth-2) + "╯\033[0m\n")
			} else if i == startY+1 {
				// Title
				s.WriteString("\033[1;36m│\033[0m \033[1m" + m.inputPrompt + "\033[0m")
				padding := modalWidth - runewidth.StringWidth(m.inputPrompt) - 4
				s.WriteString(strings.Repeat(" ", padding) + " \033[1;36m│\033[0m\n")
			} else if i == startY+3 {
				// Current value label
				s.WriteString("\033[1;36m│\033[0m \033[90mCurrent:\033[0m " + currentValue)
				padding := modalWidth - runewidth.StringWidth(currentValue) - 13
				if padding > 0 {
					s.WriteString(strings.Repeat(" ", padding))
				}
				s.WriteString(" \033[1;36m│\033[0m\n")
			} else if i == startY+5 {
				// New value label
				labelText := "New:"
				s.WriteString("\033[1;36m│\033[0m \033[1m" + labelText + "\033[0m")
				padding := modalWidth - runewidth.StringWidth(labelText) - 4
				s.WriteString(strings.Repeat(" ", padding) + " \033[1;36m│\033[0m\n")
			} else if i == startY+6 {
				// Input field
				inputVisible := m.inputValue
				if inputVisible == "" {
					inputVisible = " "
				}
				// Add cursor space to the visible portion in reverse video
				inputWithCursor := inputVisible + " "
				visibleWidth := runewidth.StringWidth(inputWithCursor)

				s.WriteString("\033[1;36m│\033[0m \033[7m")
				s.WriteString(inputWithCursor)
				s.WriteString("\033[0m")

				// Calculate padding: modalWidth - borders(2) - left_space(1) - inputWithCursor - right_space(1)
				padding := modalWidth - 4 - visibleWidth
				if padding > 0 {
					s.WriteString(strings.Repeat(" ", padding))
				}
				s.WriteString(" \033[1;36m│\033[0m\n")
			} else if i == startY+8 {
				// File count
				fileCount := fmt.Sprintf("\033[90mApplying to %d file(s)\033[0m", len(m.previewFiles))
				s.WriteString("\033[1;36m│\033[0m ")
				s.WriteString(fileCount)
				padding := modalWidth - runewidth.StringWidth(stripAnsi(fileCount)) - 4
				if padding > 0 {
					s.WriteString(strings.Repeat(" ", padding))
				}
				s.WriteString(" \033[1;36m│\033[0m\n")
			} else {
				// Empty line
				s.WriteString("\033[1;36m│\033[0m")
				s.WriteString(strings.Repeat(" ", modalWidth-2))
				s.WriteString("\033[1;36m│\033[0m\n")
			}
		} else if i == startY+modalHeight {
			// Add help text right below modal
			s.WriteString(strings.Repeat(" ", startX+2))
			s.WriteString("\033[90m[\033[33menter\033[90m]apply [\033[33mesc\033[90m]cancel\033[0m\n")
		} else {
			s.WriteString("\n")
		}
	}

	return s.String()
}

// renderFileBrowser renders the file browser modal
func (m model) renderFileBrowser() string {
	var s strings.Builder

	// Center modal
	modalWidth := 80
	modalHeight := 25
	startY := (m.height - modalHeight) / 2
	startX := (m.width - modalWidth) / 2

	// Build content list
	content := make([]string, 0)
	content = append(content, "\033[90mCurrent: \033[0m"+m.browserDir)
	content = append(content, "")

	// Add directories
	for i, dir := range m.browserDirs {
		idx := i
		var line strings.Builder

		if idx == m.browserCursor {
			line.WriteString("\033[36m> \033[0m")
		} else {
			line.WriteString("  ")
		}

		line.WriteString("\033[34m📁 ")
		if idx == m.browserCursor {
			line.WriteString("\033[1m")
		}
		line.WriteString(dir)
		if idx == m.browserCursor {
			line.WriteString("\033[0m")
		} else {
			line.WriteString("\033[0m")
		}

		content = append(content, line.String())
	}

	// Add files
	for i, file := range m.browserFiles {
		idx := i + len(m.browserDirs)
		var line strings.Builder

		if idx == m.browserCursor {
			line.WriteString("\033[36m> \033[0m")
		} else {
			line.WriteString("  ")
		}

		line.WriteString("\033[32m🖼️  ")
		if idx == m.browserCursor {
			line.WriteString("\033[1m")
		}
		line.WriteString(file)
		if idx == m.browserCursor {
			line.WriteString("\033[0m")
		} else {
			line.WriteString("\033[0m")
		}

		content = append(content, line.String())
	}

	if len(m.browserDirs)+len(m.browserFiles) == 0 {
		content = append(content, "\033[90mNo image files found\033[0m")
	}

	// Build modal
	for i := 0; i < m.height; i++ {
		if i >= startY && i < startY+modalHeight {
			s.WriteString(strings.Repeat(" ", startX))

			if i == startY {
				// Top border
				s.WriteString("\033[1;36m╭─ \033[1mSelect Cover Image\033[0m\033[1;36m ")
				titleLen := len("Select Cover Image") + 3
				remainingWidth := modalWidth - titleLen - 2
				s.WriteString(strings.Repeat("─", remainingWidth))
				s.WriteString("╮\033[0m\n")
			} else if i == startY+modalHeight-1 {
				// Bottom border
				s.WriteString("\033[1;36m╰" + strings.Repeat("─", modalWidth-2) + "╯\033[0m\n")
			} else {
				// Content
				lineIdx := i - startY - 1
				s.WriteString("\033[1;36m│\033[0m ")

				if lineIdx < len(content) {
					line := content[lineIdx]
					lineWidth := runewidth.StringWidth(stripAnsi(line))
					s.WriteString(line)
					padding := modalWidth - 4 - lineWidth
					if padding > 0 {
						s.WriteString(strings.Repeat(" ", padding))
					}
				} else {
					s.WriteString(strings.Repeat(" ", modalWidth-4))
				}

				s.WriteString(" \033[1;36m│\033[0m\n")
			}
		} else if i == startY+modalHeight {
			// Help text below modal
			s.WriteString(strings.Repeat(" ", startX+2))
			s.WriteString("\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]select [\033[33mesc\033[90m]cancel\033[0m\n")
		} else {
			s.WriteString("\n")
		}
	}

	return s.String()
}
// renderPreview renders the Navidrome preview
func (m model) renderPreview() string {
	var s strings.Builder

	// Two column layout: Left (cover + info) | Right (tracklist)
	leftWidth := 70  // Fixed width for cover and info
	rightWidth := m.width - leftWidth - 2
	panelHeight := m.height - 2

	// Build left panel content (cover + album info)
	leftContent := make([]string, 0)

	if len(m.previewAlbums) == 0 {
		leftContent = append(leftContent, "\033[90mNo albums found\033[0m")
	} else if m.previewCursor < len(m.previewAlbums) {
		album := m.previewAlbums[m.previewCursor]

		// Show cover art if present
		if album.HasCover && len(album.Tracks) > 0 {
			if coverData, err := taglib.ReadImage(album.Tracks[0].FilePath); err == nil && len(coverData) > 0 {
				coverLines := renderCoverArtPreview(coverData)
				leftContent = append(leftContent, coverLines...)
				leftContent = append(leftContent, "")
			}
		}

		// Album info below cover
		leftContent = append(leftContent, "\033[1;36mAlbum:\033[0m \033[1m"+album.AlbumName+"\033[0m")
		leftContent = append(leftContent, "\033[1;36mArtist:\033[0m "+album.AlbumArtist)

		// Year with empty indicator
		if album.Year != "" {
			leftContent = append(leftContent, "\033[1;36mYear:\033[0m "+album.Year)
		} else {
			leftContent = append(leftContent, "\033[1;36mYear:\033[0m \033[90m<empty>\033[0m")
		}

		// Genre with empty indicator
		if album.Genre != "" {
			leftContent = append(leftContent, "\033[1;36mGenre:\033[0m "+album.Genre)
		} else {
			leftContent = append(leftContent, "\033[1;36mGenre:\033[0m \033[90m<empty>\033[0m")
		}

		leftContent = append(leftContent, "")
		leftContent = append(leftContent, "\033[90mTotal: "+strconv.Itoa(len(album.Tracks))+" track(s)\033[0m")
	}

	// Build right panel content (tracklist)
	trackListContent := make([]string, 0)
	if m.previewCursor < len(m.previewAlbums) {
		album := m.previewAlbums[m.previewCursor]

		for _, track := range album.Tracks {
			var trackLine strings.Builder

			// Track number
			if track.TrackNumber != "" {
				trackLine.WriteString("\033[36m")
				trackLine.WriteString(track.TrackNumber)
				trackLine.WriteString(".\033[0m ")
			} else {
				trackLine.WriteString("\033[90m-.\033[0m ")
			}

			// Title
			trackLine.WriteString(track.Title)

			// Artist if different from album artist
			if track.Artist != album.AlbumArtist {
				trackLine.WriteString(" \033[90m(")
				trackLine.WriteString(track.Artist)
				trackLine.WriteString(")\033[0m")
			}

			trackListContent = append(trackListContent, trackLine.String())
		}
	}

	// Create bordered panels
	leftPanel := renderBorderedPanel("Album", leftContent, leftWidth, panelHeight, true)
	rightPanel := renderBorderedPanel("Tracks", trackListContent, rightWidth, panelHeight, false)

	// Combine panels
	s.WriteString(combinePanelsSideBySide(leftPanel, rightPanel))

	// Footer
	s.WriteString("\033[90m")
	s.WriteString(strings.Repeat("─", m.width))
	s.WriteString("\033[0m\n")
	s.WriteString("\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33mp/esc\033[90m]close • \033[36mThis shows how albums will appear in Navidrome\033[0m")

	return s.String()
}
