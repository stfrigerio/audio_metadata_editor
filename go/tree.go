package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// buildTree recursively builds a tree structure from a directory path
func buildTree(path string, parent *TreeNode, depth int) *TreeNode {
	node := &TreeNode{
		Path:     path,
		Name:     filepath.Base(path),
		Expanded: depth == 0, // Root is expanded by default
		Parent:   parent,
		Depth:    depth,
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return node
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(path, name)

		if entry.IsDir() {
			child := buildTree(fullPath, node, depth+1)
			node.Children = append(node.Children, child)
		} else if isAudioFile(name) {
			node.Files = append(node.Files, fullPath)
		}
	}

	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Name < node.Children[j].Name
	})
	sort.Strings(node.Files)

	return node
}

// flattenTree converts a tree structure into a flat list of visible nodes
func flattenTree(node *TreeNode, result *[]*TreeNode) {
	*result = append(*result, node)
	if node.Expanded {
		for _, child := range node.Children {
			flattenTree(child, result)
		}
	}
}
