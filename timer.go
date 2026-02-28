package main

import (
	"fmt"
	"strings"
	"time"
)

// timerState holds the current timer state.
type timerState struct {
	running        bool
	endTime        time.Time // When the timer should end
	inputDur       string    // User input duration string
	inputMode      bool      // True if waiting for duration input
	duration       time.Duration
	alarmTriggered bool
	flashState     bool // For flashing effect
}

// timerMode handles the countdown timer functionality.
type timerMode struct {
	state  *timerState
	ticker *time.Ticker
	stopCh chan struct{}
}

// newTimerMode creates a new timer mode handler.
func newTimerMode() *timerMode {
	return &timerMode{
		state:  &timerState{},
		stopCh: make(chan struct{}),
	}
}

// GetMode returns the mode type.
func (t *timerMode) GetMode() Mode {
	return ModeTimer
}

// HandleKey handles key events in timer mode.
func (t *timerMode) HandleKey(key rune) bool {
	switch key {
	case 27: // Escape
		return false
	case ' ', 's', 'S':
		// Toggle start/pause
		if t.state.running {
			t.pause()
		} else {
			t.start()
		}
		return true
	case 'r', 'R':
		// Reset
		t.reset()
		return true
	case 't', 'T':
		// Toggle input mode
		t.state.inputMode = !t.state.inputMode
		if !t.state.inputMode {
			// Try to parse input when exiting input mode
			if err := t.parseDuration(); err != nil {
				// Keep input if invalid
				t.state.inputMode = true
			}
		}
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if t.state.inputMode {
			if len(t.state.inputDur) < 8 { // Max HH:MM:SS
				t.state.inputDur += string(key)
				// Auto-add colon after 2 or 4 digits
				if len(t.state.inputDur) == 2 || len(t.state.inputDur) == 4 {
					t.state.inputDur += ":"
				}
			}
		}
		return true
	case ':':
		if t.state.inputMode && strings.Count(t.state.inputDur, ":") < 2 {
			if !strings.HasSuffix(t.state.inputDur, ":") {
				t.state.inputDur += ":"
			}
		}
		return true
	case 127: // Backspace
		if t.state.inputMode && len(t.state.inputDur) > 0 {
			t.state.inputDur = t.state.inputDur[:len(t.state.inputDur)-1]
			// Auto-remove colon if backspacing over it
			if len(t.state.inputDur) > 0 && len(t.state.inputDur)%3 == 2 {
				t.state.inputDur = t.state.inputDur[:len(t.state.inputDur)-1]
			}
		}
		return true
	}
	return false
}

// parseDuration tries to parse the input duration string.
func (t *timerMode) parseDuration() error {
	if t.state.inputDur == "" {
		return fmt.Errorf("empty input")
	}

	// Try parsing as HH:MM:SS
	d, err := time.ParseDuration(t.state.inputDur + "s")
	if err != nil {
		// Try parsing as MM:SS
		d, err = time.ParseDuration("0m" + t.state.inputDur + "s")
		if err != nil {
			return err
		}
	}

	if d <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	t.state.duration = d
	t.state.inputDur = "" // Clear input after successful parse
	return nil
}

// start starts or resumes the timer.
func (t *timerMode) start() {
	if t.state.running {
		return
	}

	// If no duration set, try to parse input
	if t.state.duration == 0 && t.state.inputDur != "" {
		if err := t.parseDuration(); err != nil {
			return
		}
	}

	if t.state.duration == 0 {
		return // No duration to start
	}

	t.state.running = true
	t.state.endTime = time.Now().Add(t.state.duration)
	t.state.alarmTriggered = false

	// Start ticker for UI updates
	t.ticker = time.NewTicker(100 * time.Millisecond)
	go func() {
		flashTicker := time.NewTicker(500 * time.Millisecond)
		flashOn := false
		for {
			select {
			case <-t.ticker.C:
				// Check if timer expired
				if time.Now().After(t.state.endTime) && !t.state.alarmTriggered {
					t.state.alarmTriggered = true
					t.state.running = false
					t.ticker.Stop()
				}
			case <-flashTicker.C:
				if t.state.alarmTriggered {
					flashOn = !flashOn
					t.state.flashState = flashOn
				}
			case <-t.stopCh:
				flashTicker.Stop()
				return
			}
		}
	}()
}

// pause pauses the timer.
func (t *timerMode) pause() {
	if !t.state.running {
		return
	}
	t.state.running = false
	// Calculate remaining duration
	remaining := time.Until(t.state.endTime)
	if remaining > 0 {
		t.state.duration = remaining
	}
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

// reset resets the timer to zero.
func (t *timerMode) reset() {
	t.state.running = false
	t.state.duration = 0
	t.state.inputDur = ""
	t.state.endTime = time.Time{}
	t.state.inputMode = false
	t.state.alarmTriggered = false
	t.state.flashState = false
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

// getRemaining returns the remaining time.
func (t *timerMode) getRemaining() time.Duration {
	if t.state.running {
		remaining := time.Until(t.state.endTime)
		if remaining > 0 {
			return remaining
		}
		return 0
	}
	return t.state.duration
}

// formatTimerDuration formats duration for timer display (HH:MM:SS).
func formatTimerDuration(d time.Duration) string {
	if d <= 0 {
		return "00:00:00"
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// Render returns the rendered timer display.
func (t *timerMode) Render() string {
	var b strings.Builder

	// Header
	b.WriteString("\n[yellow::b]━━━ COUNTDOWN TIMER ━━━[-::-]\n\n")

	// Main timer display
	remaining := t.getRemaining()

	// Flash effect when alarm triggered
	if t.state.alarmTriggered && t.state.flashState {
		b.WriteString(fmt.Sprintf("  [::b]%s[white]\n\n", "ALARM!"))
	} else if t.state.alarmTriggered {
		b.WriteString(fmt.Sprintf("  [::b]%s[white]\n\n", "ALARM!"))
	} else {
		b.WriteString(fmt.Sprintf("  [::b]%s[-]\n\n", formatTimerDuration(remaining)))
	}

	// Status
	if t.state.running {
		b.WriteString("  Status: [green]Running[white]\n")
	} else if t.state.alarmTriggered {
		b.WriteString("  Status: [red]Alarm Triggered![white]  Press any key\n")
	} else if remaining > 0 {
		b.WriteString("  Status: [yellow]Paused[white]\n")
	} else {
		b.WriteString("  Status: [darkgray]Ready[white]\n")
	}

	// Duration input
	if t.state.inputMode {
		displayDur := t.state.inputDur
		if displayDur == "" {
			displayDur = "HH:MM:SS"
		}
		b.WriteString(fmt.Sprintf("\n  [::b]Enter duration: [yellow]%s[-]\n", displayDur))
	} else if remaining == 0 && !t.state.alarmTriggered {
		b.WriteString("\n  [darkgray]Press T to set duration[white]\n")
	}

	// Controls
	b.WriteString("\n  [darkgray]Controls:[white] T=Set Time  Space=Start/Pause  R=Reset  Esc=Exit\n")

	return b.String()
}

// GetHelpText returns the help text for timer mode.
func (t *timerMode) GetHelpText() string {
	return "[darkgray]Keys:[white] T=Set Duration  Space=Start/Pause  R=Reset  Esc=Exit"
}

// Stop stops the background ticker.
func (t *timerMode) Stop() {
	close(t.stopCh)
	if t.ticker != nil {
		t.ticker.Stop()
	}
}
