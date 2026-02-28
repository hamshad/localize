---
phase: 04-expansion
plan: 01
type: execute
wave: 1
depends_on: []
files_modified: [main.go, cities.go]
autonomous: true
requirements: [F2, F8]
---

<objective>
Implement meeting time planner and expand city list.
</objective>

<context>
@main.go
@.planning/phases/01-foundation/01-foundation-SUMMARY.md
@.planning/phases/03-visual/03-visual-SUMMARY.md
</context>

<tasks>

<task type="auto">
  <name>Expand city list to 30+ cities</name>
  <files>cities.go (new)</files>
  <action>
    Create comprehensive city database:
    - Add new regions: Hong Kong, Berlin, Toronto, Mexico City, Lagos, Johannesburg, Istanbul, Bangkok, Jakarta, Seoul, Auckland, Mexico City, Denver, Chicago, Miami, Toronto, Vancouver, Amsterdam, Stockholm, Warsaw, Athens, Tel Aviv, Riyadh, Karachi, Manila, Ho Chi Minh
    - Categorize: Americas, Europe, Middle East, Asia, Africa, Oceania
    - Each city: Name, Timezone, Country, Category, Coordinates (for future map markers)
    - Update -list flag to show categorized cities
    - Support city aliases (NYC = New York, LAX = Los Angeles)
  </action>
  <verify>go run main.go -list shows 30+ cities in categories</verify>
  <done>30+ cities available for selection</done>
</task>

<task type="auto">
  <name>Implement Meeting Time Planner</name>
  <files>meeting.go (new), main.go</files>
  <action>
    Create Meeting Time Planner:
    - Press 'M' to enter meeting planner mode
    - Select 2+ cities to compare (using arrow keys + space to select)
    - Define business hours (default 9AM-5PM, configurable)
    - Display timeline showing:
      - Green highlight: all selected cities in business hours
      - Yellow highlight: some cities in business hours
      - Red: no cities in business hours
    - Show "best meeting times" - slots where ALL cities are within business hours
    - Show "acceptable times" - slots where MOST cities are available
    - Display in 30-minute increments across 24 hours
  </action>
  <verify>Press M enters meeting planner, can select cities, shows optimal meeting times</verify>
  <done>Users can find optimal meeting times across timezones</done>
</task>

<task type="auto">
  <name>Add preset city groups</name>
  <files>preset.go (new), main.go</files>
  <action>
    Create preset city groups:
    - business: Major financial hubs (NYC, London, Tokyo, Singapore, Hong Kong, Dubai)
    - family: Cities with personal significance (user-configurable in future)
    - americas: North & South American cities
    - europe: European cities
    - asia: Asian cities
    - all: Default all available cities
    - Support: -preset business shows only business hubs
  </action>
  <verify>go run main.go -preset business shows only business cities</verify>
  <done>Users can quickly switch between city groups</done>
</task>

</tasks>

<verification>
- go run main.go -preset business shows business hubs
- go run main.go -preset asia shows Asian cities
- Meeting planner (M key) allows multi-city comparison
- Optimal meeting times are highlighted and clearly shown
</verification>

<success_criteria>
Users can access 30+ cities and find optimal meeting times across multiple timezones.
</success_criteria>

<output>
After completion, create .planning/phases/04-expansion/04-expansion-SUMMARY.md
</output>
