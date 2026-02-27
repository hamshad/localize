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

// cityMarker defines a city's position on the braille map and its display label.
type cityMarker struct {
	label string
	row   int // row in the braille output (0-indexed)
	col   int // col in the braille output (0-indexed, in characters)
	color string
}

var cityMarkers = []cityMarker{
	{label: "HNL", row: 6, col: 2, color: "turquoise"},
	{label: "ANC", row: 3, col: 7, color: "lightcyan"},
	{label: "LA", row: 4, col: 12, color: "darkcyan"},
	{label: "NYC", row: 4, col: 17, color: "dodgerblue"},
	{label: "GRU", row: 9, col: 18, color: "limegreen"},
	{label: "LON", row: 3, col: 31, color: "green"},
	{label: "PAR", row: 3, col: 35, color: "darkgreen"},
	{label: "MOS", row: 3, col: 39, color: "red"},
	{label: "CAI", row: 5, col: 33, color: "sandybrown"},
	{label: "NBO", row: 7, col: 35, color: "coral"},
	{label: "DXB", row: 5, col: 38, color: "gold"},
	{label: "BOM", row: 6, col: 41, color: "orange"},
	{label: "SIN", row: 7, col: 46, color: "darkmagenta"},
	{label: "PVG", row: 4, col: 49, color: "orangered"},
	{label: "TYO", row: 4, col: 52, color: "deeppink"},
	{label: "SYD", row: 11, col: 53, color: "yellow"},
}

