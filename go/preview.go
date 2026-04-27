package main

import (
	"path/filepath"
	"sort"
	"strconv"

	taglib "go.senan.xyz/taglib"
)

// buildAlbumPreviews analyzes metadata and groups files as Navidrome would
func (m *model) buildAlbumPreviews() {
	node := m.getCurrentNode()
	if node == nil || len(node.Files) == 0 {
		m.previewAlbums = nil
		return
	}

	// Map to group files by album
	albumMap := make(map[string]*AlbumPreview)

	for _, filePath := range node.Files {
		tags, err := taglib.ReadTags(filePath)
		if err != nil {
			continue
		}

		// Extract metadata
		album := getTagValue(tags, taglib.Album, "Unknown Album")
		albumArtist := getTagValue(tags, taglib.AlbumArtist, "")
		if albumArtist == "" {
			albumArtist = getTagValue(tags, taglib.Artist, "Unknown Artist")
		}
		artist := getTagValue(tags, taglib.Artist, "Unknown Artist")
		title := getTagValue(tags, taglib.Title, filepath.Base(filePath))
		year := getTagValue(tags, taglib.Date, "")
		genre := getTagValue(tags, taglib.Genre, "")
		trackNum := getTagValue(tags, taglib.TrackNumber, "")

		// Create unique key for album grouping (Navidrome groups by Album + AlbumArtist)
		albumKey := album + "|" + albumArtist

		// Get or create album preview
		if _, exists := albumMap[albumKey]; !exists {
			// Check if file has cover art
			hasCover := false
			if coverData, err := taglib.ReadImage(filePath); err == nil && len(coverData) > 0 {
				hasCover = true
			}

			albumMap[albumKey] = &AlbumPreview{
				AlbumName:   album,
				AlbumArtist: albumArtist,
				Year:        year,
				Genre:       genre,
				Tracks:      make([]TrackPreview, 0),
				HasCover:    hasCover,
			}
		}

		// Add track to album
		albumMap[albumKey].Tracks = append(albumMap[albumKey].Tracks, TrackPreview{
			TrackNumber: trackNum,
			Title:       title,
			Artist:      artist,
			FilePath:    filePath,
		})
	}

	// Convert map to sorted slice
	albums := make([]AlbumPreview, 0, len(albumMap))
	for _, album := range albumMap {
		// Sort tracks by track number
		sort.Slice(album.Tracks, func(i, j int) bool {
			numI, errI := strconv.Atoi(album.Tracks[i].TrackNumber)
			numJ, errJ := strconv.Atoi(album.Tracks[j].TrackNumber)

			if errI != nil || errJ != nil {
				// If not numbers, sort alphabetically
				return album.Tracks[i].TrackNumber < album.Tracks[j].TrackNumber
			}
			return numI < numJ
		})
		albums = append(albums, *album)
	}

	// Sort albums by name
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].AlbumName < albums[j].AlbumName
	})

	m.previewAlbums = albums
	m.previewCursor = 0
}

// getTagValue safely extracts a tag value with a default
func getTagValue(tags map[string][]string, key string, defaultValue string) string {
	if values, ok := tags[key]; ok && len(values) > 0 && values[0] != "" {
		return values[0]
	}
	return defaultValue
}

// handlePreviewKey handles keyboard input in preview mode
func (m *model) handlePreviewKey(key string) {
	switch key {
	case "esc", "p":
		m.showPreview = false
		m.previewAlbums = nil
		m.previewCursor = 0

	case "up":
		m.previewCursor--
		if m.previewCursor < 0 {
			m.previewCursor = 0
		}

	case "down":
		if len(m.previewAlbums) > 0 {
			m.previewCursor++
			if m.previewCursor >= len(m.previewAlbums) {
				m.previewCursor = len(m.previewAlbums) - 1
			}
		}
	}
}
