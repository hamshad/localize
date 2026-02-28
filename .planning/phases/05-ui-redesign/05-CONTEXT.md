# Phase 05: Full-Screen Map UI Redesign — Context

## Vision

Transform Localize from a split-panel dashboard into a **full-screen world map display** designed as a **Linux ricing showpiece** with in-app features accessible via overlay modals. The map should dominate the entire terminal. Features are tucked behind a menu system that opens centered overlay boxes on top of the map.

## Decisions

### Map (LOCKED)
- Higher-resolution bitmap source (more detailed coastlines, islands) — roughly 480x192 → 240x48 braille
- Dynamic scaling to fill whatever terminal size the user has
- Both: better source data AND runtime scaling to terminal dimensions

### Status Bar (LOCKED)
- Minimal 1-row bar at top or bottom
- Shows ONLY: city date and time in **12-hour format**
- No mode indicators, no verbose headers

### City Markers on Map (LOCKED)
- Each city shows its **3-letter abbreviation + current time** directly on the map
- Selected city gets a **pulsing/blinking dot effect**
- Use the City struct's Coordinates field for proper lat/lon → map position mapping

### Clock Panels (LOCKED)
- Removed from always-visible layout
- Available as a menu feature (opens as overlay box)
- Also show as subtle text overlaid on map edges (dual approach: overlay + on-map)

### Menu System (LOCKED)
- Triggered by **M, Space, or Enter** keys
- Shows a visible `[=]` indicator in a corner (no other hint text)
- Opens a menu listing all features: Converter, Stopwatch, Timer, Alarm, Meeting Planner, Clock List
- No minimal hint bar, no footer help text in normal view

### Feature Overlays (LOCKED)
- Centered modal box (popup dialog) over the map
- Each feature (Timer, Converter, Stopwatch, Alarm, Meeting) opens in its own modal
- Include an **ASCII analog clock** in relevant overlays for visual richness

### Layout (LOCKED)
- Map fills entire terminal (no fixed-height panels stealing rows)
- Remove: header bar, clock panels, navigation panel, mode panel, footer
- Replace with: fullscreen map + overlay system

## Deferred Ideas

- Semi-transparent overlays (tview doesn't support true transparency)
- Slide-in panels from sides
- Mouse interaction

## Claude's Discretion

- Exact ASCII analog clock design/size
- Menu box styling and positioning
- Pulsing effect implementation (color cycling vs character alternation)
- How city time labels avoid overlapping on the map
- Exact braille scaling algorithm (nearest-neighbor vs bilinear for bitmap scaling)
- Status bar position (top vs bottom) — lean toward bottom
