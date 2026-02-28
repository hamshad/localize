package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// converterMode handles time conversion between timezones.
type converterMode struct {
	app          *tview.Application
	inputTime    string   // HH:MM format input
	selectedZone int      // Index of selected source timezone
	zones        []Region // All available zones for conversion
	inputMode    bool     // True if waiting for time input
	timeError    string   // Error message for invalid input
}

// newConverterMode creates a new converter mode handler.
func newConverterMode(app *tview.Application) *converterMode {
	// Collect all available zones
	var allZones []Region
	allZones = append(allZones, leftRegions...)
	allZones = append(allZones, rightRegions...)

	return &converterMode{
		app:          app,
		inputTime:    "",
		selectedZone: 0,
		zones:        allZones,
		inputMode:    false,
		timeError:    "",
	}
}

// GetMode returns the mode type.
func (c *converterMode) GetMode() Mode {
	return ModeConverter
}

// HandleKey handles key events in converter mode.
func (c *converterMode) HandleKey(key rune) bool {
	switch key {
	case 27: // Escape
		return false // Let modeManager handle escape
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if c.inputMode {
			if len(c.inputTime) < 5 {
				c.inputTime += string(key)
				// Auto-add colon after 2 digits
				if len(c.inputTime) == 2 {
					c.inputTime += ":"
				}
			}
		}
		c.timeError = ""
		return true
	case ':':
		if c.inputMode && !strings.Contains(c.inputTime, ":") {
			c.inputTime += ":"
		}
		return true
	case 'c', 'C':
		// Toggle input mode
		c.inputMode = !c.inputMode
		if !c.inputMode {
			c.inputTime = ""
			c.timeError = ""
		}
		return true
	case 16: // Ctrl+P or Up arrow equivalent - previous zone
		if c.selectedZone > 0 {
			c.selectedZone--
		}
		return true
	case 14: // Ctrl+N or Down arrow equivalent - next zone
		if c.selectedZone < len(c.zones)-1 {
			c.selectedZone++
		}
		return true
	case 'j', 'J': // Down
		if c.selectedZone < len(c.zones)-1 {
			c.selectedZone++
		}
		return true
	case 'k', 'K': // Up
		if c.selectedZone > 0 {
			c.selectedZone--
		}
		return true
	case 'r', 'R':
		// Reset input
		c.inputTime = ""
		c.timeError = ""
		c.inputMode = false
		return true
	}
	return false
}

// Render returns the rendered converter display.
func (c *converterMode) Render() string {
	var b strings.Builder

	// Header
	b.WriteString("\n[yellow::b]━━━ TIME CONVERTER ━━━[-::-]\n\n")

	// Source time display
	sourceTime := c.getSourceTime()
	if sourceTime != nil {
		zone := c.zones[c.selectedZone]
		b.WriteString(fmt.Sprintf("  [dodgerblue::b]Source:[-::-] %s @ %s\n",
			sourceTime.Format("15:04:05"), zone.Name))
		b.WriteString(fmt.Sprintf("         [%s] %s\n\n",
			sourceTime.Format("Mon, 02 Jan 2006"), zone.Timezone))
	} else {
		b.WriteString("  [dodgerblue::b]Enter source time below[white]\n\n")
	}

	// Time input
	if c.inputMode {
		displayTime := c.inputTime
		if displayTime == "" {
			displayTime = "HH:MM"
		}
		b.WriteString(fmt.Sprintf("  [::b]Enter time: [yellow]%s[-]\n\n", displayTime))
		if c.timeError != "" {
			b.WriteString(fmt.Sprintf("  [red]%s[-]\n\n", c.timeError))
		}
	} else {
		b.WriteString("  [darkgray]Press C to enter time[white]\n\n")
	}

	// Converted times
	b.WriteString("  [aqua::b]Converted Times:[-::-]\n")
	b.WriteString("  [darkgray]Use ↑/↓ to select source timezone[white]\n\n")

	if sourceTime != nil {
		for i, zone := range c.zones {
			if i == c.selectedZone {
				continue // Skip source zone
			}
			converted := c.convertTime(sourceTime, zone.Timezone)
			marker := "  "
			if i == c.selectedZone+1 || (c.selectedZone == len(c.zones)-1 && i == 0) {
				marker = "> "
			}
			colorTag := colorToTag(zone.Color)
			b.WriteString(fmt.Sprintf("%s[%s::b]%-13s[-::-] %s  [%s]\n",
				marker, colorTag, zone.Name, converted.Format("15:04:05"),
				converted.Format("Mon, 02 Jan")))
		}
	}

	return b.String()
}

// getSourceTime parses the input time and returns it in the selected zone.
func (c *converterMode) getSourceTime() *time.Time {
	if len(c.inputTime) < 4 {
		return nil
	}

	// Parse time
	sourceZone := c.zones[c.selectedZone].Timezone

	// Create a time in the source timezone
	inputStr := c.inputTime
	if len(inputStr) == 4 {
		inputStr = "0" + inputStr // Pad single digit hour
	}

	loc, err := time.LoadLocation(sourceZone)
	if err != nil {
		c.timeError = "Invalid timezone"
		return nil
	}

	// Parse the input time
	parsed, err := time.ParseInLocation("15:04", inputStr, loc)
	if err != nil {
		c.timeError = "Invalid format (use HH:MM)"
		return nil
	}

	// Return time with today's date in the source zone
	now := time.Now().In(loc)
	result := time.Date(now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, loc)
	return &result
}

// convertTime converts a time to a target timezone.
func (c *converterMode) convertTime(t *time.Time, targetZone string) time.Time {
	loc, err := time.LoadLocation(targetZone)
	if err != nil {
		return *t
	}
	return t.In(loc)
}

// GetHelpText returns the help text for converter mode.
func (c *converterMode) GetHelpText() string {
	return "[darkgray]Keys:[white] C=Enter Time  ↑/↓=Select Zone  R=Reset  Esc=Exit"
}

// HandleSpecialKeyEvent handles non-rune key events (Enter, Backspace, etc.).
func (c *converterMode) HandleSpecialKeyEvent(key tcell.Key) bool {
	switch key {
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if c.inputMode && len(c.inputTime) > 0 {
			c.inputTime = c.inputTime[:len(c.inputTime)-1]
			// Remove trailing colon if backspacing exposed one
			if len(c.inputTime) > 0 && c.inputTime[len(c.inputTime)-1] == ':' {
				c.inputTime = c.inputTime[:len(c.inputTime)-1]
			}
			return true
		}
	}
	return false
}
