package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type OverlayState int

const (
	OverlayNone OverlayState = iota
	OverlayMenu
	OverlayFeature
)

type MenuItem struct {
	Label string
	Mode  Mode
	Icon  string
}

type OverlayManager struct {
	state         OverlayState
	pages         *tview.Pages
	selectedIndex int
	menuItems     []MenuItem
	activeFeature Mode
	app           *tview.Application
	mm            *modeManager
}

func NewOverlayManager(app *tview.Application, pages *tview.Pages, mm *modeManager) *OverlayManager {
	return &OverlayManager{
		state: OverlayNone,
		pages: pages,
		menuItems: []MenuItem{
			{Label: "Clocks", Mode: ModeNavigation, Icon: "🕒"},
			{Label: "Converter", Mode: ModeConverter, Icon: "🔄"},
			{Label: "Stopwatch", Mode: ModeStopwatch, Icon: "⏱️"},
			{Label: "Timer", Mode: ModeTimer, Icon: "⏲️"},
			{Label: "Alarm", Mode: ModeAlarm, Icon: "🔔"},
			{Label: "Meeting Planner", Mode: ModeMeeting, Icon: "🗓️"},
		},
		app: app,
		mm:  mm,
	}
}

func (om *OverlayManager) ShowMenu() {
	om.state = OverlayMenu
	om.selectedIndex = 0
	om.renderMenu()
}

func (om *OverlayManager) renderMenu() {
	menuContent := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(false)

	menuContent.SetBorder(true).
		SetTitle("[ Menu ]").
		SetTitleAlign(tview.AlignCenter).
		SetBorderPadding(1, 1, 2, 2)

	var sb strings.Builder
	for i, item := range om.menuItems {
		if i == om.selectedIndex {
			sb.WriteString(fmt.Sprintf("[black:white] %s %-18s [-:-]\n", item.Icon, item.Label))
		} else {
			sb.WriteString(fmt.Sprintf(" %s %-18s \n", item.Icon, item.Label))
		}
	}
	menuContent.SetText(sb.String())

	// Centered layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(menuContent, 32, 0, true).
			AddItem(nil, 0, 1, false),
			10, 0, true).
		AddItem(nil, 0, 1, false)

	om.pages.AddPage("menu", flex, true, true)
	om.app.SetFocus(menuContent)
}

func (om *OverlayManager) ShowFeature(mode Mode) {
	om.state = OverlayFeature
	om.activeFeature = mode
	om.mm.SwitchTo(mode)
	om.pages.RemovePage("menu")
	om.renderFeature()
}

func (om *OverlayManager) renderFeature() {
	featureContent := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetWrap(false)

	title := modeNames[om.activeFeature]
	featureContent.SetBorder(true).
		SetTitle(fmt.Sprintf("[ %s ]", title)).
		SetTitleAlign(tview.AlignCenter).
		SetBorderPadding(1, 1, 2, 2)

	// Get feature content
	var content string
	if om.activeFeature == ModeNavigation {
		// Use Navigation mode for "Clocks" list
		content = RenderClockList(append(leftRegions, rightRegions...))
	} else {
		content = om.mm.Render()
	}

	// Add ASCII clock for relevant features
	if om.activeFeature == ModeTimer || om.activeFeature == ModeStopwatch || om.activeFeature == ModeAlarm {
		clock := RenderASCIIClock(time.Now())
		// Layout clock and content side-by-side
		lines := strings.Split(content, "\n")
		clockLines := strings.Split(clock, "\n")

		maxLines := len(lines)
		if len(clockLines) > maxLines {
			maxLines = len(clockLines)
		}

		var combined strings.Builder
		for i := 0; i < maxLines; i++ {
			cLine := ""
			if i < len(clockLines) {
				cLine = clockLines[i]
			} else {
				cLine = strings.Repeat(" ", 27)
			}

			fLine := ""
			if i < len(lines) {
				fLine = lines[i]
			}

			combined.WriteString(cLine)
			combined.WriteString("   ")
			combined.WriteString(fLine)
			combined.WriteString("\n")
		}
		content = combined.String()
	}

	featureContent.SetText(content)

	// Centered layout - size depends on content
	width := 64
	height := 18
	if om.activeFeature == ModeNavigation {
		height = 22
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(featureContent, width, 0, true).
			AddItem(nil, 0, 1, false),
			height, 0, true).
		AddItem(nil, 0, 1, false)

	om.pages.AddPage("feature", flex, true, true)
	om.app.SetFocus(featureContent)
}

func (om *OverlayManager) CloseOverlay() {
	if om.state == OverlayFeature {
		om.pages.RemovePage("feature")
		om.mm.SwitchTo(ModeNormal)
		om.ShowMenu()
	} else if om.state == OverlayMenu {
		om.pages.RemovePage("menu")
		om.state = OverlayNone
	}
}

func (om *OverlayManager) HandleInput(event *tcell.EventKey) bool {
	if om.state == OverlayMenu {
		switch event.Key() {
		case tcell.KeyUp:
			om.selectedIndex--
			if om.selectedIndex < 0 {
				om.selectedIndex = len(om.menuItems) - 1
			}
			om.renderMenu()
			return true
		case tcell.KeyDown:
			om.selectedIndex++
			if om.selectedIndex >= len(om.menuItems) {
				om.selectedIndex = 0
			}
			om.renderMenu()
			return true
		case tcell.KeyEnter:
			om.ShowFeature(om.menuItems[om.selectedIndex].Mode)
			return true
		case tcell.KeyEscape:
			om.CloseOverlay()
			return true
		}
		if event.Rune() == 'm' || event.Rune() == 'M' || event.Rune() == ' ' {
			om.CloseOverlay()
			return true
		}
	} else if om.state == OverlayFeature {
		if event.Key() == tcell.KeyEscape {
			om.CloseOverlay()
			return true
		}
		// Delegate to mode handlers
		if event.Key() == tcell.KeyRune {
			return om.mm.HandleKey(event.Rune())
		}
		// Handle special keys like Enter, Backspace for converter/meeting
		return om.mm.HandleSpecialKeyEvent(event.Key())
	}
	return false
}

func RenderClockList(regions []Region) string {
	var sb strings.Builder

	// Group by category if possible, or just Americas/Europe vs Others
	sb.WriteString("[yellow::b]Americas & Europe[-]\n")
	for _, r := range regions {
		city := GetCityByName(r.Name)
		if city == nil || (city.Category != "Americas" && city.Category != "Europe") {
			continue
		}
		loc, _ := time.LoadLocation(r.Timezone)
		now := time.Now().In(loc)
		sb.WriteString(fmt.Sprintf("  [white]%s  %-14s [green]%8s  [silver]%s[-]\n",
			city.Abbreviation(), r.Name, now.Format("3:04 PM"), now.Format("Mon Jan 2")))
	}

	sb.WriteString("\n[yellow::b]Asia, Africa & Oceania[-]\n")
	for _, r := range regions {
		city := GetCityByName(r.Name)
		if city == nil || (city.Category == "Americas" || city.Category == "Europe") {
			continue
		}
		loc, _ := time.LoadLocation(r.Timezone)
		now := time.Now().In(loc)
		sb.WriteString(fmt.Sprintf("  [white]%s  %-14s [green]%8s  [silver]%s[-]\n",
			city.Abbreviation(), r.Name, now.Format("3:04 PM"), now.Format("Mon Jan 2")))
	}

	return sb.String()
}
