# 🌍 localize

**A terminal-based world clock dashboard for the modern era**

![License](https://img.shields.io/github/license/hamshad/localize?color=blue)
![Go Version](https://img.shields.io/github/go-mod/go-version/hamshad/localize)
![Last Commit](https://img.shields.io/github/last-commit/hamshad/localize)

---

**localize** is a beautiful, functional terminal application that brings the world's timezones to your fingertips. Built with Go and TUI frameworks, it features a stunning Braille-rendered world map, real-time clocks, and powerful productivity tools for global collaboration.

<div align="center">
  
![Screenshot](https://placehold.co/800x400/1a1a2e/FFF?text=localize+-+World+Time+Dashboard+Preview)
  
</div>

---

## ✨ Features

### 🌐 World Time at a Glance
- **137+ cities** across 6 regions (Americas, Europe, Middle East, Asia, Africa, Oceania)
- **Braille world map** with color-coded city markers
- **Real-time updates** every second
- **Day/night overlay** showing global daylight distribution
- **Relative time offsets** comparing cities to each other

### 🛠 Productivity Tools
- **Meeting Planner** — Find overlapping business hours across timezones
- **Time Converter** — Convert times between any two cities
- **Stopwatch & Timer** — Track time with precision
- **Alarm System** — Set timezone-aware alarms with notifications

### ⌨ Keyboard-First Interface
```
Navigation:  ↑↓ arrow keys, Tab to switch panels, Enter for details
Modes:       c=Converter, s=Stopwatch, t=Timer, a=Alarm, m=Meeting
Toggle:      d=Day/Night overlay
Exit:        Q or Esc
```

### 🔧 Customization
- **Preset city groups** — business, family, americas, europe, asia, africa, oceania
- **Custom city selection** via CLI flags
- **Config persistence** to `~/.localize/config.json`
- **Color-coded regions** for visual clarity

---

## 🚀 Getting Started

### Prerequisites
- Go 1.25.6+
- Terminal with Unicode (Braille) support

### Installation

#### From Source
```bash
git clone https://github.com/hamshad/localize.git
cd localize
go build -o localize .
./localize
```

#### Using Go Install
```bash
go install github.com/hamshad/localize@latest
```

### Usage

#### Default Mode (Show All Cities)
```bash
./localize
```

#### City Presets
```bash
./localize -preset business     # Business hubs
./localize -preset family       # Family-friendly timezones
./localize -preset americas     # Americas only
./localize -preset europe       # Europe only
./localize -preset asia         # Asia only
```

#### Custom City Selection
```bash
./localize -cities "Tokyo,London,New York"
./localize -preset business -cities "Dubai,Singapore"
```

#### List Available Cities
```bash
./localize -list
```

---

## 📁 Project Structure

```
.
├── main.go           # Core application logic & UI
├── cities.go         # City database (137+ locations)
├── config.go         # Configuration persistence
├── navigation.go     # Keyboard navigation
├── mode.go           # Mode system (converter, timer, etc.)
├── converter.go      # Time converter implementation
├── stopwatch.go      # Stopwatch functionality
├── timer.go          # Countdown timer
├── alarm.go          # Alarm system
├── meeting.go        # Meeting planner
├── daynight.go       # Day/night overlay logic
├── go.mod            # Go module definition
└── LICENCE           # MIT Licence
```

---

## 🔌 Key Press Reference

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate cities |
| `Tab` | Switch between left/right panels |
| `Enter` | Toggle city details panel |
| `Esc` | Exit navigation / quit app |
| `c` | Switch to Time Converter mode |
| `s` | Switch to Stopwatch mode |
| `t` | Switch to Timer mode |
| `a` | Switch to Alarm mode |
| `m` | Switch to Meeting Planner mode |
| `d` | Toggle Day/Night overlay |
| `Q` / `q` | Quit application |

---

## 📦 Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/rivo/tview` | Terminal UI framework |
| `github.com/gdamore/tcell/v2` | Terminal cell management |
| Go standard library | Time handling, file I/O, etc. |

---

## 🌟 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Questions?** Open an [issue](https://github.com/hamshad/localize/issues) to discuss before implementing major changes.

---

## 📄 Licence

This project is licensed under the **MIT Licence** — see the [LICENCE](/LICENCE) file for details.

```
MIT Licence

Copyright (c) 2026 Hamshad

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 Acknowledgements

- Inspired by terminal ricing culture and CLI-first tools
- Built with [tview](https://github.com/rivo/tview) — a beautiful TUI framework
- City data sourced from IANA timezone database

---

<div align="center">
  
**Happy time-tracking across timezones! 🕒✨**

[GitHub](https://github.com/hamshad/localize) • [Issues](https://github.com/hamshad/localize/issues) • [Discussions](https://github.com/hamshad/localize/discussions)

</div>
