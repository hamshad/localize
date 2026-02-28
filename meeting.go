package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MeetingPlanner manages the meeting time planning functionality.
type MeetingPlanner struct {
	app            *tview.Application
	selectedCities []City
	businessStart  int // Hour (0-23) when business hours start
	businessEnd    int // Hour (0-23) when business hours end
	selectedIndex  int // Current selection index in city list
	mode           int // 0 = selecting cities, 1 = viewing timeline
	timelineStart  int // Starting hour for timeline view
}

// NewMeetingPlanner creates a new MeetingPlanner instance.
func NewMeetingPlanner(app *tview.Application) *MeetingPlanner {
	return &MeetingPlanner{
		app:            app,
		selectedCities: []City{},
		businessStart:  9,  // 9 AM
		businessEnd:    17, // 5 PM
		selectedIndex:  0,
		mode:           0,
		timelineStart:  0,
	}
}

// AddCity adds a city to the meeting planner selection.
func (mp *MeetingPlanner) AddCity(city City) {
	// Check if already selected
	for _, c := range mp.selectedCities {
		if c.Name == city.Name {
			return
		}
	}
	mp.selectedCities = append(mp.selectedCities, city)
}

// RemoveCity removes a city from the selection.
func (mp *MeetingPlanner) RemoveCity(city City) {
	for i, c := range mp.selectedCities {
		if c.Name == city.Name {
			mp.selectedCities = append(mp.selectedCities[:i], mp.selectedCities[i+1:]...)
			return
		}
	}
}

// ToggleCity toggles a city in the selection.
func (mp *MeetingPlanner) ToggleCity(city City) {
	for i, c := range mp.selectedCities {
		if c.Name == city.Name {
			mp.selectedCities = append(mp.selectedCities[:i], mp.selectedCities[i+1:]...)
			return
		}
	}
	mp.selectedCities = append(mp.selectedCities, city)
}

// IsSelected checks if a city is selected.
func (mp *MeetingPlanner) IsSelected(city City) bool {
	for _, c := range mp.selectedCities {
		if c.Name == city.Name {
			return true
		}
	}
	return false
}

// ClearSelection clears all selected cities.
func (mp *MeetingPlanner) ClearSelection() {
	mp.selectedCities = []City{}
	mp.selectedIndex = 0
	mp.mode = 0
}

// SetBusinessHours sets the business hours for the meeting planner.
func (mp *MeetingPlanner) SetBusinessHours(start, end int) {
	mp.businessStart = start
	mp.businessEnd = end
}

// GetBestMeetingTimes calculates the best meeting times across all selected cities.
// Hours are in UTC. For each UTC hour, we check what local time it would be in each city.
func (mp *MeetingPlanner) GetBestMeetingTimes() []struct {
	Hour         int
	AllAvailable bool
	Count        int
} {
	var results []struct {
		Hour         int
		AllAvailable bool
		Count        int
	}

	now := time.Now()

	for utcHour := 0; utcHour < 24; utcHour++ {
		count := 0
		for _, city := range mp.selectedCities {
			loc, err := time.LoadLocation(city.Timezone)
			if err != nil {
				continue
			}
			// Create a time at the specified UTC hour, then convert to city's timezone
			utcTime := time.Date(now.Year(), now.Month(), now.Day(), utcHour, 0, 0, 0, time.UTC)
			cityTime := utcTime.In(loc)
			hourInCity := cityTime.Hour()
			if hourInCity >= mp.businessStart && hourInCity < mp.businessEnd {
				count++
			}
		}
		results = append(results, struct {
			Hour         int
			AllAvailable bool
			Count        int
		}{
			Hour:         utcHour,
			AllAvailable: count == len(mp.selectedCities) && len(mp.selectedCities) > 0,
			Count:        count,
		})
	}
	return results
}

