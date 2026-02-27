package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Region represents a major world region with timezone info.
type Region struct {
	Name     string
	Timezone string
	Color    tcell.Color
}

var leftRegions = []Region{
	{Name: "Honolulu", Timezone: "Pacific/Honolulu", Color: tcell.ColorTurquoise},
	{Name: "Anchorage", Timezone: "America/Anchorage", Color: tcell.ColorLightCyan},
	{Name: "Los Angeles", Timezone: "America/Los_Angeles", Color: tcell.ColorDarkCyan},
	{Name: "New York", Timezone: "America/New_York", Color: tcell.ColorDodgerBlue},
	{Name: "Sao Paulo", Timezone: "America/Sao_Paulo", Color: tcell.ColorLimeGreen},
	{Name: "London", Timezone: "Europe/London", Color: tcell.ColorGreen},
	{Name: "Paris", Timezone: "Europe/Paris", Color: tcell.ColorDarkGreen},
	{Name: "Moscow", Timezone: "Europe/Moscow", Color: tcell.ColorRed},
}

var rightRegions = []Region{
	{Name: "Cairo", Timezone: "Africa/Cairo", Color: tcell.ColorSandyBrown},
	{Name: "Nairobi", Timezone: "Africa/Nairobi", Color: tcell.ColorCoral},
	{Name: "Dubai", Timezone: "Asia/Dubai", Color: tcell.ColorGold},
	{Name: "Mumbai", Timezone: "Asia/Kolkata", Color: tcell.ColorOrange},
	{Name: "Singapore", Timezone: "Asia/Singapore", Color: tcell.ColorDarkMagenta},
	{Name: "Shanghai", Timezone: "Asia/Shanghai", Color: tcell.ColorOrangeRed},
	{Name: "Tokyo", Timezone: "Asia/Tokyo", Color: tcell.ColorDeepPink},
	{Name: "Sydney", Timezone: "Australia/Sydney", Color: tcell.ColorYellow},
}

func getWorldMap() string {
	lines := []string{
		"                                                                                            ",
		"                 ___.__                                                  _____              ",
		"             ___/       \\___                          ___               /     \\____         ",
		"           _/               \\__            __________/   \\____    _____/          \\__       ",
		"          /    NORTH          \\     __    /                   \\__/                   |      ",
		"         |    AMERICA          |   /  \\__/     EUROPE           \\     ASIA           |     ",
		"         |              *NYC   |  |   *LON  *PAR                 \\               *TYO|     ",
		"          \\   *LA             /   |          *MOS    *DXB         \\    *PVG        __/     ",
		"           \\                 /    |  *CAI          *BOM            |           ___/        ",
		"            \\___     ___    /      \\                    *SIN       |          /            ",
		"                \\   /   \\  /    AFRICA                             |    _____/             ",
		"                 \\_/  *GRU\\         *NBO                           |   /                   ",
		"                  |       |\\                                       |  /                    ",
		"                  |SOUTH  | \\                                     /  |                     ",
		"                  |AMERICA|  \\                              *SYD |   |                     ",
		"                   \\      /   \\                          AUSTRALIA   |                     ",
		"                    \\    /     \\                                /    /                     ",
		"                     \\__/      \\______________________________/  __/                      ",
		"                                                                                          ",
		"                              A  N  T  A  R  C  T  I  C  A                                ",
	}
	return strings.Join(lines, "\n")
}

