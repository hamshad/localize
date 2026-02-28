package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Alarm represents a single alarm.
type Alarm struct {
	ID       string `json:"id"`
	Time     string `json:"time"` // HH:MM format
	Timezone string `json:"timezone"`
	Repeat   string `json:"repeat"` // "once", "daily", "weekday"
	Enabled  bool   `json:"enabled"`
	CityName string `json:"city_name"` // Display name for the city
}

// AlarmConfig holds all alarms.
type AlarmConfig struct {
	Alarms []Alarm `json:"alarms"`
}

// alarmMode handles the alarm functionality.
type alarmMode struct {
	config         *AlarmConfig
	inputMode      string // "none", "add"
	inputStep      int    // 0=select timezone, 1=enter time, 2=select repeat
	currentZone    int
	selectedZone   string
	inputTime      string
	selectedRepeat string
	ticker         *time.Ticker
	stopCh         chan struct{}
	triggered      map[string]bool // Tracks which alarms have triggered this minute
}

// newAlarmMode creates a new alarm mode handler.
func newAlarmMode() *alarmMode {
	am := &alarmMode{
		config:    &AlarmConfig{},
		inputMode: "none",
		stopCh:    make(chan struct{}),
		triggered: make(map[string]bool),
	}
	am.loadAlarms()
	return am
}

// GetMode returns the mode type.
func (am *alarmMode) GetMode() Mode {
	return ModeAlarm
}

// HandleKey handles key events in alarm mode.
func (am *alarmMode) HandleKey(key rune) bool {
	switch key {
	case 27: // Escape
		am.inputMode = "none"
		am.inputStep = 0
		am.inputTime = ""
		return false
	case 'a', 'A':
		// Start adding new alarm
		if am.inputMode == "none" {
			am.inputMode = "add"
			am.inputStep = 0
			am.selectedZone = ""
			am.inputTime = ""
			am.selectedRepeat = ""
			am.currentZone = 0
		}
		return true
	case 'd', 'D':
		// Delete alarm - simple implementation: delete last alarm
		if len(am.config.Alarms) > 0 && am.inputMode == "none" {
			am.config.Alarms = am.config.Alarms[:len(am.config.Alarms)-1]
			am.saveAlarms()
		}
		return true
	case ' ', 's', 'S':
		// Toggle start/pause
		return true
	case 'j', 'J':
		// Navigate down in selection lists
		if am.inputMode == "add" && am.inputStep == 0 {
			allZones := am.getAllZones()
			if am.currentZone < len(allZones)-1 {
				am.currentZone++
			}
		}
		return true
	case 'k', 'K':
		// Navigate up in selection lists
		if am.inputMode == "add" && am.inputStep == 0 {
			if am.currentZone > 0 {
				am.currentZone--
			}
		}
		return true
	case 13: // Enter
		if am.inputMode == "add" {
			am.handleEnterKey()
		}
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Handle digit input
		if am.inputMode == "add" && am.inputStep == 2 {
			// Repeat selection: 1=once, 2=daily, 3=weekday
			switch key {
			case '1':
				am.selectedRepeat = "once"
				am.addAlarm()
			case '2':
				am.selectedRepeat = "daily"
				am.addAlarm()
			case '3':
				am.selectedRepeat = "weekday"
				am.addAlarm()
			}
		} else if am.inputMode == "add" && am.inputStep == 1 {
			// Time input
			if len(am.inputTime) < 5 {
				am.inputTime += string(key)
				if len(am.inputTime) == 2 {
					am.inputTime += ":"
				}
			}
		}
		return true
	case ':':
		if am.inputMode == "add" && am.inputStep == 1 && !strings.Contains(am.inputTime, ":") {
			am.inputTime += ":"
		}
		return true
	case 127: // Backspace
		if am.inputMode == "add" && am.inputStep == 1 && len(am.inputTime) > 0 {
			am.inputTime = am.inputTime[:len(am.inputTime)-1]
			if len(am.inputTime) == 2 && am.inputTime[1] == ':' {
				am.inputTime = am.inputTime[:1]
			}
		}
		return true
	}
	return false
}

