package main

import (
	"strings"
)

// ScaleBitmap scales a string bitmap to the target dimensions using nearest-neighbor sampling.
// Input: src is the source bitmap (each string is a row of '0'/'1' characters).
// Output: 2D bool grid where true = land.
func ScaleBitmap(src []string, srcW, srcH, dstW, dstH int) [][]bool {
	result := make([][]bool, dstH)
	for i := range result {
		result[i] = make([]bool, dstW)
	}

	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			// Map destination pixel to source pixel
			srcX := (x * srcW) / dstW
			srcY := (y * srcH) / dstH

			// Clamp to source bounds
			if srcX >= srcW {
				srcX = srcW - 1
			}
			if srcY >= srcH {
				srcY = srcH - 1
			}

			// Get pixel value
			if srcY < len(src) && srcX < len(src[srcY]) {
				result[y][x] = src[srcY][srcX] == '1'
			}
		}
	}

	return result
}

// GetBrailleGridSize returns the braille grid dimensions for a given terminal size.
// Subtracts 2 for border (top/bottom) and 1 for status bar.
func GetBrailleGridSize(termWidth, termHeight int) (cols, rows int) {
	cols = termWidth - 2  // Border left/right
	rows = termHeight - 3 // Border top/bottom + status bar
	return cols, rows
}

// LatLonToBraille converts geographic coordinates to braille grid position.
// Uses equirectangular projection.
func LatLonToBraille(lat, lon float64, brailleCols, brailleRows int) (col, row int) {
	// Equirectangular projection
	col = int((lon + 180.0) / 360.0 * float64(brailleCols))
	row = int((90.0 - lat) / 180.0 * float64(brailleRows))

	// Clamp to grid bounds
	if col < 0 {
		col = 0
	}
	if col >= brailleCols {
		col = brailleCols - 1
	}
	if row < 0 {
		row = 0
	}
	if row >= brailleRows {
		row = brailleRows - 1
	}

	return col, row
}

// RenderBrailleMap renders the world map to a braille string for the given terminal dimensions.
// Returns a multi-line braille string (no color tags - colorization is separate).
func RenderBrailleMap(termWidth, termHeight int) string {
	// Calculate available braille grid
	brailleCols, brailleRows := GetBrailleGridSize(termWidth, termHeight)

	// Calculate pixel dimensions (each braille char is 2x4 pixels)
	pixelW := brailleCols * 2
	pixelH := brailleRows * 4

	// Scale bitmap to target size
	scaled := ScaleBitmap(WorldBitmap, WorldBitmapWidth, WorldBitmapHeight, pixelW, pixelH)

	// Encode to braille characters
	// Braille encoding: each braille char is 2 cols x 4 rows of dots
	// Dot positions and their bit values:
	// (0,0)=0x01  (1,0)=0x08
	// (0,1)=0x02  (1,1)=0x10
	// (0,2)=0x04  (1,2)=0x20
	// (0,3)=0x40  (1,3)=0x80
	brailleBase := rune(0x2800)

	var sb strings.Builder
	for by := 0; by < brailleRows; by++ {
		for bx := 0; bx < brailleCols; bx++ {
			var code rune
			// Map the 2x4 block to braille dots
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					py := by*4 + dy
					px := bx*2 + dx
					if py < pixelH && px < pixelW && scaled[py][px] {
						// Braille dot encoding
						switch {
						case dx == 0 && dy == 0:
							code |= 0x01
						case dx == 0 && dy == 1:
							code |= 0x02
						case dx == 0 && dy == 2:
							code |= 0x04
						case dx == 1 && dy == 0:
							code |= 0x08
						case dx == 1 && dy == 1:
							code |= 0x10
						case dx == 1 && dy == 2:
							code |= 0x20
						case dx == 0 && dy == 3:
							code |= 0x40
						case dx == 1 && dy == 3:
							code |= 0x80
						}
					}
				}
			}
			sb.WriteRune(brailleBase + code)
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}