// getBrailleWorldMap returns the world map rendered in Braille Unicode characters.
// Each braille character encodes a 2-wide x 4-tall pixel grid.
// The bitmap below is 120 cols x 56 rows (=> 60 braille chars wide x 14 braille chars tall).
func getBrailleWorldMap() string {
	// The bitmap: 1 = land, 0 = water
	// 120 columns x 56 rows — a simplified but recognizable equirectangular world map
	// Rows go top-to-bottom (north pole -> south pole)
	bitmap := []string{
		//0         1         2         3         4         5         6         7         8         9         10        11
		//0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 0
		"000000000000000000000000000000000000000000000000000000000000000000001100000000000000000000001111111111000000000000000000000000000000", // 1
		"000000000000001100000000000000000000000000000000000000000000000001111111100000000000000000011111111111100000000000000000000000000000", // 2
		"000000000000111111100000000000000000000000000000000000000000001111111111110000000000000001111111111111111000000000000000000000000000", // 3
		"000000000011111111111100000000000000000000000000000000000000011111111111111100000000000011111111111111111100000000000000000000000000", // 4
		"000000000111111111111111000000000000000000000000000000000000111111111111111100000001101111111111111111111110000000000000000000000000", // 5
		"000000001111111111111111100000000000000000000000000000000001111111111111111111111111111111111111111111111111000000000000000000000000", // 6
		"000000001111111111111111110000000000000000000000000000000011111111111111111111111111111111111111111111111111100000000000000000000000", // 7
		"000000011111111111111111111000000000000000000000000000000111111111111111111111111111111111111111111111111111110000000000000000000000", // 8
		"000000011111111111111111111000000000000000000000000000011111111111111111111111111111111111111111111111111111111000000000000000000000", // 9
		"000000111111111111111111111100000000000000000000000000111111111111111111111111111111111111111111111111111111111110000000000000000000", // 10
		"000000111111111111111111111110000000000000000000000001111111111111111111111111111111111111111111111111111111111111000000000000000000", // 11
		"000000111111111111111111111111000000000000000000000001111111111111111111111111111111111111111111111111111111111111100000000000000000", // 12
		"000001111111111111111111111111000000000000000000000001111111111111111111111111111111111111111111111111111111111111110000000000000000", // 13
		"000001111111111111111111111111100000000000000000000001111111111111111111111111111111111111111111111111111111111111111000000000000000", // 14
		"000001111111111111111111111111100000000000000000000001111111111111111111111111111111111111111111111111111111111111111100000000000000", // 15
		"000001111111111111111111111111100000000000000000000001111111111111111111111111111111111111111111111111111111111111111110000000000000", // 16
		"000001111111111111111111111111110000000000000000000000111111111111111111111111111111111111111111111111111111111111111110000000000000", // 17
		"000000111111111111111111111111110000000000000000000000011111111111111111111111111111111111111111111111111111111111111110000000000000", // 18
		"000000011111111111111111111111110000000000000000000000001111111111111111111111111111111111111111111111111111111111111100000000000000", // 19
		"000000001111111111111111111111110000000000000000000000000011111111111101111111111111111111111111111111111111111111111000000000000000", // 20
		"000000000111111111111111111111000000000000000000000000000001111111111000011111111111111111111111111111111111111111100000000000000000", // 21
		"000000000011111111111111111111000000000000000000000000000000011111110000001111111111111111111111111111111011111110000000000000000000", // 22
		"000000000001111111111111111110000000000000000000000000000000001111000000011111111111111111111111111111110001111100000000000000000000", // 23
		"000000000000111111110011111100000000000000000000000000000000000000000000111111111111111111111111111111000000110000000000000000000000", // 24
		"000000000000111111100001111100000000000000000000000000000000000000000001111111111111111011111111111100000000000000000000000000000000", // 25
		"000000000000011111000001111100000000000000000000000000000000000000000011111111111111100001111111111000000000000000000000000000000000", // 26
		"000000000000001111100001111110000000000000000000000000000000000000000011111111111110000000111111110000000000000000000000000000000000", // 27
		"000000000000000111100001111110000000000000000000000000000000000000000001111111110000000000011111100000000000000000000000000000000000", // 28
		"000000000000000111100000111110000000000000000000000000000000000000000000111111100000000000011111000000000000000000000000000000000000", // 29
		"000000000000000011110000111111000000000000000000000000000000000000000000001110000000000000011110000000000000000000000000000000000000", // 30
		"000000000000000001110000011111000000000000000000000000000000000000000000000000000000000000001100000000000000000000000000000000000000", // 31
		"000000000000000001110000011111100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 32
		"000000000000000000111000001111100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 33
		"000000000000000000111000000111110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 34
		"000000000000000000011100000011111000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 35
		"000000000000000000001100000011111100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 36
		"000000000000000000001110000001111100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 37
		"000000000000000000000110000001111110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 38
		"000000000000000000000111000000111111000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 39
		"000000000000000000000011000000011111000000000000000000000000000000000000000000000000000000000000000011111111100000000000000000000000", // 40
		"000000000000000000000011100000011111100000000000000000000000000000000000000000000000000000000000000111111111110000000000000000000000", // 41
		"000000000000000000000001100000001111100000000000000000000000000000000000000000000000000000000000001111111111111000000000000000000000", // 42
		"000000000000000000000001100000000111100000000000000000000000000000000000000000000000000000000000001111111111111100000000000000000000", // 43
		"000000000000000000000000110000000011110000000000000000000000000000000000000000000000000000000000001111111111111100000000000000000000", // 44
		"000000000000000000000000010000000001110000000000000000000000000000000000000000000000000000000000000111111111111000000000000000000000", // 45
		"000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000011111111110000000000000000000000", // 46
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001111111100000000000000000000000", // 47
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011100000000000000000000000000", // 48
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 49
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 50
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", // 51
		"000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111111111111111111111000000000000000000000", // 52
		"000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111000000000000000000", // 53
		"000000000000000000001111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111000000000000000", // 54
		"000000000000000001111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111110000000000000", // 55
	}

	rows := len(bitmap)
	cols := len(bitmap[0])

	// Braille encoding: each braille char is 2 cols x 4 rows of dots
	// Dot positions and their bit values:
	// (0,0)=0x01  (1,0)=0x08
	// (0,1)=0x02  (1,1)=0x10
	// (0,2)=0x04  (1,2)=0x20
	// (0,3)=0x40  (1,3)=0x80
	brailleBase := rune(0x2800)

	brailleRows := rows / 4
	brailleCols := cols / 2

	var sb strings.Builder
	for by := 0; by < brailleRows; by++ {
		for bx := 0; bx < brailleCols; bx++ {
			var code rune
			// Map the 2x4 block to braille dots
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					py := by*4 + dy
					px := bx*2 + dx
					if py < rows && px < cols && bitmap[py][px] == '1' {
						// Braille dot encoding
						switch {
						case dx == 0 && dy == 0:
							code |= 0x01
						case dx == 0 && dy == 1:
							code |= 0x02
						case dx == 0 && dy == 2:
							code |= 0x04
						case dx == 1 && dy == 0:
							code |= 0x08
						case dx == 1 && dy == 1:
							code |= 0x10
						case dx == 1 && dy == 2:
							code |= 0x20
						case dx == 0 && dy == 3:
							code |= 0x40
						case dx == 1 && dy == 3:
							code |= 0x80
						}
					}
				}
			}
			sb.WriteRune(brailleBase + code)
		}
		sb.WriteRune('\n')
	}

	return sb.String()
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
		AddItem(mapView, 18, 0, false).
		AddItem(clockRow, 0, 1, false).
		AddItem(footer, 1, 0, false)

	// ── UPDATE FUNCTION ──
	updateUI := func() {
		now := time.Now()
		utcNow := now.UTC()

		// Header
		headerText := fmt.Sprintf(
			"\n[::b][dodgerblue]   LOCALIZE [white]- World Time Dashboard[::-]   |   [yellow]UTC: %s[white]   |   [aqua]%s",
			utcNow.Format("15:04:05"),
			utcNow.Format("Monday, 02 January 2006"),
		)
		header.SetText(headerText)

		// Braille Map with overlaid city markers
		brailleMap := getBrailleWorldMap()
		coloredMap := colorizeBrailleMap(brailleMap)
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

// colorizeBrailleMap takes the raw braille map string and adds color tags +
// overlays city markers at their positions.
func colorizeBrailleMap(brailleMap string) string {
	lines := strings.Split(strings.TrimRight(brailleMap, "\n"), "\n")

	// Overlay city markers onto map lines
	// We replace braille chars at specific positions with colored city labels
	for _, m := range cityMarkers {
		if m.row < len(lines) {
			runes := []rune(lines[m.row])
			label := m.label
			// Ensure the label fits
			if m.col+len(label) <= len(runes) {
				// Build the replacement: color tag + label + reset + rest
				before := string(runes[:m.col])
				after := string(runes[m.col+len(label):])
				lines[m.row] = before + fmt.Sprintf("[%s::b]%s[-::-]", m.color, label) + after
			}
		}
	}

	var sb strings.Builder
	for _, line := range lines {
		// Wrap each line in green color for the braille land dots
		// but preserve any already-injected color tags for markers
		sb.WriteString("[green]")
		sb.WriteString(line)
		sb.WriteString("[-]\n")
	}

	return sb.String()
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
