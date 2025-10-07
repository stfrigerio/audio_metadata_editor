package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/disintegration/imaging"
)

// renderImageAsBlocks converts an image to colored block characters for terminal display
func renderImageAsBlocks(imageData []byte, width, height int) []string {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return []string{"\033[90mError loading image\033[0m"}
	}

	// Resize to fit terminal (half height because terminal chars are taller than wide)
	resized := imaging.Fit(img, width, height*2, imaging.Lanczos)

	lines := make([]string, 0)
	bounds := resized.Bounds()

	// Process two rows at a time (using half-block characters)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		var line strings.Builder

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get top and bottom pixel colors
			topColor := resized.At(x, y)
			tr, tg, tb, _ := topColor.RGBA()

			// Check if there's a bottom pixel
			var br, bg, bb uint32
			if y+1 < bounds.Max.Y {
				bottomColor := resized.At(x, y+1)
				br, bg, bb, _ = bottomColor.RGBA()
			} else {
				br, bg, bb = tr, tg, tb
			}

			// Convert to 8-bit color
			tr8, tg8, tb8 := uint8(tr>>8), uint8(tg>>8), uint8(tb>>8)
			br8, bg8, bb8 := uint8(br>>8), uint8(bg>>8), uint8(bb>>8)

			// Use upper half block (▀) with foreground color for top, background for bottom
			line.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀\033[0m",
				tr8, tg8, tb8, br8, bg8, bb8))
		}

		lines = append(lines, line.String())
	}

	return lines
}

// renderCoverArtPreview creates a small preview of the cover art
func renderCoverArtPreview(imageData []byte) []string {
	// Use much higher resolution for better quality
	return renderImageAsBlocks(imageData, 60, 30)
}
