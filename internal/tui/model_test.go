package tui

import (
	"strings"
	"testing"

	"github.com/Flontistacks/adl/internal/aria2"
	"github.com/Flontistacks/adl/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TestRenderTerminalScreenUsesMoleStyleRedraw(t *testing.T) {
	got := renderTerminalScreen("content")

	if !strings.HasPrefix(got, "\x1b[H") {
		t.Fatal("screen must move the cursor home before rendering")
	}
	if strings.Contains(got, "\x1b[2J") {
		t.Fatal("screen must not erase the entire display into scrollback")
	}
	if !strings.HasSuffix(got, "\x1b[J") {
		t.Fatal("screen must erase only the area below the rendered content")
	}
}

func TestCtrlCQuitsProgram(t *testing.T) {
	m := NewModel(config.Config{}, ViewMenu)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("ctrl+c must return a command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("ctrl+c must return tea.Quit")
	}
}

func TestQuestionMarkCanBeTypedInDownloadURL(t *testing.T) {
	m := NewModel(config.Config{}, ViewDownload)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	got := updated.(Model)

	if got.screen != ScreenNewDownload {
		t.Fatal("question mark must not open help while editing a URL")
	}
	if got.urlInput.Value() != "?" {
		t.Fatalf("URL input = %q, want ?", got.urlInput.Value())
	}
}

func TestDownloadsUpdateClampsSelection(t *testing.T) {
	m := NewModel(config.Config{}, ViewList)
	m.selectedDL = 4

	updated, cmd := m.Update(downloadsMsg{items: []aria2.Status{{GID: "only"}}})
	got := updated.(Model)

	if cmd != nil {
		t.Fatal("a downloads response must not create another polling timer")
	}
	if got.selectedDL != 0 {
		t.Fatalf("selected index = %d, want 0", got.selectedDL)
	}
}

func TestDownloadsUpdateRefreshesOpenDetail(t *testing.T) {
	m := NewModel(config.Config{}, ViewList)
	m.screen = ScreenDetail
	m.detail = aria2.Status{GID: "gid-1", Completed: 1}

	updated, _ := m.Update(downloadsMsg{items: []aria2.Status{{
		GID:       "gid-1",
		Completed: 99,
	}}})
	got := updated.(Model)

	if got.detail.Completed != 99 {
		t.Fatalf("detail progress = %d, want 99", got.detail.Completed)
	}
}

func TestSanitizeTerminalTextRemovesControlCharacters(t *testing.T) {
	input := "safe\x1b]52;c;secret\a\nnext"
	got := sanitizeTerminalText(input)

	if strings.ContainsAny(got, "\x1b\a\n\r") {
		t.Fatalf("sanitized text still contains terminal control characters: %q", got)
	}
	if !strings.Contains(got, "safe") || !strings.Contains(got, "next") {
		t.Fatalf("sanitized text lost printable content: %q", got)
	}
}
