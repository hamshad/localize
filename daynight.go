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
// Returns: "day", "night", "dawn", "dusk", "morning", "evening"
func getDayPhaseForLocation(t time.Time, longitude float64) string {
	// Calculate local solar hour at this longitude
	// Each 15 degrees = 1 hour time difference from UTC
	utcHour := float64(t.UTC().Hour()) + float64(t.UTC().Minute())/60.0
	localHour := utcHour + longitude/15.0
	// Normalize to 0-24 range
	for localHour < 0 {
		localHour += 24
	}
	for localHour >= 24 {
		localHour -= 24
	}
	intHour := int(localHour)

	switch {
	case intHour >= 6 && intHour < 7:
		return "dawn"
	case intHour >= 7 && intHour < 8:
		return "morning"
	case intHour >= 8 && intHour < 17:
		return "day"
	case intHour >= 17 && intHour < 18:
		return "dusk"
	case intHour >= 18 && intHour < 20:
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
	// Each braille column represents about 3 degrees of longitude (360/120 = 3)
	brailleCols := 120 // 240 / 2
	degreesPerCol := 360.0 / float64(brailleCols)

	for row := 0; row < len(lines); row++ {
		line := lines[row]
		runes := []rune(line)

		var newLine strings.Builder
		brailleCol := 0 // Track the braille column for longitude calc
		col := 0
		for col < len(runes) {
			ch := runes[col]

			// If we hit a '[', this is a tview color/style tag - pass it through verbatim
			if ch == '[' {
				// Find the matching ']'
				end := col + 1
				for end < len(runes) && runes[end] != ']' {
					end++
				}
				if end < len(runes) {
					end++ // include the ']'
					newLine.WriteString(string(runes[col:end]))
					col = end
					continue
				}
			}

			// Braille character - apply day/night color
			if ch >= 0x2800 && ch <= 0x28FF {
				longitude := -180.0 + float64(brailleCol)*degreesPerCol
				phase := getDayPhaseForLocation(now, longitude)
				color := getColorForDayPhase(phase)
				newLine.WriteString("[")
				newLine.WriteString(color)
				newLine.WriteString("]")
				newLine.WriteRune(ch)
				newLine.WriteString("[-]")
				brailleCol++
				col++
				continue
			}

			// Any other character (like city label letters) - pass through
			// These replace braille chars, so they occupy a braille column position
			newLine.WriteRune(ch)
			brailleCol++
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
