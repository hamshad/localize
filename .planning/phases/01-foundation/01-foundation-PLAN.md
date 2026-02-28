---
phase: 01-foundation
plan: 01
type: execute
wave: 1
depends_on: []
files_modified: [main.go]
autonomous: true
requirements: [F1, F9]
---

<objective>
Implement user-configurable cities via CLI flags and config file persistence.
</objective>

<context>
@main.go
</context>

<tasks>

<task type="auto">
  <name>Add CLI flag parsing for city selection</name>
  <files>main.go</files>
  <action>
    Add flag parsing using the standard library's flag package. Create:
    - `-cities` flag: comma-separated list of city names (e.g., "Tokyo,London,New York")
    - `-preset` flag: use predefined city groups ("business", "family", "americas", "europe", "asia")
    - `-list` flag: show all available cities and exit
    Default behavior remains unchanged (show all 16 cities).
  </action>
  <verify>go run main.go -list shows all available cities; go run main.go -cities "Tokyo,London" shows only those two</verify>
  <done>User can specify which cities to display via CLI flags</done>
</task>

<task type="auto">
  <name>Create config file structure and persistence</name>
  <files>config.go (new file)</files>
  <action>
    Create config.go with:
    - Config struct: Cities []string, Preset string, Layout preferences
    - Default config path: ~/.localize/config.json
    - LoadConfig() function: reads from file, falls back to defaults
    - SaveConfig() function: writes current config to file
    - GetAvailableCities() function: returns list of all supported cities with their timezones
  </action>
  <verify>Config file created at ~/.localize/config.json with default settings</verify>
  <done>User preferences persist between sessions</done>
</task>

<task type="auto">
  <name>Integrate config with main app</name>
  <files>main.go</files>
  <action>
    Modify main.go to:
    - Load config on startup (config.LoadConfig())
    - Use config.Cities if provided, otherwise CLI flags, otherwise defaults
    - Auto-save config if user makes changes (future feature)
    - Create ~/.localize directory if it doesn't exist
  </action>
  <verify>App reads from config file and displays configured cities</verify>
  <done>Config and CLI flags work together seamlessly</done>
</task>

</tasks>

<verification>
- CLI flag -list displays all 16+ cities
- CLI flag -cities "Tokyo,London" shows only those cities
- Config file ~/.localize/config.json is created on first run
- App respects both CLI flags and config file preferences
</verification>

<success_criteria>
Users can customize displayed cities via CLI or config file, and preferences persist across sessions.
</success_criteria>

<output>
After completion, create .planning/phases/01-foundation/01-foundation-SUMMARY.md
</output>
