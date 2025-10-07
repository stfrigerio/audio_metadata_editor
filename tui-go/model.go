package main

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// initialModel creates the initial application state
func initialModel(dir string) model {
	absDir, _ := filepath.Abs(dir)
	tree := buildTree(absDir, nil, 0)

	flattened := make([]*TreeNode, 0)
	flattenTree(tree, &flattened)

	return model{
		rootDir:         absDir,
		tree:            tree,
		flattenedTree:   flattened,
		cursor:          0,
		selectedFiles:   make(map[string]bool),
		inFilePanel:     false,
		fileCursor:      0,
		showEditMenu:    false,
		editCursor:      0,
		showInputPrompt: false,
		inputValue:      "",
		inputPrompt:     "",
		pendingAction:   "",
		errorMessage:    "",
	}
}

// Init initializes the BubbleTea program
func (m model) Init() tea.Cmd {
	return nil
}

// refreshFlattenedTree rebuilds the flattened tree view
func (m *model) refreshFlattenedTree() {
	m.flattenedTree = make([]*TreeNode, 0)
	flattenTree(m.tree, &m.flattenedTree)
}

// getCurrentNode returns the currently selected tree node
func (m *model) getCurrentNode() *TreeNode {
	if m.cursor >= 0 && m.cursor < len(m.flattenedTree) {
		return m.flattenedTree[m.cursor]
	}
	return nil
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle preview mode first
		if m.showPreview {
			m.handlePreviewKey(msg.String())
			return m, nil
		}

		// Handle file browser
		if m.showFileBrowser {
			m.handleBrowserKey(msg.String())
			return m, nil
		}

		// Handle input prompt
		if m.showInputPrompt {
			m.handleInputPrompt(msg.String())
			return m, nil
		}

		key := msg.String()

		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.showEditMenu {
				m.showEditMenu = false
				m.editCursor = 0
			}
			if m.errorMessage != "" {
				m.errorMessage = ""
			}

		case "tab":
			if !m.showEditMenu {
				node := m.getCurrentNode()
				if node != nil && len(node.Files) > 0 {
					m.inFilePanel = !m.inFilePanel
					if m.inFilePanel {
						m.fileCursor = 0
					}
				}
			}

		case "down":
			m.handleDown()

		case "up":
			m.handleUp()

		case "e":
			if m.inFilePanel && len(m.selectedFiles) > 0 {
				m.showEditMenu = true
				m.editCursor = 0
			}

		case "enter", " ":
			m.handleEnter()

		case "a":
			m.handleToggleAll()

		case "p":
			// Toggle preview mode
			node := m.getCurrentNode()
			if node != nil && len(node.Files) > 0 {
				if !m.showPreview {
					m.buildAlbumPreviews()
					m.showPreview = true
				} else {
					m.showPreview = false
					m.previewAlbums = nil
					m.previewCursor = 0
				}
			}
		}
	}

	return m, nil
}

