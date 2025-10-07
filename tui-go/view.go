package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mattn/go-runewidth"
	taglib "go.senan.xyz/taglib"
)

// View renders the current state of the application
func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Show input prompt overlay if active
	if m.showInputPrompt {
		return m.renderInputPrompt()
	}

	// Get current node for file display
	currentNode := m.getCurrentNode()

	// Calculate panel dimensions (with borders)
	leftWidth := m.width/2 - 1
	rightWidth := m.width - leftWidth - 3
	panelHeight := m.height - 3 // Leave room for footer

	return m.renderBorderedLayout(currentNode, leftWidth, rightWidth, panelHeight)
}

// renderBorderedLayout renders the main layout with bordered panels
func (m model) renderBorderedLayout(currentNode *TreeNode, leftWidth, rightWidth, panelHeight int) string {
	var s strings.Builder

	// Get content for both panels
	treeContent := m.renderTreeContent()
	filesContent := m.renderFilesContent(currentNode)

	// Determine titles and focus
	treeTitle := "Files"
	filesTitle := "Details"
	treeFocused := !m.inFilePanel
	filesFocused := m.inFilePanel

	// Create bordered boxes
	leftBox := renderBorderedBox(treeTitle, treeContent, leftWidth, panelHeight, treeFocused)
	rightBox := renderBorderedBox(filesTitle, filesContent, rightWidth, panelHeight, filesFocused)

	// Render side by side
	s.WriteString(renderSideBySide(leftBox, rightBox))

	// If edit menu is shown, render it as an overlay at the bottom left
	if m.showEditMenu {
		// Three panels: left (tree), bottom left (edit menu overlay), right (files)
		leftWidth := m.width / 3
		rightWidth := m.width - leftWidth - 1

		leftLines = m.renderTree()
		editLines := m.renderEditMenu()
		rightLines = m.renderFiles(currentNode)

		availHeight := m.height - 2
		editHeight := len(editLines)
		treeOnlyHeight := availHeight - editHeight

		// Render side by side
		for i := 0; i < availHeight; i++ {
			if i < treeOnlyHeight {
				// Top left: tree only
				if i < len(leftLines) {
					line := leftLines[i]
					lineWidth := runewidth.StringWidth(stripAnsi(line))
					if lineWidth > leftWidth {
						stripped := stripAnsi(line)
						truncated := runewidth.Truncate(stripped, leftWidth, "")
						line = truncated
						lineWidth = leftWidth
					}
					s.WriteString(line)
					s.WriteString(strings.Repeat(" ", leftWidth-lineWidth))
				} else {
					s.WriteString(strings.Repeat(" ", leftWidth))
				}
			} else {
				// Bottom left: edit menu (overlaid on tree)
				lineIdx := i - treeOnlyHeight
				if lineIdx < len(editLines) {
					line := editLines[lineIdx]
					lineWidth := runewidth.StringWidth(stripAnsi(line))
					if lineWidth > leftWidth {
						stripped := stripAnsi(line)
						truncated := runewidth.Truncate(stripped, leftWidth, "")
						line = truncated
						lineWidth = leftWidth
					}
					s.WriteString(line)
					s.WriteString(strings.Repeat(" ", leftWidth-lineWidth))
				} else {
					s.WriteString(strings.Repeat(" ", leftWidth))
				}
			}

			s.WriteString("│")

			// Right side (files) - full height
			if i < len(rightLines) {
				line := rightLines[i]
				lineWidth := runewidth.StringWidth(stripAnsi(line))
				if lineWidth > rightWidth {
					stripped := stripAnsi(line)
					truncated := runewidth.Truncate(stripped, rightWidth, "")
					line = truncated
				}
				s.WriteString(line)
			}

			s.WriteString("\n")
		}
	} else {
		// Normal view: left panel (tree), right panel (files)
		leftWidth := m.width / 2
		rightWidth := m.width - leftWidth - 1

		leftLines = m.renderTree()
		rightLines = m.renderFiles(currentNode)

		// Calculate max lines
		maxLines := len(leftLines)
		if len(rightLines) > maxLines {
			maxLines = len(rightLines)
		}

		// Render side by side
		for i := 0; i < maxLines && i < m.height-2; i++ {
			// Left side (tree)
			if i < len(leftLines) {
				line := leftLines[i]
				lineWidth := runewidth.StringWidth(stripAnsi(line))
				if lineWidth > leftWidth {
					stripped := stripAnsi(line)
					truncated := runewidth.Truncate(stripped, leftWidth, "")
					line = truncated
					lineWidth = leftWidth
				}
				s.WriteString(line)
				s.WriteString(strings.Repeat(" ", leftWidth-lineWidth))
			} else {
				s.WriteString(strings.Repeat(" ", leftWidth))
			}

			s.WriteString("│")

			// Right side (files)
			if i < len(rightLines) {
				line := rightLines[i]
				lineWidth := runewidth.StringWidth(stripAnsi(line))
				if lineWidth > rightWidth {
					stripped := stripAnsi(line)
					truncated := runewidth.Truncate(stripped, rightWidth, "")
					line = truncated
				}
				s.WriteString(line)
			}

			s.WriteString("\n")
		}
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

// renderTree generates the left panel tree view
func (m model) renderTree() []string {
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

// renderFiles generates the file list panel
func (m model) renderFiles(node *TreeNode) []string {
	lines := make([]string, 0)

	if node == nil {
		lines = append(lines, "\033[90m  No folder selected\033[0m")
		return lines
	}

	if len(node.Files) == 0 {
		lines = append(lines, "\033[90m  No audio files\033[0m")
		return lines
	}

	lines = append(lines, fmt.Sprintf("\033[1mFiles in: \033[36m%s\033[0m", node.Name))
	lines = append(lines, "")

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

// renderEditMenu generates the edit menu panel
func (m model) renderEditMenu() []string {
	lines := make([]string, 0)

	lines = append(lines, "\033[1mEdit Options\033[0m")
	lines = append(lines, "")

	options := []string{
		"1. Strip text from titles",
		"2. Edit Title",
		"3. Edit Artist",
		"4. Edit Album",
		"5. Edit Year",
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
			return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]toggle [\033[33ma\033[90m]toggle all [\033[33me\033[90m]edit [\033[33mtab\033[90m]tree [\033[33mq\033[90m]quit\033[0m"
		} else {
			return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]toggle [\033[33ma\033[90m]toggle all [\033[33mtab\033[90m]tree [\033[33mq\033[90m]quit\033[0m"
		}
	} else {
		return "\033[90m[\033[33m↑/↓\033[90m]navigate [\033[33menter\033[90m]open/toggle [\033[33mtab\033[90m]files [\033[33mq\033[90m]quit\033[0m"
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
	var s strings.Builder

	// Split screen: left = input prompt, right = preview
	leftWidth := 50
	rightWidth := m.width - leftWidth - 1

	// Build preview lines
	previewLines := m.renderStripPreview()

	// Build prompt lines
	promptLines := make([]string, 0)
	promptLines = append(promptLines, "\033[1;36m╭"+strings.Repeat("─", leftWidth-2)+"╮\033[0m")
	promptLines = append(promptLines, "\033[1;36m│\033[0m "+m.inputPrompt+strings.Repeat(" ", leftWidth-len(m.inputPrompt)-4)+" \033[1;36m│\033[0m")
	promptLines = append(promptLines, "\033[1;36m│\033[0m"+strings.Repeat(" ", leftWidth-2)+"\033[1;36m│\033[0m")

	// Input field - properly account for visible width only
	inputVisible := m.inputValue + " " // visible characters only
	visibleWidth := runewidth.StringWidth(inputVisible)

	// Build the line with proper spacing
	var inputDisplay strings.Builder
	inputDisplay.WriteString("\033[1;36m│\033[0m ")
	inputDisplay.WriteString("\033[7m")
	inputDisplay.WriteString(inputVisible)
	inputDisplay.WriteString("\033[0m")

	// Add padding to reach the right border
	paddingNeeded := leftWidth - visibleWidth - 4 // -4 for "│ " and " │"
	if paddingNeeded > 0 {
		inputDisplay.WriteString(strings.Repeat(" ", paddingNeeded))
	}
	inputDisplay.WriteString(" \033[1;36m│\033[0m")

	promptLines = append(promptLines, inputDisplay.String())

	promptLines = append(promptLines, "\033[1;36m│\033[0m"+strings.Repeat(" ", leftWidth-2)+"\033[1;36m│\033[0m")
	helpText := "\033[90m[enter]apply [esc]cancel\033[0m"
	helpLen := runewidth.StringWidth(stripAnsi(helpText))
	promptLines = append(promptLines, "\033[1;36m│\033[0m "+helpText+strings.Repeat(" ", leftWidth-helpLen-4)+" \033[1;36m│\033[0m")
	promptLines = append(promptLines, "\033[1;36m╰"+strings.Repeat("─", leftWidth-2)+"╯\033[0m")

	// Render side by side
	maxLines := len(promptLines)
	if len(previewLines) > maxLines {
		maxLines = len(previewLines)
	}

	for i := 0; i < m.height-1; i++ {
		// Left side (prompt)
		if i < len(promptLines) {
			line := promptLines[i]
			lineWidth := runewidth.StringWidth(stripAnsi(line))
			s.WriteString(line)
			if lineWidth < leftWidth {
				s.WriteString(strings.Repeat(" ", leftWidth-lineWidth))
			}
		} else {
			s.WriteString(strings.Repeat(" ", leftWidth))
		}

		s.WriteString("│")

		// Right side (preview)
		if i < len(previewLines) {
			line := previewLines[i]
			lineWidth := runewidth.StringWidth(stripAnsi(line))
			if lineWidth > rightWidth {
				stripped := stripAnsi(line)
				truncated := runewidth.Truncate(stripped, rightWidth, "")
				line = truncated
			}
			s.WriteString(line)
		}

		s.WriteString("\n")
	}

	return s.String()
}

// renderStripPreview shows files with highlighted text that will be removed
func (m model) renderStripPreview() []string {
	lines := make([]string, 0)

	lines = append(lines, "\033[1mPreview - Text to remove highlighted:\033[0m")
	lines = append(lines, "")

	if len(m.previewFiles) == 0 {
		lines = append(lines, "\033[90m  No files selected\033[0m")
		return lines
	}

	for _, file := range m.previewFiles {
		filename := filepath.Base(file)
		nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

		// Highlight matches if input is not empty
		if m.inputValue != "" && strings.Contains(nameWithoutExt, m.inputValue) {
			// Replace matched text with highlighted version
			highlighted := strings.ReplaceAll(nameWithoutExt, m.inputValue, "\033[41;1m"+m.inputValue+"\033[0m")
			lines = append(lines, "  "+highlighted+filepath.Ext(filename))
		} else {
			lines = append(lines, "  \033[90m"+filename+"\033[0m")
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
				padding := modalWidth - len(m.inputPrompt) - 4
				s.WriteString(strings.Repeat(" ", padding) + " \033[1;36m│\033[0m\n")
			} else if i == startY+3 {
				// Current value label
				s.WriteString("\033[1;36m│\033[0m \033[90mCurrent:\033[0m " + currentValue)
				padding := modalWidth - len(currentValue) - 13
				if padding > 0 {
					s.WriteString(strings.Repeat(" ", padding))
				}
				s.WriteString(" \033[1;36m│\033[0m\n")
			} else if i == startY+5 {
				// New value label
				s.WriteString("\033[1;36m│\033[0m \033[1mNew:\033[0m")
				s.WriteString(strings.Repeat(" ", modalWidth-9) + " \033[1;36m│\033[0m\n")
			} else if i == startY+6 {
				// Input field
				inputVisible := m.inputValue + " "
				visibleWidth := runewidth.StringWidth(inputVisible)

				s.WriteString("\033[1;36m│\033[0m  \033[7m")
				s.WriteString(inputVisible)
				s.WriteString("\033[0m")

				padding := modalWidth - visibleWidth - 6
				if padding > 0 {
					s.WriteString(strings.Repeat(" ", padding))
				}
				s.WriteString("  \033[1;36m│\033[0m\n")
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
			} else if i == startY+modalHeight-2 {
				// Help text
				helpText := "\033[90m[enter]apply [esc]cancel\033[0m"
				helpLen := runewidth.StringWidth(stripAnsi(helpText))
				s.WriteString("\033[1;36m│\033[0m ")
				s.WriteString(helpText)
				padding := modalWidth - helpLen - 4
				s.WriteString(strings.Repeat(" ", padding) + " \033[1;36m│\033[0m\n")
			} else {
				// Empty line
				s.WriteString("\033[1;36m│\033[0m")
				s.WriteString(strings.Repeat(" ", modalWidth-2))
				s.WriteString("\033[1;36m│\033[0m\n")
			}
		} else {
			s.WriteString("\n")
		}
	}

	return s.String()
}
