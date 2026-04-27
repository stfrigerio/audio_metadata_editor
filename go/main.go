package main

import (
	"fmt"
	"io/fs"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: audio-metadata-editor-tui <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Directory does not exist: %s\n", dir)
		} else if errors, ok := err.(*fs.PathError); ok {
			fmt.Printf("Cannot access directory: %s (%v)\n", dir, errors)
		} else {
			fmt.Printf("Error accessing directory: %s (%v)\n", dir, err)
		}
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(dir), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
