package main

// TreeNode represents a folder in the tree
type TreeNode struct {
	Path     string
	Name     string
	Expanded bool
	Children []*TreeNode
	Files    []string
	Parent   *TreeNode
	Depth    int
}

// model represents the application state
type model struct {
	rootDir         string
	tree            *TreeNode
	flattenedTree   []*TreeNode
	cursor          int
	selectedFiles   map[string]bool
	width           int
	height          int
	inFilePanel     bool
	fileCursor      int
	showEditMenu    bool
	editCursor      int
	showInputPrompt bool
	inputValue      string
	inputPrompt     string
	pendingAction   string // "strip", "edit-title", etc.
	errorMessage    string
	previewFiles    []string // Files to preview during input
	currentNodePath string   // Store current folder path to restore state
}