// handleDown handles the down arrow key
func (m *model) handleDown() {
	if m.showEditMenu {
		m.editCursor++
		if m.editCursor >= 6 { // 6 edit options
			m.editCursor = 5
		}
	} else if m.inFilePanel {
		node := m.getCurrentNode()
		if node != nil {
			m.fileCursor++
			if m.fileCursor >= len(node.Files) {
				m.fileCursor = len(node.Files) - 1
			}
		}
	} else {
		m.cursor++
		if m.cursor >= len(m.flattenedTree) {
			m.cursor = len(m.flattenedTree) - 1
		}
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.fileCursor < 0 {
		m.fileCursor = 0
	}
}

// handleUp handles the up arrow key
func (m *model) handleUp() {
	if m.showEditMenu {
		m.editCursor--
		if m.editCursor < 0 {
			m.editCursor = 0
		}
	} else if m.inFilePanel {
		m.fileCursor--
		if m.fileCursor < 0 {
			m.fileCursor = 0
		}
	} else {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
}

// handleEnter handles the enter key
func (m *model) handleEnter() {
	if m.showEditMenu {
		// Execute the selected edit operation
		switch m.editCursor {
		case 0: // Strip text from titles
			// Store current state
			node := m.getCurrentNode()
			if node != nil {
				m.currentNodePath = node.Path
			}

			// Prepare preview files
			m.previewFiles = make([]string, 0, len(m.selectedFiles))
			for file := range m.selectedFiles {
				m.previewFiles = append(m.previewFiles, file)
			}

			m.pendingAction = "strip"
			m.inputPrompt = "Enter text to remove:"
			m.inputValue = ""
			m.showInputPrompt = true
			m.showEditMenu = false
		case 1: // Add cover image
			// Open file browser instead of text input
			m.openFileBrowser()
			m.showEditMenu = false
		case 2: // Edit Title
			m.startEditAction("edit-title", "Enter new title:")
		case 3: // Edit Artist
			m.startEditAction("edit-artist", "Enter new artist:")
		case 4: // Edit Album
			m.startEditAction("edit-album", "Enter new album:")
		case 5: // Edit Year
			m.startEditAction("edit-year", "Enter new year:")
		}
	} else if m.inFilePanel {
		// Toggle file selection in file panel
		node := m.getCurrentNode()
		if node != nil && m.fileCursor < len(node.Files) {
			file := node.Files[m.fileCursor]
			m.selectedFiles[file] = !m.selectedFiles[file]
		}
	} else {
		// In tree panel
		node := m.getCurrentNode()
		if node != nil {
			// If folder has files, navigate to file panel
			if len(node.Files) > 0 {
				m.inFilePanel = true
				m.fileCursor = 0
			} else if len(node.Children) > 0 {
				// Otherwise toggle expand/collapse if it has subfolders
				node.Expanded = !node.Expanded
				m.refreshFlattenedTree()
			}
		}
	}
}

// startEditAction starts an edit action with the given type and prompt
func (m *model) startEditAction(actionType string, prompt string) {
	// Store current state
	node := m.getCurrentNode()
	if node != nil {
		m.currentNodePath = node.Path
	}

	// Prepare preview files
	m.previewFiles = make([]string, 0, len(m.selectedFiles))
	for file := range m.selectedFiles {
		m.previewFiles = append(m.previewFiles, file)
	}

	m.pendingAction = actionType
	m.inputPrompt = prompt
	m.inputValue = ""
	m.showInputPrompt = true
	m.showEditMenu = false
}

// handleToggleAll handles the 'a' key for selecting/deselecting all files
func (m *model) handleToggleAll() {
	if m.inFilePanel && !m.showEditMenu {
		node := m.getCurrentNode()
		if node != nil {
			// Check if all are selected
			allSelected := true
			for _, file := range node.Files {
				if !m.selectedFiles[file] {
					allSelected = false
					break
				}
			}

			// Toggle: if all selected, deselect all; otherwise select all
			if allSelected {
				for _, file := range node.Files {
					delete(m.selectedFiles, file)
				}
			} else {
				for _, file := range node.Files {
					m.selectedFiles[file] = true
				}
			}
		}
	}
}

// handleInputPrompt handles keyboard input when the input prompt is shown
func (m *model) handleInputPrompt(key string) {
	switch key {
	case "esc":
		// Cancel input
		m.showInputPrompt = false
		m.inputValue = ""
		m.pendingAction = ""
		m.previewFiles = nil
		m.currentNodePath = ""

	case "enter":
		// Execute the pending action with the input value
		// Note: executeAction will close the prompt via defer
		m.executeAction()

	case "backspace":
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}

	default:
		// Add character to input (filter out special keys)
		if len(key) == 1 {
			m.inputValue += key
		}
	}
}

// executeAction executes the pending action with the current input value
func (m *model) executeAction() {
	// Always close the input prompt, even if execution fails
	defer func() {
		m.showInputPrompt = false
		m.pendingAction = ""
		m.inputValue = ""
		m.previewFiles = nil
	}()

	if m.inputValue == "" {
		m.errorMessage = "Input cannot be empty"
		return
	}

	// Get list of selected files
	selectedFiles := make([]string, 0, len(m.selectedFiles))
	for file := range m.selectedFiles {
		selectedFiles = append(selectedFiles, file)
	}

	var errors []error

	switch m.pendingAction {
	case "strip":
		errors = StripTextFromFiles(selectedFiles, m.inputValue)
	case "add-cover":
		errors = AddCoverImageToFiles(selectedFiles, m.inputValue)
	case "edit-title":
		errors = EditTitleForFiles(selectedFiles, m.inputValue)
	case "edit-artist":
		errors = EditArtistForFiles(selectedFiles, m.inputValue)
	case "edit-album":
		errors = EditAlbumForFiles(selectedFiles, m.inputValue)
	case "edit-year":
		errors = EditYearForFiles(selectedFiles, m.inputValue)
	}

	if len(errors) > 0 {
		m.errorMessage = errors[0].Error()
	} else {
		// Refresh the tree (only needed for strip which renames files)
		if m.pendingAction == "strip" {
			m.tree = buildTree(m.rootDir, nil, 0)
			m.refreshFlattenedTree()
		}

		// Restore the folder state
		m.restoreFolderState()

		// Clear selection since files may have been renamed
		m.selectedFiles = make(map[string]bool)

		// If no files remain, switch back to tree panel
		node := m.getCurrentNode()
		if node == nil || len(node.Files) == 0 {
			m.inFilePanel = false
		}

		m.errorMessage = ""
	}
}

// restoreFolderState restores the tree to show the previously opened folder
func (m *model) restoreFolderState() {
	if m.currentNodePath == "" {
		return
	}

	// Find the node in the tree (recursively) and expand all parents
	targetNode := m.findAndExpandNode(m.tree, m.currentNodePath)
	if targetNode != nil {
		// Rebuild flattened tree with expanded nodes
		m.refreshFlattenedTree()

		// Find the node in the flattened tree
		for i, node := range m.flattenedTree {
			if node.Path == m.currentNodePath {
				m.cursor = i
				m.inFilePanel = true
				m.fileCursor = 0
				break
			}
		}
	}

	m.currentNodePath = ""
}

// findAndExpandNode recursively finds a node by path and expands all parents
func (m *model) findAndExpandNode(node *TreeNode, targetPath string) *TreeNode {
	if node.Path == targetPath {
		return node
	}

	// Search in children
	for _, child := range node.Children {
		if result := m.findAndExpandNode(child, targetPath); result != nil {
			// Expand this node since it's a parent of the target
			node.Expanded = true
			return result
		}
	}

	return nil
}
