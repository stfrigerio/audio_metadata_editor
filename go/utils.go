package main

import (
	"path/filepath"
	"strings"
)

var supportedExtensions = map[string]bool{
	".mp3":  true,
	".m4a":  true,
	".flac": true,
	".ogg":  true,
}

// isAudioFile checks if a filename has a supported audio extension
func isAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return supportedExtensions[ext]
}

// stripAnsi removes ANSI escape codes from a string
func stripAnsi(str string) string {
	result := str
	for {
		start := strings.Index(result, "\033[")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "m")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return result
}
