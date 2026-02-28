---
phase: 03-visual
plan: 01
type: execute
wave: 1
depends_on: []
files_modified: [main.go, daynight.go (new)]
autonomous: true
requirements: [F6, F7, F10]
---

<objective>
Implement world day/night overlay, relative time display, and keyboard navigation.
</objective>

<context>
@main.go
@.planning/phases/01-foundation/01-foundation-SUMMARY.md
</context>

<tasks>

<task type="auto">
  <name>Implement World Day/Night Overlay</name>
  <files>daynight.go (new), main.go</files>
  <action>
    Modify the braille map rendering to show day/night:
    - Calculate current UTC hour for each longitude position
    - Daytime (6AM-6PM local): render with full brightness green
    - Nighttime (6PM-6AM local): render with dimmed color (dark gray/blue)
    - Twilight (6-7AM, 5-6PM): render with amber tint
    - Add toggle: press 'D' to enable/disable day/night overlay
    - Header shows current overlay state
  </action>
  <verify>Press D toggles day/night overlay; map shows light/dark regions based on time</verify>
  <done>Visual representation of day/night on world map</done>
</task>

<task type="auto">
  <name>Add Relative Time Display</name>
  <files>main.go</files>
  <action>
    Add relative time offset display next to each city:
    - Determine user's local timezone (or first city in list as reference)
    - Show "+5h", "-3h", "+0h" next to each city clock
    - Color code: green for same day, yellow for different day (+/- 1 day)
    - Format: "New York [white::b]14:30[-::-] [silver]+0h"
    - First city in list is the "reference" (shows +0h)
  </action>
  <verify>Each city shows relative time offset from reference city</verify>
  <done>Users can quickly see time differences at a glance</done>
</task>

<task type="auto">
  <name>Implement Keyboard Navigation & City Selection</name>
  <files>navigation.go (new), main.go</files>
  <action>
    Implement keyboard navigation:
    - Press arrow keys (Up/Down) to select cities in the clock panel
    - Selected city gets highlighted border/color
    - Press Enter on selected city to show details panel:
      - Full timezone name (e.g., "America/New_York")
      - UTC offset (+05:30 format)
      - Day of week
      - Is DST active?
    - Selected city indicator persists across mode switches
    - Press Escape to deselect
  </action>
  <verify>Arrow keys navigate between cities, Enter shows details, Escape closes</verify>
  <done>Full keyboard navigation for city selection and details viewing</done>
</task>

</tasks>

<verification>
- D key toggles day/night overlay on map
- Each city clock shows relative offset from reference
- Arrow keys navigate cities, Enter shows details
- Navigation works in all modes (Normal, Converter, etc.)
</verification>

<success_criteria>
World map shows day/night regions, users can see relative times, and full keyboard navigation is available.
</success_criteria>

<output>
After completion, create .planning/phases/03-visual/03-visual-SUMMARY.md
</output>
