# Localize - Project Roadmap

## Overview

**Project:** localize - Terminal World Clock Dashboard  
**Type:** Terminal Application (Go/tview)  
**Goal:** Transform the static world clock into a fully configurable, feature-rich timezone management tool

## Current State

- Working terminal app displaying 16 cities across major timezones
- Braille-rendered world map with city markers
- Real-time clock updates every second
- Day/night phase indicators

---

## Phase Progress

| Phase | Name | Plans | Status | Completed |
|-------|------|-------|--------|-----------|
| 01 | Foundation & Configuration | 1/1 | Complete | 2026-02-28 |
| 02 | Core Utility Features | 0/1 | Not Started | - |
| 03 | Map & Visual Enhancements | 0/1 | Not Started | - |
| 04 | Planning & Expansion | 0/1 | Not Started | - |

---

## Requirements

### Phase 1: Foundation & Configuration
**Requirements:** [F1, F9]

1. **F1:** User-Configurable Cities — Allow users to specify which cities to display via CLI flags or config file
2. **F9:** Config Persistence — Save user preferences (selected cities, layout) to a config file

### Phase 2: Core Utility Features
**Requirements:** [F3, F4, F5]

3. **F3:** Time Converter — Input a time in one timezone, display what time it is in all other cities
4. **F4:** Stopwatch/Timer — Add stopwatch and countdown timer with keyboard shortcuts
5. **F5:** Alarm Functionality — Set alarms in any timezone with visual/audio notification

### Phase 3: Map & Visual Enhancements
**Requirements:** [F6, F7, F10]

6. **F6:** World Day/Night Overlay — Visually show which regions are in daylight vs nighttime on the map
7. **F7:** Relative Time Display — Show time in other cities relative to user's local timezone
8. **F10:** Keyboard Navigation — Select/highlight cities to see detailed information

### Phase 4: Planning & Expansion
**Requirements:** [F2, F8]

9. **F2:** Meeting Time Planner — Highlight overlapping business hours across selected timezones
10. **F8:** Add More Cities — Expand city list with additional business hubs

---

## Feature Details

### F1: User-Configurable Cities
- **Description:** Allow specifying which cities to display via CLI flags (-cities) or config file
- **Priority:** HIGH (foundation for personalization)
- **Input:** Comma-separated city names or preset groups (e.g., "business", "family")

### F9: Config Persistence
- **Description:** Save/load user preferences to ~/.localize/config.json
- **Priority:** HIGH (enables F1 persistence)
- **Storage:** JSON config file with selected cities, layout preferences

### F3: Time Converter
- **Description:** Interactive mode where user enters a time + source timezone, displays converted times
- **Priority:** MEDIUM
- **Input:** Keyboard input for time and timezone selection

### F4: Stopwatch/Timer
- **Description:** Full-featured stopwatch and countdown timer with start/stop/reset
- **Priority:** MEDIUM
- **Controls:** Keyboard shortcuts (S=start, P=pause, R=reset, T=timer mode)

### F5: Alarm Functionality
- **Description:** Set alarms in any timezone with terminal bell notification
- **Priority:** MEDIUM
- **Storage:** Saved in config, persistent across sessions

### F6: World Day/Night Overlay
- **Description:** Calculate sun position and color-code map regions (light=destination, dark=night)
- **Priority:** MEDIUM
- **Algorithm:** Simple day/night based on UTC hour per longitude

### F7: Relative Time Display
- **Description:** Show +/- hours difference from local timezone next to each city clock
- **Priority:** LOW
- **Format:** "+5h", "-3h", "same"

### F10: Keyboard Navigation
- **Description:** Arrow keys to select cities, Enter for details, Esc to exit
- **Priority:** MEDIUM
- **Details Panel:** Shows additional info (UTC offset, day of week, etc.)

### F2: Meeting Time Planner
- **Description:** Visual highlight of overlapping business hours (9AM-5PM) across selected zones
- **Priority:** LOW
- **Display:** Color-coded "green" when all selected cities are within business hours

### F8: Add More Cities
- **Description:** Expand city list to 30+ major business hubs
- **Priority:** LOW
- **Categories:** Americas, Europe, Asia, Africa, Oceania, Middle East

---

## Technical Notes

### Dependencies
- F1 → F9 (config needed for persistence)
- F6 uses F1 (day/night calc per configured cities)
- F2 uses F1 (meeting planner operates on selected cities)

### Architecture
- All features add to existing tview layout
- New keyboard modes (navigation, converter, timer)
- Config struct for persistence

---

## Success Criteria

1. Users can customize displayed cities and have preferences persist
2. All utility features (converter, timer, alarm) work reliably
3. Map enhancements are visually clear
4. Keyboard navigation is intuitive
5. No performance degradation from additional features