func main() {
	app := tview.NewApplication()

	// ── MAP VIEW ──
	mapView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)
	mapView.SetBorder(true).
		SetTitle(" [ World Map ] ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorDodgerBlue)

	// ── CLOCK PANELS ──
	leftClockView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false)
	leftClockView.SetBorder(true).
		SetTitle(" Americas & Europe ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorGreen)

	rightClockView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false)
	rightClockView.SetBorder(true).
		SetTitle(" Asia, Africa & Oceania ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorOrange)

	// ── HEADER ──
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)

	// ── FOOTER ──
	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)

	// ── LAYOUT ──
	clockRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftClockView, 0, 1, false).
		AddItem(rightClockView, 0, 1, false)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(mapView, 24, 0, false).
		AddItem(clockRow, 0, 1, false).
		AddItem(footer, 1, 0, false)

	// ── UPDATE FUNCTION ──
	updateUI := func() {
		now := time.Now()
		utcNow := now.UTC()

		// Header with app title and UTC time
		headerText := fmt.Sprintf(
			"\n[::b][dodgerblue]   LOCALIZE [white]- World Time Dashboard[::-]   |   [yellow]UTC: %s[white]   |   [aqua]%s",
			utcNow.Format("15:04:05"),
			utcNow.Format("Monday, 02 January 2006"),
		)
		header.SetText(headerText)

		// Map - colorize it
		rawMap := getWorldMap()
		coloredMap := colorizeMap(rawMap)
		mapView.SetText(coloredMap)

		// Clocks
		leftClockView.SetText(formatClockPanel(leftRegions))
		rightClockView.SetText(formatClockPanel(rightRegions))

		// Footer
		footer.SetText("[darkgray]Press [white::b]Q[darkgray::] or [white::b]Esc[darkgray::] to quit   |   Updates every second   |   [dodgerblue]localize")
	}

	// ── KEY BINDINGS ──
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' || event.Rune() == 'Q' {
				app.Stop()
				return nil
			}
		}
		return event
	})

	// ── TICKER ──
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		app.QueueUpdateDraw(func() { updateUI() })

		for range ticker.C {
			app.QueueUpdateDraw(func() { updateUI() })
		}
	}()

	// ── RUN ──
	if err := app.SetRoot(mainLayout, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

// colorizeMap adds tview color tags to the ASCII world map.
func colorizeMap(rawMap string) string {
	var b strings.Builder
	lines := strings.Split(rawMap, "\n")

	// Region markers and their colors on the map
	markers := map[string]string{
		"*NYC": "[dodgerblue::b]*NYC[-::-]",
		"*LA":  "[darkcyan::b]*LA[-::-]",
		"*LON": "[green::b]*LON[-::-]",
		"*PAR": "[darkgreen::b]*PAR[-::-]",
		"*MOS": "[red::b]*MOS[-::-]",
		"*DXB": "[gold::b]*DXB[-::-]",
		"*BOM": "[orange::b]*BOM[-::-]",
		"*SIN": "[darkmagenta::b]*SIN[-::-]",
		"*TYO": "[deeppink::b]*TYO[-::-]",
		"*SYD": "[yellow::b]*SYD[-::-]",
		"*GRU": "[limegreen::b]*GRU[-::-]",
		"*CAI": "[sandybrown::b]*CAI[-::-]",
		"*NBO": "[coral::b]*NBO[-::-]",
		"*PVG": "[orangered::b]*PVG[-::-]",
	}

	// Continent labels
	continentLabels := map[string]string{
		"NORTH":      "[white::b]NORTH[-::-]",
		"AMERICA":    "[white::b]AMERICA[-::-]",
		"SOUTH":      "[white::b]SOUTH[-::-]",
		"EUROPE":     "[white::b]EUROPE[-::-]",
		"ASIA":       "[white::b]ASIA[-::-]",
		"AFRICA":     "[white::b]AFRICA[-::-]",
		"AUSTRALIA":  "[white::b]AUSTRALIA[-::-]",
		"ANTARCTICA": "[white::d]A  N  T  A  R  C  T  I  C  A[-::-]",
	}

	for _, line := range lines {
		coloredLine := "[green]" + tview.Escape(line) + "[-]"

		// Replace markers with colored versions
		for marker, colored := range markers {
			escaped := tview.Escape(marker)
			coloredLine = strings.ReplaceAll(coloredLine, escaped, colored)
		}

		// Replace continent labels
		for label, colored := range continentLabels {
			if label == "ANTARCTICA" {
				escaped := tview.Escape("A  N  T  A  R  C  T  I  C  A")
				coloredLine = strings.ReplaceAll(coloredLine, escaped, colored)
			} else {
				coloredLine = strings.ReplaceAll(coloredLine, label, colored)
			}
		}

		b.WriteString(coloredLine)
		b.WriteString("\n")
	}

	return b.String()
}

// formatClockPanel builds a formatted clock display for a set of regions.
func formatClockPanel(regs []Region) string {
	var b strings.Builder
	b.WriteString("\n")

	for _, r := range regs {
		loc, err := time.LoadLocation(r.Timezone)
		if err != nil {
			b.WriteString(fmt.Sprintf("  [red]%-13s  ERROR[-]\n", r.Name))
			continue
		}

		now := time.Now().In(loc)
		timeStr := now.Format("15:04:05")
		dateStr := now.Format("Mon, 02 Jan 2006")
		offsetStr := now.Format("-07:00")
		dayPhase := getDayPhase(now)
		colorTag := colorToTag(r.Color)

		// Format: City name | Time | Date | UTC offset | Day phase
		b.WriteString(fmt.Sprintf(
			"  [%s::b]%-13s[-::-] [white::b]%s[-::-]  [silver]%s  [darkgray]UTC%s  %s\n",
			colorTag,
			r.Name,
			timeStr,
			dateStr,
			offsetStr,
			dayPhase,
		))
	}

	return b.String()
}

// getDayPhase returns a colored label for the time of day.
func getDayPhase(t time.Time) string {
	hour := t.Hour()
	switch {
	case hour >= 5 && hour < 8:
		return "[#FFA07A]Dawn[-]"
	case hour >= 8 && hour < 12:
		return "[yellow]Morning[-]"
	case hour >= 12 && hour < 17:
		return "[#FFD700]Afternoon[-]"
	case hour >= 17 && hour < 20:
		return "[orange]Evening[-]"
	case hour >= 20 && hour < 22:
		return "[#CD853F]Dusk[-]"
	default:
		return "[#6A5ACD]Night[-]"
	}
}

// colorToTag maps tcell.Color to tview color tag strings.
func colorToTag(c tcell.Color) string {
	colorMap := map[tcell.Color]string{
		tcell.ColorDodgerBlue:  "dodgerblue",
		tcell.ColorDarkCyan:    "darkcyan",
		tcell.ColorGreen:       "green",
		tcell.ColorDarkGreen:   "darkgreen",
		tcell.ColorRed:         "red",
		tcell.ColorGold:        "gold",
		tcell.ColorOrange:      "orange",
		tcell.ColorDarkMagenta: "darkmagenta",
		tcell.ColorDeepPink:    "deeppink",
		tcell.ColorYellow:      "yellow",
		tcell.ColorLimeGreen:   "limegreen",
		tcell.ColorSandyBrown:  "sandybrown",
		tcell.ColorCoral:       "coral",
		tcell.ColorLightCyan:   "lightcyan",
		tcell.ColorTurquoise:   "turquoise",
		tcell.ColorOrangeRed:   "orangered",
	}
	if tag, ok := colorMap[c]; ok {
		return tag
	}
	return "white"
}
