package tui

import (
	"fmt"
	"strings"
)

type menuEntry struct {
	title  string
	desc   string
	screen Screen
	quit   bool
}

var mainMenuEntries = []menuEntry{
	{"New Download", "Add HTTP, magnet, or torrent", ScreenNewDownload, false},
	{"Active Downloads", "Progress and controls", ScreenActive, false},
	{"Settings", "Default folder and aria2c path", ScreenSettings, false},
	{"Help", "Keybindings and commands", ScreenHelp, false},
	{"Quit", "Stop daemon and exit", ScreenMenu, true},
}

func (m Model) viewMenu() string {
	const titleWidth = 22
	var b strings.Builder
	b.WriteString("\n")

	for i, item := range mainMenuEntries {
		num := menuNumStyle.Render(fmt.Sprintf("%d.", i+1))
		if i == m.menuIndex {
			cursor := cursorStyle.Render(">")
			title := menuActive.Render(padRight(item.title, titleWidth))
			b.WriteString(fmt.Sprintf("%s %s %s%s\n", cursor, num, title, menuDescStyle.Render(item.desc)))
		} else {
			title := menuDescStyle.Render(padRight(item.title, titleWidth))
			b.WriteString(fmt.Sprintf("  %s %s%s\n", num, title, menuDescStyle.Render(item.desc)))
		}
	}

	b.WriteString("\n")
	b.WriteString(renderFooter("↑↓", "Enter", "? Help", "Q Quit"))
	return b.String()
}
