# Localize - Project Specification

## Basic Info

- **Project Name:** localize
- **Type:** Terminal Application (Go)
- **Language:** Go 1.25.6

## Tech Stack

- **UI Framework:** tview (github.com/rivo/tview)
- **Terminal:** tcell (github.com/gdamore/tcell/v2)
- **Time Handling:** Go standard library (time package with IANA timezone names)

## Current Codebase

**main.go (449 lines):**
- Terminal UI with tview layout
- Braille-rendered world map (240x96 bitmap)
- 16 pre-configured cities across Americas, Europe, Asia, Africa, Oceania
- Real-time clock updates (1 second interval)
- Day/night phase indicators

### Key Components

| Component | Purpose |
|-----------|---------|
| `Region` struct | City name, timezone, display color |
| `getBrailleWorldMap()` | Renders ASCII braille world map |
| `colorizeBrailleMap()` | Adds colors + city markers |
| `formatClockPanel()` | Displays formatted time for each region |
| `getDayPhase()` | Returns dawn/morning/afternoon/evening/dusk/night |

### City Data (Current)

**Americas & Europe:**
- Honolulu (Pacific/Honolulu)
- Anchorage (America/Anchorage)
- Los Angeles (America/Los_Angeles)
- New York (America/New_York)
- Sao Paulo (America/Sao_Paulo)
- London (Europe/London)
- Paris (Europe/Paris)
- Moscow (Europe/Moscow)

**Asia, Africa & Oceania:**
- Cairo (Africa/Cairo)
- Nairobi (Africa/Nairobi)
- Dubai (Asia/Dubai)
- Mumbai (Asia/Kolkata)
- Singapore (Asia/Singapore)
- Shanghai (Asia/Shanghai)
- Tokyo (Asia/Tokyo)
- Sydney (Australia/Sydney)

## Dependencies

```
github.com/gdamore/tcell/v2    v2.13.8
github.com/rivo/tview          v0.42.0
```

## Build & Run

```bash
go run main.go
# Press Q or Esc to quit
```

## Project Characteristics

- **Brownfield:** Yes (existing working codebase)
- **Testing:** None currently
- **Config:** None (hardcoded cities)
- **CLI Flags:** None currently
