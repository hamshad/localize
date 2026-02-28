package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NavigationState tracks the current navigation selection
type NavigationState struct {
	selectedIndex    int            // -1 means none selected
	selectedPanel    string         // "left" or "right"
	detailsVisible   bool           // whether details panel is showing
	selectedCity     *Region        // currently selected city details
	selectedTime     time.Time      // time when selected
	selectedTimezone *time.Location // timezone of selected city
}

// global navigation state
var navState = &NavigationState{
	selectedIndex:  -1,
	selectedPanel:  "left",
	detailsVisible: false,
}

// NavigationView is the details panel for selected city
var NavigationView *tview.TextView

// InitNavigation creates and initializes the navigation view
func InitNavigation() *tview.TextView {
	NavigationView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false)
	NavigationView.SetBorder(true).
		SetTitle(" [ City Details ] ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorYellow)
	return NavigationView
}

// GetNavigationView returns the navigation details view
func GetNavigationView() *tview.TextView {
	return NavigationView
}

// IsNavigationActive returns whether any city is selected
func IsNavigationActive() bool {
	return navState.selectedIndex >= 0
}

// IsDetailsVisible returns whether the details panel is visible
func IsDetailsVisible() bool {
	return navState.detailsVisible
}

// SelectCity selects a city at the given index in the specified panel
func SelectCity(panel string, index int, regions []Region) {
	navState.selectedPanel = panel
	navState.selectedIndex = index
	navState.detailsVisible = false
	navState.selectedCity = nil

	if index >= 0 && index < len(regions) {
		navState.selectedCity = &regions[index]
		loc, err := time.LoadLocation(regions[index].Timezone)
		if err == nil {
			navState.selectedTime = time.Now().In(loc)
			navState.selectedTimezone = loc
		}
	}

	updateNavigationView()
}

// NavigateUp moves selection up in the current panel
func NavigateUp(leftRegions []Region, rightRegions []Region) bool {
	regions := getCurrentPanelRegions(leftRegions, rightRegions)
	if len(regions) == 0 {
		return false
	}

	newIndex := navState.selectedIndex - 1
	if newIndex < 0 {
		newIndex = len(regions) - 1 // Wrap to bottom
	}

	SelectCity(navState.selectedPanel, newIndex, regions)
	return true
}

// NavigateDown moves selection down in the current panel
func NavigateDown(leftRegions []Region, rightRegions []Region) bool {
	regions := getCurrentPanelRegions(leftRegions, rightRegions)
	if len(regions) == 0 {
		return false
	}

	newIndex := navState.selectedIndex + 1
	if newIndex >= len(regions) {
		newIndex = 0 // Wrap to top
	}

	SelectCity(navState.selectedPanel, newIndex, regions)
	return true
}

// SwitchPanel switches to the other clock panel
func SwitchPanel(leftRegions []Region, rightRegions []Region) {
	if navState.selectedPanel == "left" {
		navState.selectedPanel = "right"
		if navState.selectedIndex >= len(rightRegions) {
			navState.selectedIndex = len(rightRegions) - 1
		}
		if navState.selectedIndex >= 0 && navState.selectedIndex < len(rightRegions) {
			navState.selectedCity = &rightRegions[navState.selectedIndex]
		}
	} else {
		navState.selectedPanel = "left"
		if navState.selectedIndex >= len(leftRegions) {
			navState.selectedIndex = len(leftRegions) - 1
		}
		if navState.selectedIndex >= 0 && navState.selectedIndex < len(leftRegions) {
			navState.selectedCity = &leftRegions[navState.selectedIndex]
		}
	}

	updateNavigationView()
}

// ToggleDetails shows/hides the city details panel
func ToggleDetails(leftRegions []Region, rightRegions []Region) bool {
	regions := getCurrentPanelRegions(leftRegions, rightRegions)
	if navState.selectedIndex < 0 || navState.selectedIndex >= len(regions) {
		return false
	}

	navState.detailsVisible = !navState.detailsVisible
	updateNavigationView()
	return true
}

