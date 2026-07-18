package tui

import (
	"strings"
	"testing"
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