// handleEnterKey handles Enter key in add mode.
func (am *alarmMode) handleEnterKey() {
	switch am.inputStep {
	case 0:
		// Select timezone - move to time input
		allZones := am.getAllZones()
		if am.currentZone < len(allZones) {
			am.selectedZone = allZones[am.currentZone].Timezone
			am.inputStep = 1
		}
	case 1:
		// Enter time - validate and move to repeat selection
		if len(am.inputTime) >= 4 {
			am.inputStep = 2
		}
	}
}

// getAllZones returns all available zones.
func (am *alarmMode) getAllZones() []Region {
	var allZones []Region
	allZones = append(allZones, leftRegions...)
	allZones = append(allZones, rightRegions...)
	return allZones
}

// addAlarm adds a new alarm.
func (am *alarmMode) addAlarm() {
	alarm := Alarm{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Time:     am.inputTime,
		Timezone: am.selectedZone,
		Repeat:   am.selectedRepeat,
		Enabled:  true,
		CityName: am.getCityForZone(am.selectedZone),
	}
	am.config.Alarms = append(am.config.Alarms, alarm)
	am.saveAlarms()

	// Reset input mode
	am.inputMode = "none"
	am.inputStep = 0
	am.inputTime = ""
	am.selectedZone = ""
	am.selectedRepeat = ""
}

// getCityForZone returns the city name for a timezone.
func (am *alarmMode) getCityForZone(tz string) string {
	for _, r := range leftRegions {
		if r.Timezone == tz {
			return r.Name
		}
	}
	for _, r := range rightRegions {
		if r.Timezone == tz {
			return r.Name
		}
	}
	return tz
}

// loadAlarms loads alarms from the config file.
func (am *alarmMode) loadAlarms() {
	configPath := alarmConfigPath()
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, am.config)
	}
}

// saveAlarms saves alarms to the config file.
func (am *alarmMode) saveAlarms() {
	data, _ := json.MarshalIndent(am.config, "", "  ")
	os.WriteFile(alarmConfigPath(), data, 0644)
}

// alarmConfigPath returns the path to the alarm config file.
func alarmConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".localize", "alarms.json")
}

// CheckAlarms checks if any alarms should trigger.
func (am *alarmMode) CheckAlarms() []Alarm {
	var triggered []Alarm
	now := time.Now()

	for i, alarm := range am.config.Alarms {
		if !alarm.Enabled {
			continue
		}

		// Check if alarm time matches current time
		alarmTime := alarm.Time
		nowTime := now.Format("15:04")

		// Create a unique key for this alarm + minute
		checkKey := fmt.Sprintf("%s-%s", alarm.ID, nowTime)

		// Skip if already triggered this minute
		if am.triggered[checkKey] {
			continue
		}

		// Check if alarm should trigger
		shouldTrigger := false
		if alarmTime == nowTime {
			switch alarm.Repeat {
			case "once":
				shouldTrigger = true
				// Disable after triggering once
				am.config.Alarms[i].Enabled = false
			case "daily":
				shouldTrigger = true
			case "weekday":
				// 0 = Sunday, 6 = Saturday
				if now.Weekday() >= 1 && now.Weekday() <= 5 {
					shouldTrigger = true
				}
			}
		}

		if shouldTrigger {
			triggered = append(triggered, alarm)
			am.triggered[checkKey] = true
		}
	}

	// Save if we disabled any alarms
	if len(triggered) > 0 {
		am.saveAlarms()
	}

	return triggered
}

