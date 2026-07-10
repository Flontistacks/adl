package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const bannerArt = ` █████╗ ██████╗ ██╗
██╔══██╗██╔══██╗██║
███████║██║  ██║██║
██╔══██║██║  ██║██║
██║  ██║██████╔╝██║
╚═╝  ╚═╝╚═════╝ ╚═╝`

func renderHeader() string {
	art := bannerStyle.Render(bannerArt)
	link := linkStyle.Render(githubURL)
	tagline := taglineStyle.Render("Terminal download manager powered by aria2c.")

	artLines := strings.Split(art, "\n")
	right := []string{link, "", tagline}
	for len(right) < len(artLines) {
		right = append(right, "")
	}

	var b strings.Builder
	for i, line := range artLines {
		b.WriteString(line)
		if i < len(right) && right[i] != "" {
			b.WriteString("  ")
			b.WriteString(right[i])
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderFooter(keys ...string) string {
	return footerStyle.Render(strings.Join(keys, " | "))
}

func padRight(s string, width int) string {
	if lipgloss.Width(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-lipgloss.Width(s))
}
