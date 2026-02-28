package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// stopwatchState holds the current stopwatch state.
type stopwatchState struct {
	running   bool
	startTime time.Time
	elapsed   time.Duration // Accumulated elapsed time when paused
	laps      []time.Duration
}

// stopwatchMode handles the stopwatch functionality.
type stopwatchMode struct {
	state  *stopwatchState
	ticker *time.Ticker
	stopCh chan struct{}
}

// newStopwatchMode creates a new stopwatch mode handler.
func newStopwatchMode() *stopwatchMode {
	return &stopwatchMode{
		state:  &stopwatchState{},
		stopCh: make(chan struct{}),
	}
}

// GetMode returns the mode type.
func (s *stopwatchMode) GetMode() Mode {
	return ModeStopwatch
}

// HandleKey handles key events in stopwatch mode.
func (s *stopwatchMode) HandleKey(key rune) bool {
	switch key {
	case 27: // Escape
		// Stop the stopwatch in background
		return false
	case ' ', 's', 'S':
		// Toggle start/pause
		if s.state.running {
			s.pause()
		} else {
			s.start()
		}
		return true
	case 'r', 'R':
		// Reset
		s.reset()
		return true
	case 'l', 'L':
		// Lap
		if s.state.running {
			s.lap()
		}
		return true
	}
	return false
}

// start starts or resumes the stopwatch.
func (s *stopwatchMode) start() {
	if s.state.running {
		return
	}
	s.state.running = true
	s.state.startTime = time.Now()

	// Stop any previous goroutine
	if s.ticker != nil {
		s.ticker.Stop()
	}
	// Create a new stop channel for this run
	s.stopCh = make(chan struct{})

	// Start ticker for UI updates
	s.ticker = time.NewTicker(10 * time.Millisecond)
	stopCh := s.stopCh // capture for goroutine
	go func() {
		for {
			select {
			case <-s.ticker.C:
				// UI updates handled by main ticker
			case <-stopCh:
				return
			}
		}
	}()
}

// pause pauses the stopwatch.
func (s *stopwatchMode) pause() {
	if !s.state.running {
		return
	}
	s.state.running = false
	s.state.elapsed += time.Since(s.state.startTime)
	if s.ticker != nil {
		s.ticker.Stop()
	}
	// Signal the goroutine to stop
	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
}

// reset resets the stopwatch to zero.
func (s *stopwatchMode) reset() {
	s.state.running = false
	s.state.elapsed = 0
	s.state.startTime = time.Time{}
	s.state.laps = nil
	if s.ticker != nil {
		s.ticker.Stop()
	}
	// Signal the goroutine to stop
	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
}

// lap records a lap time.
func (s *stopwatchMode) lap() {
	elapsed := s.getElapsed()
	s.state.laps = append(s.state.laps, elapsed)
}

// getElapsed returns the current elapsed time.
func (s *stopwatchMode) getElapsed() time.Duration {
	if s.state.running {
		return s.state.elapsed + time.Since(s.state.startTime)
	}
	return s.state.elapsed
}

// formatDuration formats a duration as HH:MM:SS.ss.
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	hundredths := int(d.Milliseconds()/10) % 100
	return fmt.Sprintf("%02d:%02d:%02d.%02d", hours, minutes, seconds, hundredths)
}

// Render returns the rendered stopwatch display.
func (s *stopwatchMode) Render() string {
	var b strings.Builder

	// Header
	b.WriteString("\n[yellow::b]━━━ STOPWATCH ━━━[-::-]\n\n")

	// Main stopwatch display
	elapsed := s.getElapsed()
	status := "[green]Running[white]"
	if !s.state.running && elapsed > 0 {
		status = "[yellow]Paused[white]"
	} else if !s.state.running && elapsed == 0 {
		status = "[darkgray]Ready[white]"
	}

	b.WriteString(fmt.Sprintf("  [::b]%s[-]\n\n", formatDuration(elapsed)))
	b.WriteString(fmt.Sprintf("  Status: %s\n\n", status))

	// Controls help
	b.WriteString("  [darkgray]Controls:[white] Space=Start/Pause  R=Reset  L=Lap  Esc=Exit\n\n")

	// Laps
	if len(s.state.laps) > 0 {
		b.WriteString("  [aqua::b]Lap Times:[-::-]\n")
		for i, lap := range s.state.laps {
			b.WriteString(fmt.Sprintf("    Lap %d: %s\n", i+1, formatDuration(lap)))
		}
	}

	return b.String()
}

// GetHelpText returns the help text for stopwatch mode.
func (s *stopwatchMode) GetHelpText() string {
	return "[darkgray]Keys:[white] Space=Start/Pause  R=Reset  L=Lap  Esc=Exit (continues)"
}

// HandleSpecialKeyEvent handles non-rune key events (Enter, Backspace, etc.).
func (s *stopwatchMode) HandleSpecialKeyEvent(key tcell.Key) bool {
	// Stopwatch doesn't use Enter or Backspace
	return false
}

// Stop stops the background ticker.
func (s *stopwatchMode) Stop() {
	close(s.stopCh)
	if s.ticker != nil {
		s.ticker.Stop()
	}
}

// IsRunning returns whether the stopwatch is running.
func (s *stopwatchMode) IsRunning() bool {
	return s.state.running
}

// GetElapsed returns the current elapsed time.
func (s *stopwatchMode) GetElapsed() time.Duration {
	return s.getElapsed()
}