// Deselect clears the current selection
func Deselect() {
	navState.selectedIndex = -1
	navState.selectedPanel = "left"
	navState.detailsVisible = false
	navState.selectedCity = nil
	updateNavigationView()
}

// getCurrentPanelRegions returns the regions for the currently selected panel
func getCurrentPanelRegions(leftRegions []Region, rightRegions []Region) []Region {
	if navState.selectedPanel == "left" {
		return leftRegions
	}
	return rightRegions
}

// updateNavigationView updates the navigation details panel
func updateNavigationView() {
	if NavigationView == nil {
		return
	}

	if navState.selectedIndex < 0 {
		NavigationView.SetText("[darkgray]Use arrow keys to navigate cities[white]")
		NavigationView.SetBorderColor(tcell.ColorYellow)
		return
	}

	if navState.selectedCity == nil {
		return
	}

	r := *navState.selectedCity
	loc, _ := time.LoadLocation(r.Timezone)
	now := time.Now().In(loc)

	// Format timezone info
	_, offset := now.Zone()
	offsetHours := offset / 3600
	offsetMins := (offset % 3600) / 60
	offsetStr := fmt.Sprintf("%+03d:%02d", offsetHours, offsetMins)

	// Day of week
	dayOfWeek := now.Weekday().String()

	// Check DST (simplified - compare January and July offsets)
	isDST := false
	_, janOffset := time.Date(2024, 1, 15, 12, 0, 0, 0, loc).Zone()
	_, julOffset := time.Date(2024, 7, 15, 12, 0, 0, 0, loc).Zone()
	if janOffset != julOffset {
		// DST exists in this timezone
		_, currentOffset := now.Zone()
		// The larger offset is typically DST
		maxOffset := janOffset
		if julOffset > maxOffset {
			maxOffset = julOffset
		}
		isDST = currentOffset == maxOffset
	}

	// Date formatting
	dateStr := now.Format("Monday, 02 January 2006")
	timeStr := now.Format("15:04:05")

	panelArrow := "→"
	if navState.selectedPanel == "right" {
		panelArrow = "←"
	}

	var detailsText string
	if navState.detailsVisible {
		detailsText = fmt.Sprintf(`[yellow]%s[white] selected %s

[aqua]City:[white]      %s
[aqua]Timezone:[white]   %s
[aqua]Current Time:[white] %s
[aqua]Date:[white]       %s
[aqua]UTC Offset:[white]  %s
[aqua]Day of Week:[white] %s
[aqua]DST Active:[white]  %s

[darkgray]Press Enter to hide details[white]`,
			panelArrow,
			navState.selectedPanel+" panel",
			r.Name,
			r.Timezone,
			timeStr,
			dateStr,
			offsetStr,
			dayOfWeek,
			formatBool(isDST),
		)
	} else {
		detailsText = fmt.Sprintf(`[yellow]%s[white] %s selected

[aqua]City:[white]    %s
[aqua]Timezone:[white] %s
[aqua]UTC Offset:[white] %s

[darkgray]Press Enter for full details[white]`,
			panelArrow,
			navState.selectedPanel+" panel",
			r.Name,
			r.Timezone,
			offsetStr,
		)
	}

	NavigationView.SetText(detailsText)

	// Set border color based on selection
	if navState.selectedPanel == "left" {
		NavigationView.SetBorderColor(tcell.ColorGreen)
	} else {
		NavigationView.SetBorderColor(tcell.ColorOrange)
	}
}

// formatBool returns Yes/No string for boolean
func formatBool(b bool) string {
	if b {
		return "[green]Yes[white]"
	}
	return "[darkgray]No[white]"
}

// GetNavigationHelpText returns help text for navigation mode
func GetNavigationHelpText() string {
	if navState.selectedIndex < 0 {
		return "[darkgray]Keys:[white] ↑/↓=Navigate  Enter=Select  Tab=Switch Panel  Q/Esc=Quit"
	}
	return "[darkgray]Keys:[white] ↑/↓=Navigate  Enter=Toggle Details  Tab=Switch Panel  Esc=Deselect"
}
