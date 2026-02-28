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
| 01 | 1/1 | Complete    | 2026-02-28 | 2026-02-28 |
| 02 | 1/1 | Complete    | 2026-02-28 | 2026-02-28 |
| 03 | 1/1 | Complete    | 2026-02-28 | 2026-02-28 |
| 04 | 1/1 | Complete    | 2026-02-28 | 2026-02-28 |
| 05 | Full-Screen Map UI Redesign | 3 plans | In Progress | - |

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
**Requirements:** [F2, F8] ✓ COMPLETE

9. **F2:** Meeting Time Planner — Highlight overlapping business hours across selected timezones ✓
10. **F8:** Add More Cities — Expand city list with additional business hubs ✓

### Phase 5: Full-Screen Map UI Redesign
**Requirements:** [UI-01, UI-02, UI-03, UI-04, UI-05, UI-06, UI-07, UI-08]
**Goal:** Transform the app into a fullscreen world map ricing showpiece with overlay-based feature access

11. **UI-01:** High-Resolution Map — Create detailed 480x192 bitmap with better coastlines, islands, peninsulas ✓
12. **UI-02:** Dynamic Map Scaling — Braille map scales to fill any terminal size ✓
13. **UI-03:** Fullscreen Layout — Remove all panels, map fills entire terminal
14. **UI-04:** City Time Markers — Show 3-letter city codes with live 12hr time on the map
15. **UI-05:** Pulsing City Selection — Selected city marker blinks/pulses, minimal 1-row status bar with 12hr format
16. **UI-06:** Menu Overlay System — M/Space/Enter opens centered menu listing all features, [=] indicator in corner
17. **UI-07:** Feature Modal Boxes — Each feature opens as centered popup over the map
18. **UI-08:** ASCII Analog Clock — Rendered inside relevant feature overlays for visual richness

Plans:
- [x] 05-01-PLAN.md — High-res map bitmap + dynamic scaling engine
- [ ] 05-02-PLAN.md — Fullscreen layout, city markers with time, status bar
- [ ] 05-03-PLAN.md — Menu overlay system, feature modals, ASCII clock

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
