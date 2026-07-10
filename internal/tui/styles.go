package tui

import "github.com/charmbracelet/lipgloss"

const githubURL = "https://github.com/Flontistacks/adl"

var (
	colorGreen  = lipgloss.Color("82")
	colorBlue   = lipgloss.Color("39")
	colorCyan   = lipgloss.Color("87")
	colorMuted  = lipgloss.Color("245")
	colorWhite  = lipgloss.Color("255")
	colorRed    = lipgloss.Color("196")
	colorYellow = lipgloss.Color("220")

	bannerStyle   = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	linkStyle     = lipgloss.NewStyle().Foreground(colorBlue).Underline(true)
	taglineStyle  = lipgloss.NewStyle().Foreground(colorGreen)
	sectionStyle  = lipgloss.NewStyle().Foreground(colorWhite).Bold(true)
	cursorStyle   = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	menuNumStyle  = lipgloss.NewStyle().Foreground(colorMuted)
	menuActive    = lipgloss.NewStyle().Foreground(colorWhite).Bold(true)
	menuDescStyle = lipgloss.NewStyle().Foreground(colorMuted)
	footerStyle   = lipgloss.NewStyle().Foreground(colorMuted)
	helpStyle     = lipgloss.NewStyle().Foreground(colorMuted)
	errStyle      = lipgloss.NewStyle().Foreground(colorRed)
	okStyle       = lipgloss.NewStyle().Foreground(colorGreen)
	boxStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(colorMuted).Padding(0, 1)
)
