---
phase: 02-utilities
plan: 01
type: execute
wave: 1
depends_on: []
files_modified: [main.go, mode.go (new)]
autonomous: true
requirements: [F3, F4, F5]
---

<objective>
Implement time converter, stopwatch/timer, and alarm functionality.
</objective>

<context>
@main.go
@.planning/phases/01-foundation/01-foundation-SUMMARY.md (after phase 1)
</context>

<tasks>

<task type="auto">
  <name>Create keyboard mode system</name>
  <files>mode.go (new)</files>
  <action>
    Create mode.go to handle different keyboard interaction modes:
    - Mode enum: Normal, Navigation, Converter, Stopwatch, Timer, Alarm
    - Mode indicator in header showing current mode
    - Mode-specific key handlers registered per mode
    - Escape key returns to Normal mode from any mode
  </action>
  <verify>Header shows current mode, Escape returns to Normal mode from any mode</verify>
  <done>Framework for mode-specific keyboard handling exists</done>
</task>

<task type="auto">
  <name>Implement Time Converter mode</name>
  <files>converter.go (new), main.go</files>
  <action>
    Create Time Converter mode:
    - Press 'C' to enter converter mode
    - Display input prompt: "Enter time (HH:MM) and select source timezone"
    - Show all configured cities with converted times in real-time
    - Use arrow keys to select different source timezones
    - Display format: "Tokyo: 14:30 → London: 05:30, New York: 00:30"
    - Press Escape to exit converter mode
  </action>
  <verify>Press C enters converter mode, times convert correctly across all timezones</verify>
  <done>Users can convert any time to all displayed timezones</done>
</task>

<task type="auto">
  <name>Implement Stopwatch mode</name>
  <files>stopwatch.go (new), main.go</files>
  <action>
    Create Stopwatch mode:
    - Press 'S' to enter stopwatch mode
    - Display large stopwatch: "00:00:00.00"
    - Controls: Space=Start/Pause, R=Reset
    - Lap functionality: L key records lap times
    - Show lap times in a scrollable list below stopwatch
    - Press Escape to exit (stopwatch continues in background)
  </action>
  <verify>Press S shows stopwatch, Space starts/pauses, R resets, L records laps</verify>
  <done>Functional stopwatch with start/pause/reset/lap features</done>
</task>

<task type="auto">
  <name>Implement Countdown Timer mode</name>
  <files>timer.go (new), main.go</files>
  <action>
    Create Timer (countdown) mode:
    - Press 'T' to enter timer mode
    - Input: Enter duration in HH:MM:SS format
    - Display countdown with large digits
    - Controls: Space=Start/Pause, R=Reset
    - When timer reaches zero: flash display, play terminal bell
    - Press Escape to exit
  </action>
  <verify>Press T enters timer mode, user can set duration, countdown works correctly</verify>
  <done>Functional countdown timer with alarm notification</done>
</task>

<task type="auto">
  <name>Implement Alarm system</name>
  <files>alarm.go (new), main.go</files>
  <action>
    Create Alarm functionality:
    - Press 'A' to add new alarm
    - Input: Select city/timezone, set time (HH:MM), set repeat (once/daily/weekday)
    - Store alarms in config file
    - Display active alarms in a panel
    - When alarm triggers: flash alarm panel, play terminal bell, show "ALARM!" message
    - Press any key to dismiss alarm
    - Alarms persist across sessions
  </action>
  <verify>Press A allows setting alarm, alarm triggers at correct time with notification</verify>
  <done>Functional alarm system with persistent storage</done>
</task>

</tasks>

<verification>
- C key enters time converter mode
- S key enters stopwatch mode  
- T key enters timer mode
- A key allows setting alarms
- All modes accessible from Normal mode via single keypress
- Escape returns to Normal mode from any mode
</verification>

<success_criteria>
Users can convert times, use stopwatch, set countdown timers, and create persistent alarms.
</success_criteria>

<output>
After completion, create .planning/phases/02-utilities/02-utilities-SUMMARY.md
</output>
