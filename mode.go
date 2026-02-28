package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Mode represents the different interaction modes in the app.
type Mode int

const (
	ModeNormal Mode = iota
	ModeConverter
	ModeStopwatch
	ModeTimer
	ModeAlarm
	ModeNavigation // Not used yet, for future map navigation
	ModeMeeting
)

// modeNames provides human-readable names for each mode.
var modeNames = map[Mode]string{
	ModeNormal:     "Normal",
	ModeConverter:  "Converter",
	ModeStopwatch:  "Stopwatch",
	ModeTimer:      "Timer",
	ModeAlarm:      "Alarm",
	ModeNavigation: "Navigation",
	ModeMeeting:    "Meeting",
}

// ModeHandler defines the interface for mode-specific behavior.
type ModeHandler interface {
	GetMode() Mode
	HandleKey(key rune) bool
	HandleSpecialKeyEvent(key tcell.Key) bool // Handle non-rune keys (Enter, Backspace, etc.)
	Render() string
	GetHelpText() string
}

// modeManager handles switching between different modes.
type modeManager struct {
	currentMode   Mode
	previousMode  Mode
	handlers      map[Mode]ModeHandler
	modeIndicator *tview.TextView
}

// newModeManager creates a new mode manager.
func newModeManager() *modeManager {
	mm := &modeManager{
		currentMode:   ModeNormal,
		previousMode:  ModeNormal,
		handlers:      make(map[Mode]ModeHandler),
		modeIndicator: tview.NewTextView(),
	}
	mm.modeIndicator.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)
	mm.updateModeIndicator()
	return mm
}

// RegisterHandler registers a mode handler for a specific mode.
func (mm *modeManager) RegisterHandler(mode Mode, handler ModeHandler) {
	mm.handlers[mode] = handler
}

// GetCurrentMode returns the current mode.
func (mm *modeManager) GetCurrentMode() Mode {
	return mm.currentMode
}

// SwitchTo switches to a new mode, returning true if successful.
func (mm *modeManager) SwitchTo(mode Mode) bool {
	// Don't switch if already in that mode
	if mm.currentMode == mode {
		return true
	}
	mm.previousMode = mm.currentMode
	mm.currentMode = mode
	mm.updateModeIndicator()
	return true
}

// SwitchToPrevious returns to the previous mode.
func (mm *modeManager) SwitchToPrevious() bool {
	if mm.previousMode != mm.currentMode {
		mm.currentMode, mm.previousMode = mm.previousMode, mm.currentMode
		mm.updateModeIndicator()
		return true
	}
	return false
}

// Escape returns to Normal mode from any mode.
func (mm *modeManager) Escape() bool {
	if mm.currentMode != ModeNormal {
		mm.previousMode = mm.currentMode
		mm.currentMode = ModeNormal
		mm.updateModeIndicator()
		return true
	}
	return false
}

// HandleKey delegates key handling to the current mode's handler.
func (mm *modeManager) HandleKey(key rune) bool {
	if handler, ok := mm.handlers[mm.currentMode]; ok {
		return handler.HandleKey(key)
	}
	return false
}

// HandleSpecialKeyEvent delegates non-rune key events (Enter, Backspace, etc.)
// to the current mode's handler.
func (mm *modeManager) HandleSpecialKeyEvent(key tcell.Key) bool {
	if handler, ok := mm.handlers[mm.currentMode]; ok {
		return handler.HandleSpecialKeyEvent(key)
	}
	return false
}

// HandleSpecialKey delegates special key handling to the current mode's handler.
func (mm *modeManager) HandleSpecialKey(key tcell.Key) bool {
	if handler, ok := mm.handlers[mm.currentMode]; ok {
		// Check if handler implements special key handling
		if h, ok := handler.(*MeetingMode); ok {
			return h.planner.HandleSpecialKey(key)
		}
	}
	return false
}

// Render returns the rendered content for the current mode.
func (mm *modeManager) Render() string {
	if handler, ok := mm.handlers[mm.currentMode]; ok {
		return handler.Render()
	}
	return ""
}

// GetModeIndicator returns the mode indicator text view.
func (mm *modeManager) GetModeIndicator() *tview.TextView {
	return mm.modeIndicator
}

// GetHelpText returns the help text for the current mode.
func (mm *modeManager) GetHelpText() string {
	if handler, ok := mm.handlers[mm.currentMode]; ok {
		return handler.GetHelpText()
	}
	return getNormalModeHelp()
}

// updateModeIndicator updates the mode indicator display.
func (mm *modeManager) updateModeIndicator() {
	var status string
	if mm.currentMode == ModeNormal {
		status = "[darkgray]Normal[white]"
	} else {
		status = fmt.Sprintf("[yellow]%s[white]", modeNames[mm.currentMode])
	}
	mm.modeIndicator.SetText(status)
}

// getNormalModeHelp returns the help text for Normal mode.
func getNormalModeHelp() string {
	if IsNavigationActive() {
		return GetNavigationHelpText()
	}
	return "[darkgray]Keys:[white] C=Converter  S=Stopwatch  T=Timer  A=Alarm  D=Day/Night  ↑/↓=Navigate  Q/Esc=Quit"
}