// RenderCitySelection renders the city selection view.
func (mp *MeetingPlanner) RenderCitySelection() string {
	var b strings.Builder
	b.WriteString("[yellow::b]Meeting Time Planner[::-]\n\n")
	b.WriteString("[::b]Select cities (space to toggle, Enter to view timeline):[::-]\n\n")

	// Group cities by category
	categories := []string{"Americas", "Europe", "MiddleEast", "Africa", "Asia", "Oceania"}

	for _, category := range categories {
		cities := GetCitiesByCategory(category)
		if len(cities) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf("[%s::b]%s:[::-]\n", getCategoryColor(category), category))

		for _, city := range cities {
			selected := mp.IsSelected(city)
			marker := "[ ] "
			if selected {
				marker = "[green]✓[white] "
			}

			// Highlight the current selection
			prefix := "  "
			if city.Name == AllCities[mp.selectedIndex].Name {
				prefix = "[yellow]►[white] "
			}

			b.WriteString(fmt.Sprintf("%s%s%s (%s)\n", prefix, marker, city.Name, city.Timezone))
		}
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("\n[::b]Selected: %d cities[::-]\n", len(mp.selectedCities)))
	b.WriteString("[silver]Press Enter to view timeline | Space to toggle | C to clear[::-]\n")

	return b.String()
}