// GetTriggeredAlarms returns currently triggered alarms for display.
func (am *alarmMode) GetTriggeredAlarms() []Alarm {
	var triggered []Alarm
	now := time.Now()
	nowTime := now.Format("15:04")

	for _, alarm := range am.config.Alarms {
		if !alarm.Enabled {
			continue
		}
		if alarm.Time == nowTime {
			triggered = append(triggered, alarm)
		}
	}
	return triggered
}

// Render returns the rendered alarm display.
func (am *alarmMode) Render() string {
	var b strings.Builder

	// Header
	b.WriteString("\n[yellow::b]━━━ ALARMS ━━━[-::-]\n\n")

	// Check for triggered alarms
	triggered := am.GetTriggeredAlarms()
	if len(triggered) > 0 {
		b.WriteString("  [red::b]⚠ ALARM! ⚠[-::-]\n")
		for _, alarm := range triggered {
			b.WriteString(fmt.Sprintf("  [red]%s @ %s (%s)[-]\n",
				alarm.Time, alarm.CityName, alarm.Repeat))
		}
		b.WriteString("\n  [darkgray]Press any key to dismiss[white]\n\n")
	}

	// Input mode for adding alarm
	if am.inputMode == "add" {
		b.WriteString("  [aqua::b]Add New Alarm:[-::-]\n\n")
		allZones := am.getAllZones()

		switch am.inputStep {
		case 0: // Select timezone
			b.WriteString("  [darkgray]Select city/timezone: (↑/↓ to navigate, Enter to select)[-]\n\n")
			for i, zone := range allZones {
				marker := "  "
				if i == am.currentZone {
					marker = "> "
				}
				b.WriteString(fmt.Sprintf("%s%s (%s)\n", marker, zone.Name, zone.Timezone))
			}
		case 1: // Enter time
			b.WriteString(fmt.Sprintf("  City: %s\n", am.getCityForZone(am.selectedZone)))
			displayTime := am.inputTime
			if displayTime == "" {
				displayTime = "HH:MM"
			}
			b.WriteString(fmt.Sprintf("  [::b]Enter time: [yellow]%s[-]\n\n", displayTime))
		case 2: // Select repeat
			b.WriteString(fmt.Sprintf("  City: %s, Time: %s\n\n", am.getCityForZone(am.selectedZone), am.inputTime))
			b.WriteString("  [::b]Select repeat:[-]\n")
			b.WriteString("    1) Once\n")
			b.WriteString("    2) Daily\n")
			b.WriteString("    3) Weekday (Mon-Fri)\n")
		}

		b.WriteString("\n  [darkgray]Press Esc to cancel[white]\n")
		return b.String()
	}

	// Display alarms
	if len(am.config.Alarms) == 0 {
		b.WriteString("  [darkgray]No alarms set.[white]\n")
		b.WriteString("  Press A to add a new alarm.\n")
	} else {
		b.WriteString(fmt.Sprintf("  [aqua::b]Active Alarms (%d):[-::-]\n\n", len(am.config.Alarms)))
		for _, alarm := range am.config.Alarms {
			enabled := "[green]●[white]"
			if !alarm.Enabled {
				enabled = "[darkgray]○[white]"
			}
			repeatLabel := alarm.Repeat
			if repeatLabel == "weekday" {
				repeatLabel = "weekdays"
			}
			b.WriteString(fmt.Sprintf("  %s %s @ %s (%s) [%s]\n",
				enabled, alarm.Time, alarm.CityName, alarm.Timezone, repeatLabel))
		}
	}

	b.WriteString("\n  [darkgray]Controls:[white] A=Add  D=Delete Last  Esc=Exit\n")

	return b.String()
}

// GetHelpText returns the help text for alarm mode.
func (am *alarmMode) GetHelpText() string {
	return "[darkgray]Keys:[white] A=Add Alarm  D=Delete Last  Esc=Exit"
}

// Stop stops background processes.
func (am *alarmMode) Stop() {
	close(am.stopCh)
	if am.ticker != nil {
		am.ticker.Stop()
	}
}
