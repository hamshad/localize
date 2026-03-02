package main

import (
	"math"
	"strings"
	"time"
)

// RenderASCIIClock returns a multi-line ASCII art analog clock face.
// Size: ~13 rows x 27 columns.
func RenderASCIIClock(t time.Time) string {
	// Clock face template (13x27)
	//        .---.
	//      /  12  \
	//    /11   |  1\
	//   |10  --+-- 2|
	//   | 9   /|  3 |
	//    \ 8 / | 4 /
	//      \  6  /
	//        '---'

	width := 27
	height := 13
	centerX := 13
	centerY := 6

	// Initialize grid
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Draw frame (box-drawing or simple ASCII)
	// Circumference points (approximate circle)
	frame := []struct {
		x, y int
		char string
	}{
		{11, 1, "."}, {12, 1, "-"}, {13, 1, "-"}, {14, 1, "-"}, {15, 1, "."},
		{9, 2, "/"}, {17, 2, "\\"},
		{7, 3, "/"}, {19, 3, "\\"},
		{6, 4, "|"}, {20, 4, "|"},
		{6, 5, "|"}, {20, 5, "|"},
		{6, 6, "|"}, {20, 6, "|"},
		{6, 7, "|"}, {20, 7, "|"},
		{6, 8, "|"}, {20, 8, "|"},
		{7, 9, "\\"}, {19, 9, "/"},
		{9, 10, "\\"}, {17, 10, "/"},
		{11, 11, "'"}, {12, 11, "-"}, {13, 11, "-"}, {14, 11, "-"}, {15, 11, "'"},
	}

	for _, p := range frame {
		grid[p.y][p.x] = "[white]" + p.char + "[-]"
	}

	// Draw numbers
	numbers := []struct {
		x, y int
		val  string
	}{
		{12, 2, "12"}, {17, 3, "1"}, {19, 5, "2"}, {20, 6, "3"}, {19, 7, "4"}, {17, 9, "5"},
		{13, 10, "6"}, {9, 9, "7"}, {7, 7, "8"}, {6, 6, "9"}, {7, 5, "10"}, {9, 3, "11"},
	}
	for _, n := range numbers {
		grid[n.y][n.x] = "[white]" + n.val + "[-]"
		if len(n.val) > 1 {
			grid[n.y][n.x+1] = "" // prevent overwriting next cell for 10, 11, 12
		}
	}

	// Center point
	grid[centerY][centerX] = "[white]+[-]"

	// Calculate hand positions
	hour := float64(t.Hour()%12) + float64(t.Minute())/60.0
	minute := float64(t.Minute())

	// Angles in radians (0 is up, clockwise)
	hourAngle := (hour / 12.0) * 2.0 * math.Pi
	minuteAngle := (minute / 60.0) * 2.0 * math.Pi

	// Hour hand (length 2-3)
	drawHand(grid, centerX, centerY, hourAngle, 2, "[yellow::b]")

	// Minute hand (length 4-5)
	drawHand(grid, centerX, centerY, minuteAngle, 4, "[green]")

	// Build final string
	var sb strings.Builder
	for _, row := range grid {
		for _, cell := range row {
			if cell != "" {
				sb.WriteString(cell)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func drawHand(grid [][]string, cx, cy int, angle float64, length int, color string) {
	// Sin/Cos for direction (0 is up)
	dx := math.Sin(angle)
	dy := -math.Cos(angle)

	for i := 1; i <= length; i++ {
		// Adjust for character aspect ratio (width ~ 0.5 * height)
		px := cx + int(math.Round(dx*float64(i)*2))
		py := cy + int(math.Round(dy*float64(i)))

		if px >= 0 && px < len(grid[0]) && py >= 0 && py < len(grid) {
			char := getHandChar(angle)
			grid[py][px] = color + char + "[-]"
		}
	}
}

func getHandChar(angle float64) string {
	// Normalize angle to [0, 2pi)
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Determine best character for direction
	deg := angle * 180 / math.Pi
	switch {
	case deg >= 337.5 || deg < 22.5:
		return "|"
	case deg >= 22.5 && deg < 67.5:
		return "/"
	case deg >= 67.5 && deg < 112.5:
		return "-"
	case deg >= 112.5 && deg < 157.5:
		return "\\"
	case deg >= 157.5 && deg < 202.5:
		return "|"
	case deg >= 202.5 && deg < 247.5:
		return "/"
	case deg >= 247.5 && deg < 292.5:
		return "-"
	case deg >= 292.5 && deg < 337.5:
		return "\\"
	}
	return "."
}
