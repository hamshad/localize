package main

import (
	"fmt"
	"strings"
	"time"
)

// dayNightState tracks whether day/night overlay is enabled
var dayNightOverlayEnabled = false

// ToggleDayNightOverlay toggles the day/night overlay
func ToggleDayNightOverlay() {
	dayNightOverlayEnabled = !dayNightOverlayEnabled
}

// IsDayNightOverlayEnabled returns whether the overlay is enabled
func IsDayNightOverlayEnabled() bool {
	return dayNightOverlayEnabled
}

// getDayPhaseForLocation determines the day phase for a specific longitude.
// Returns: "day", "night", "dawn", "dusk"
func getDayPhaseForLocation(t time.Time, longitude float64) string {
	// Calculate local hour at this longitude
	// Each 15 degrees = 1 hour time difference
	utcHour := t.UTC().Hour()
	localHour := (utcHour + int(longitude/15)) % 24
	if localHour < 0 {
		localHour += 24
	}

	switch {
	case localHour >= 6 && localHour < 7:
		return "dawn"
	case localHour >= 7 && localHour < 8:
		return "morning"
	case localHour >= 8 && localHour < 17:
		return "day"
	case localHour >= 17 && localHour < 18:
		return "dusk"
	case localHour >= 18 && localHour < 20:
		return "evening"
	default:
		return "night"
	}
}

// getColorForDayPhase returns the appropriate color tag for a day phase
func getColorForDayPhase(phase string) string {
	switch phase {
	case "dawn":
		return "#FFA500" // Amber for dawn
	case "morning":
		return "yellow" // Yellow for morning
	case "day":
		return "green" // Full bright green for day
	case "dusk":
		return "#CD853F" // Peru/dusk color
	case "evening":
		return "#6A5ACD" // Slate blue for evening
	case "night":
		return "#2F4F4F" // Dark slate gray for night
	default:
		return "green"
	}
}

// colorizeBrailleMapWithDayNight adds day/night coloring to the braille map.
// Each longitude position represents a different local time.
func colorizeBrailleMapWithDayNight(brailleMap string) string {
	lines := strings.Split(strings.TrimRight(brailleMap, "\n"), "\n")

	// Map dimensions: 240 columns (each 2 chars = 1 braille), 96 rows
	// Longitude ranges from -180 to 180
	// We need to determine the local time for each column position

	now := time.Now().UTC()

	// First, overlay city markers onto map lines
	for _, m := range cityMarkers {
		if m.row < len(lines) {
			runes := []rune(lines[m.row])
			label := m.label
			if m.col+len(label) <= len(runes) {
				before := string(runes[:m.col])
				after := string(runes[m.col+len(label):])
				lines[m.row] = before + fmt.Sprintf("[%s::b]%s[-::-]", m.color, label) + after
			}
		}
	}

	// Now colorize based on day/night if enabled
	if !dayNightOverlayEnabled {
		// Original behavior: all green
		var sb strings.Builder
		for _, line := range lines {
			sb.WriteString("[green]")
			sb.WriteString(line)
			sb.WriteString("[-]\n")
		}
		return sb.String()
	}

	// Apply day/night coloring
	// Each braille column represents about 3 degrees of longitude (180/60 = 3)
	brailleCols := 120 // 240 / 2
	degreesPerCol := 360.0 / float64(brailleCols)

	for row := 0; row < len(lines); row++ {
		line := lines[row]
		runes := []rune(line)

		var newLine strings.Builder
		col := 0
		for col < len(runes) {
			// Calculate longitude for this braille column position
			longitude := -180 + float64(col*2)*degreesPerCol
			phase := getDayPhaseForLocation(now, longitude)
			color := getColorForDayPhase(phase)

			// Check if this is a color tag we inserted
			if col < len(runes)-3 && string(runes[col:col+3]) == "[-]" {
				// We're at end of a color tag, continue
				newLine.WriteRune(runes[col])
				col++
				continue
			}

			// Check if we're inside a color tag we added (like city markers)
			inColorTag := false
			tagEnd := -1
			for tagLen := 5; tagLen <= 20 && col+tagLen <= len(runes); tagLen++ {
				testStr := string(runes[col : col+tagLen])
				if strings.HasPrefix(testStr, "[") && strings.Contains(testStr, "]") {
					inColorTag = true
					tagEnd = col + strings.Index(testStr, "[-]")
					if tagEnd >= col {
						tagEnd += 3 // include the [-] part
						break
					}
				}
			}

			if inColorTag && tagEnd > col {
				newLine.WriteString(string(runes[col:tagEnd]))
				col = tagEnd
				continue
			}

			// Regular character - apply day/night color
			// Check if it's a braille character (unicode range)
			if col < len(runes) {
				ch := runes[col]
				if ch >= 0x2800 && ch <= 0x28FF {
					newLine.WriteString(color)
					newLine.WriteRune(ch)
					newLine.WriteString("[-]")
					col++
					continue
				}
			}

			// Regular character
			newLine.WriteRune(runes[col])
			col++
		}
		lines[row] = newLine.String()
	}

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteRune('\n')
	}

	return sb.String()
}

// GetDayNightStatus returns the current day/night overlay status text for the header
func GetDayNightStatus() string {
	if dayNightOverlayEnabled {
		return "[yellow]Day/Night ON[white]"
	}
	return "[darkgray]Day/Night OFF[white]"
}
