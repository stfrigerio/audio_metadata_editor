package main

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// renderBorderedPanel renders content inside a bordered box with rounded corners
func renderBorderedPanel(title string, content []string, width, height int, focused bool) []string {
	lines := make([]string, 0, height)

	// Border color based on focus
	borderColor := "\033[90m" // Gray for unfocused
	if focused {
		borderColor = "\033[36m" // Cyan for focused
	}

	// Top border with title
	topBorder := borderColor + "╭─"
	if title != "" {
		topBorder += " \033[1m" + title + "\033[0m" + borderColor + " "
		titleLen := len(title) + 3
		remainingWidth := width - titleLen - 2
		if remainingWidth > 0 {
			topBorder += strings.Repeat("─", remainingWidth)
		}
	} else {
		repeatCount := width - 3
		if repeatCount > 0 {
			topBorder += strings.Repeat("─", repeatCount)
		}
	}
	topBorder += "╮\033[0m"
	lines = append(lines, topBorder)

	// Content lines (height - 2 for top and bottom borders)
	contentHeight := height - 2
	for i := 0; i < contentHeight; i++ {
		line := borderColor + "│\033[0m "

		if i < len(content) {
			contentLine := content[i]
			lineWidth := runewidth.StringWidth(stripAnsi(contentLine))

			// Available width for content (width - borders - padding)
			availWidth := width - 4

			if lineWidth > availWidth {
				stripped := stripAnsi(contentLine)
				truncated := runewidth.Truncate(stripped, availWidth, "")
				contentLine = truncated
				lineWidth = availWidth
			}

			line += contentLine
			padding := availWidth - lineWidth
			if padding > 0 {
				line += strings.Repeat(" ", padding)
			}
		} else {
			emptyWidth := width - 4
			if emptyWidth > 0 {
				line += strings.Repeat(" ", emptyWidth)
			}
		}

		line += " " + borderColor + "│\033[0m"
		lines = append(lines, line)
	}

	// Bottom border
	bottomWidth := width - 2
	if bottomWidth < 0 {
		bottomWidth = 0
	}
	bottomBorder := borderColor + "╰" + strings.Repeat("─", bottomWidth) + "╯\033[0m"
	lines = append(lines, bottomBorder)

	return lines
}

// combinePanelsSideBySide combines two bordered panels side by side
func combinePanelsSideBySide(left, right []string) string {
	var s strings.Builder

	maxLines := len(left)
	if len(right) > maxLines {
		maxLines = len(right)
	}

	// Calculate left panel width from first line (if it exists)
	leftWidth := 0
	if len(left) > 0 {
		leftWidth = runewidth.StringWidth(stripAnsi(left[0]))
	}

	for i := 0; i < maxLines; i++ {
		// Left panel
		if i < len(left) {
			s.WriteString(left[i])
		} else {
			// Pad with spaces if left panel is shorter
			s.WriteString(strings.Repeat(" ", leftWidth))
		}

		// Single space separator
		s.WriteString(" ")

		// Right panel
		if i < len(right) {
			s.WriteString(right[i])
		}

		s.WriteString("\n")
	}

	return s.String()
}