// RenderTimeline renders the meeting timeline view.
func (mp *MeetingPlanner) RenderTimeline() string {
	var b strings.Builder

	b.WriteString("[yellow::b]Meeting Timeline[::-]\n\n")

	// Show selected cities
	b.WriteString("[::b]Cities:[::-] ")
	for i, city := range mp.selectedCities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("[%s]%s[white]", colorToTag(city.Color), city.Name))
	}
	b.WriteString(fmt.Sprintf("\n[::b]Business hours:[::-] %d:00 - %d:00\n\n", mp.businessStart, mp.businessEnd))

	// Calculate and display best times
	times := mp.GetBestMeetingTimes()

	// Find best slots (all available)
	var bestHours []int
	var goodHours []int

	for _, t := range times {
		if t.AllAvailable {
			bestHours = append(bestHours, t.Hour)
		} else if t.Count >= len(mp.selectedCities)/2 && len(mp.selectedCities) > 0 {
			goodHours = append(goodHours, t.Hour)
		}
	}

	// Display best meeting times
	if len(bestHours) > 0 {
		b.WriteString("[green::b]✓ BEST Times (all cities in business hours):[::-]\n")
		for i, h := range bestHours {
			if i > 0 && i%6 == 0 {
				b.WriteString("\n")
			}
			b.WriteString(fmt.Sprintf(" %02d:00 ", h))
		}
		b.WriteString("\n\n")
	}

	if len(goodHours) > 0 {
		b.WriteString("[yellow::b]⚠ Acceptable Times (most cities available):[::-]\n")
		for i, h := range goodHours {
			if i > 0 && i%6 == 0 {
				b.WriteString("\n")
			}
			b.WriteString(fmt.Sprintf(" %02d:00 ", h))
		}
		b.WriteString("\n\n")
	}

	if len(bestHours) == 0 && len(goodHours) == 0 {
		b.WriteString("[red]No overlapping business hours found.[-]\n\n")
	}

	// Show detailed timeline
	b.WriteString("[::b]24-Hour Timeline (Green=All, Yellow=Some, Red=None):[::-]\n")
	b.WriteString("     ")

	// Header row with hours
	for h := 0; h < 24; h++ {
		b.WriteString(fmt.Sprintf("%02d", h))
	}
	b.WriteString("\n")

	// For each city, show their availability
	for _, city := range mp.selectedCities {
		loc, err := time.LoadLocation(city.Timezone)
		if err != nil {
			continue
		}

		now := time.Now()

		colorTag := colorToTag(city.Color)
		displayName := city.Name
		if len(displayName) > 12 {
			displayName = displayName[:12]
		}
		b.WriteString(fmt.Sprintf("[%s]%-12s[white] ", colorTag, displayName))

		for h := 0; h < 24; h++ {
			// Calculate what local hour it would be in this city when UTC is h:00
			utcTime := time.Date(now.Year(), now.Month(), now.Day(), h, 0, 0, 0, time.UTC)
			cityTime := utcTime.In(loc)
			displayHour := cityTime.Hour()

			if displayHour >= mp.businessStart && displayHour < mp.businessEnd {
				b.WriteString("[green]▀[white]")
			} else if displayHour >= mp.businessStart-2 && displayHour < mp.businessEnd+2 {
				b.WriteString("[yellow]▀[white]")
			} else {
				b.WriteString("[red]▀[white]")
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("\n[yellow::b]Your Reference Time (UTC):[::-] ")
	utcNow := time.Now().UTC()
	b.WriteString(utcNow.Format("15:04"))
	b.WriteString("\n\n")

	b.WriteString("[silver]Press Enter to edit selection | B to adjust business hours | Escape to exit[::-]\n")

	return b.String()
}

// Render returns the appropriate view based on mode.
func (mp *MeetingPlanner) Render() string {
	if mp.mode == 0 || len(mp.selectedCities) == 0 {
		return mp.RenderCitySelection()
	}
	return mp.RenderTimeline()
}

// HandleKey handles key input for the meeting planner.
// Takes the rune character for character input.
func (mp *MeetingPlanner) HandleKey(ch rune) bool {
	switch ch {
	case ' ':
		// Toggle current city
		if mp.selectedIndex < len(AllCities) {
			mp.ToggleCity(AllCities[mp.selectedIndex])
		}
		return true
	case 'c', 'C':
		mp.ClearSelection()
		return true
	case 'b', 'B':
		// Cycle through business hour presets
		if mp.businessStart == 9 {
			mp.businessStart = 8
			mp.businessEnd = 16
		} else if mp.businessStart == 8 {
			mp.businessStart = 10
			mp.businessEnd = 18
		} else {
			mp.businessStart = 9
			mp.businessEnd = 17
		}
		return true
	}
	return false
}

// HandleSpecialKey handles special keys (not character input).
func (mp *MeetingPlanner) HandleSpecialKey(key tcell.Key) bool {
	switch key {
	case tcell.KeyEscape:
		mp.mode = 0
		mp.selectedCities = []City{}
		return true
	case tcell.KeyEnter:
		if mp.mode == 0 && len(mp.selectedCities) > 0 {
			mp.mode = 1
		} else if mp.mode == 1 {
			mp.mode = 0
		}
		return true
	case tcell.KeyUp:
		if mp.selectedIndex > 0 {
			mp.selectedIndex--
		}
		return true
	case tcell.KeyDown:
		if mp.selectedIndex < len(AllCities)-1 {
			mp.selectedIndex++
		}
		return true
	}
	return false
}

// HasSelectedCities returns whether any cities are selected.
func (mp *MeetingPlanner) HasSelectedCities() bool {
	return len(mp.selectedCities) > 0
}

// GetSelectedCities returns the list of selected cities.
func (mp *MeetingPlanner) GetSelectedCities() []City {
	return mp.selectedCities
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getCategoryColor returns the color for a category.
func getCategoryColor(category string) string {
	switch category {
	case "Americas":
		return "dodgerblue"
	case "Europe":
		return "green"
	case "MiddleEast":
		return "gold"
	case "Africa":
		return "sandybrown"
	case "Asia":
		return "orange"
	case "Oceania":
		return "yellow"
	default:
		return "white"
	}
}

// MeetingMode wraps MeetingPlanner to implement ModeHandler interface.
type MeetingMode struct {
	planner *MeetingPlanner
}

// NewMeetingMode creates a new MeetingMode instance.
func NewMeetingMode(app *tview.Application) *MeetingMode {
	return &MeetingMode{
		planner: NewMeetingPlanner(app),
	}
}

// GetMode returns the mode type.
func (m *MeetingMode) GetMode() Mode {
	return ModeMeeting
}

// HandleKey handles key input for the meeting mode.
func (m *MeetingMode) HandleKey(key rune) bool {
	return m.planner.HandleKey(key)
}

// Render returns the rendered content for the meeting mode.
func (m *MeetingMode) Render() string {
	return m.planner.Render()
}

// GetHelpText returns the help text for the meeting mode.
func (m *MeetingMode) GetHelpText() string {
	if m.planner.mode == 0 {
		return "[darkgray]Keys:[white] ↑/↓=Navigate  Space=Toggle  Enter=View Timeline  C=Clear  Esc=Exit"
	}
	return "[darkgray]Keys:[white] Enter=Back to Selection  B=Change Hours  C=Clear  Esc=Exit"
}

// HandleSpecialKeyEvent handles non-rune key events (Enter, Backspace, etc.).
func (m *MeetingMode) HandleSpecialKeyEvent(key tcell.Key) bool {
	// Meeting mode handles special keys via HandleSpecialKey in the main input handler
	return false
}
