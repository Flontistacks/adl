package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Flontistacks/adl/internal/aria2"
	"github.com/Flontistacks/adl/internal/config"
	"github.com/Flontistacks/adl/internal/download"
)

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenNewDownload
	ScreenActive
	ScreenSettings
	ScreenHelp
	ScreenDetail
)

type StartView string

const (
	ViewMenu     StartView = "menu"
	ViewDownload StartView = "download"
	ViewList     StartView = "list"
	ViewSettings StartView = "settings"
)

type tickMsg struct{}
type downloadsMsg struct {
	items []aria2.Status
	err   error
}

type Model struct {
	cfg      config.Config
	daemon   *aria2.Daemon
	screen   Screen
	width    int
	height   int
	errMsg   string
	status   string

	menuIndex int

	urlInput  textinput.Model
	destInput textinput.Model
	newStep   int // 0=url, 1=dest

	browsing   bool
	browseDir  string
	browseList list.Model

	downloads  []aria2.Status
	selectedDL int
	detail     aria2.Status

	settingsDir   textinput.Model
	settingsAria2 textinput.Model
	settingsFocus int
}

func styleTextInput(ti *textinput.Model) {
	green := lipgloss.NewStyle().Foreground(colorGreen)
	ti.PromptStyle = green
	ti.TextStyle = lipgloss.NewStyle().Foreground(colorWhite)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(colorGreen)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorMuted)
}

func NewModel(cfg config.Config, start StartView) Model {
	url := textinput.New()
	url.Placeholder = "https://... or magnet:?xt=... or /path/file.torrent"
	url.CharLimit = 2048
	url.Width = 60
	url.Prompt = "URL: "
	styleTextInput(&url)

	dest := textinput.New()
	dest.CharLimit = 512
	dest.Width = 60
	dest.Prompt = "Dest: "
	styleTextInput(&dest)

	sDir := textinput.New()
	sDir.Prompt = "Download dir: "
	sDir.Width = 50
	sDir.SetValue(cfg.DownloadDir)
	styleTextInput(&sDir)

	sAria := textinput.New()
	sAria.Prompt = "aria2c path: "
	sAria.Width = 50
	sAria.SetValue(cfg.Aria2Path)
	styleTextInput(&sAria)

	m := Model{
		cfg:           cfg,
		screen:        ScreenMenu,
		urlInput:      url,
		destInput:     dest,
		settingsDir:   sDir,
		settingsAria2: sAria,
	}
	m.destInput.SetValue(cfg.DownloadDir)

	switch start {
	case ViewDownload:
		m.screen = ScreenNewDownload
		m.newStep = 0
		m.urlInput.Focus()
	case ViewList:
		m.screen = ScreenActive
	case ViewSettings:
		m.screen = ScreenSettings
		m.settingsDir.Focus()
	default:
		m.screen = ScreenMenu
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.handleGlobalKeys(msg) {
			return m, nil
		}
		return m.updateKey(msg)

	case tickMsg:
		if m.daemon != nil && (m.screen == ScreenActive || m.screen == ScreenDetail) {
			return m, m.fetchDownloads()
		}
		return m, tickCmd()

	case downloadsMsg:
		if msg.err == nil {
			m.downloads = msg.items
		}
		return m, tickCmd()
	}

	var cmd tea.Cmd
	switch m.screen {
	case ScreenNewDownload:
		if m.browsing {
			m.browseList, cmd = m.browseList.Update(msg)
		} else if m.newStep == 0 {
			m.urlInput, cmd = m.urlInput.Update(msg)
		} else {
			m.destInput, cmd = m.destInput.Update(msg)
		}
	case ScreenSettings:
		if m.settingsFocus == 0 {
			m.settingsDir, cmd = m.settingsDir.Update(msg)
		} else {
			m.settingsAria2, cmd = m.settingsAria2.Update(msg)
		}
	}
	return m, cmd
}

func (m *Model) handleGlobalKeys(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "ctrl+c":
		m.cleanup()
		return true
	case "?":
		if m.screen != ScreenHelp {
			m.screen = ScreenHelp
		} else {
			m.screen = ScreenMenu
		}
		return true
	}
	return false
}

func (m Model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case ScreenMenu:
		return m.updateMenu(msg)
	case ScreenNewDownload:
		return m.updateNewDownload(msg)
	case ScreenActive:
		return m.updateActive(msg)
	case ScreenSettings:
		return m.updateSettings(msg)
	case ScreenHelp:
		return m.updateHelp(msg)
	case ScreenDetail:
		return m.updateDetail(msg)
	}
	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < len(mainMenuEntries)-1 {
			m.menuIndex++
		}
	case "q", "Q":
		m.cleanup()
		return m, tea.Quit
	case "enter", " ":
		item := mainMenuEntries[m.menuIndex]
		if item.quit {
			m.cleanup()
			return m, tea.Quit
		}
		switch item.screen {
		case ScreenNewDownload:
			m.screen = ScreenNewDownload
			m.newStep = 0
			m.urlInput.Focus()
		case ScreenActive:
			m.screen = ScreenActive
			m.selectedDL = 0
			return m, m.fetchDownloads()
		case ScreenSettings:
			m.screen = ScreenSettings
			m.settingsFocus = 0
			m.settingsDir.Focus()
		case ScreenHelp:
			m.screen = ScreenHelp
		}
	}
	return m, nil
}

func (m Model) updateNewDownload(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.browsing {
		switch msg.String() {
		case "esc":
			m.browsing = false
			m.destInput.Focus()
		case "enter":
			idx := m.browseList.Index()
			item := m.browseList.Items()[idx].(browseItem)
			if item.name == ".." {
				m.browseDir = filepath.Dir(m.browseDir)
			} else {
				m.browseDir = filepath.Join(m.browseDir, item.name)
			}
			m.refreshBrowse()
		case " ":
			m.destInput.SetValue(m.browseDir)
			m.browsing = false
			m.destInput.Focus()
		}
		var cmd tea.Cmd
		m.browseList, cmd = m.browseList.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "esc":
		m.screen = ScreenMenu
		m.errMsg = ""
	case "b":
		if m.newStep == 1 {
			m.browsing = true
			m.browseDir = m.destInput.Value()
			if m.browseDir == "" {
				m.browseDir = m.cfg.DownloadDir
			}
			m.refreshBrowse()
		}
	case "enter":
		if m.newStep == 0 {
			m.newStep = 1
			if m.destInput.Value() == "" {
				m.destInput.SetValue(m.cfg.DownloadDir)
			}
			m.destInput.Focus()
		} else {
			return m.startDownload()
		}
	default:
		if m.newStep == 0 {
			var cmd tea.Cmd
			m.urlInput, cmd = m.urlInput.Update(msg)
			return m, cmd
		}
		var cmd tea.Cmd
		m.destInput, cmd = m.destInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) startDownload() (tea.Model, tea.Cmd) {
	if m.daemon == nil {
		m.errMsg = "daemon not running"
		return m, nil
	}
	in, err := download.Parse(m.urlInput.Value())
	if err != nil {
		m.errMsg = err.Error()
		return m, nil
	}
	dir := strings.TrimSpace(m.destInput.Value())
	if dir == "" {
		dir = m.cfg.DownloadDir
	}
	client := m.daemon.Client()
	switch in.Kind {
	case download.KindHTTP, download.KindMagnet:
		_, err = client.AddURI(in.Raw, dir)
	case download.KindTorrentFile:
		_, err = client.AddTorrent(in.Torrent, dir)
	}
	if err != nil {
		m.errMsg = err.Error()
		return m, nil
	}
	m.errMsg = ""
	m.status = "Download started"
	m.urlInput.SetValue("")
	m.newStep = 0
	m.screen = ScreenActive
	m.selectedDL = 0
	return m, m.fetchDownloads()
}

func (m Model) updateActive(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.screen = ScreenMenu
	case "j", "down":
		if m.selectedDL < len(m.downloads)-1 {
			m.selectedDL++
		}
	case "k", "up":
		if m.selectedDL > 0 {
			m.selectedDL--
		}
	case "p":
		if len(m.downloads) > 0 && m.daemon != nil {
			_ = m.daemon.Client().Pause(m.downloads[m.selectedDL].GID)
		}
	case "r":
		if len(m.downloads) > 0 && m.daemon != nil {
			_ = m.daemon.Client().Unpause(m.downloads[m.selectedDL].GID)
		}
	case "x":
		if len(m.downloads) > 0 && m.daemon != nil {
			_ = m.daemon.Client().Remove(m.downloads[m.selectedDL].GID)
		}
		return m, m.fetchDownloads()
	case "enter":
		if len(m.downloads) > 0 {
			m.detail = m.downloads[m.selectedDL]
			m.screen = ScreenDetail
		}
	}
	return m, m.fetchDownloads()
}

func (m Model) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" || msg.String() == "enter" {
		m.screen = ScreenActive
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = ScreenMenu
	case "tab":
		m.settingsFocus = 1 - m.settingsFocus
		if m.settingsFocus == 0 {
			m.settingsDir.Focus()
			m.settingsAria2.Blur()
		} else {
			m.settingsAria2.Focus()
			m.settingsDir.Blur()
		}
	case "enter":
		m.cfg.DownloadDir = strings.TrimSpace(m.settingsDir.Value())
		m.cfg.Aria2Path = strings.TrimSpace(m.settingsAria2.Value())
		if err := config.Save(m.cfg); err != nil {
			m.errMsg = err.Error()
		} else {
			m.status = "Settings saved"
			m.errMsg = ""
		}
	default:
		if m.settingsFocus == 0 {
			var cmd tea.Cmd
			m.settingsDir, cmd = m.settingsDir.Update(msg)
			return m, cmd
		}
		var cmd tea.Cmd
		m.settingsAria2, cmd = m.settingsAria2.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" || msg.String() == "q" || msg.String() == "enter" {
		m.screen = ScreenMenu
	}
	return m, nil
}

func (m *Model) fetchDownloads() tea.Cmd {
	return func() tea.Msg {
		if m.daemon == nil {
			return downloadsMsg{}
		}
		items, err := m.daemon.Client().TellActive()
		return downloadsMsg{items: items, err: err}
	}
}

func (m *Model) refreshBrowse() {
	entries, _ := os.ReadDir(m.browseDir)
	items := []list.Item{browseItem{name: "..", dir: true}}
	for _, e := range entries {
		if e.IsDir() {
			items = append(items, browseItem{name: e.Name(), dir: true})
		}
	}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(colorGreen)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(colorMuted)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(colorWhite)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(colorMuted)
	l := list.New(items, delegate, 50, 12)
	l.Title = sectionStyle.Render("Browse: " + m.browseDir)
	l.SetShowStatusBar(false)
	m.browseList = l
}

type browseItem struct {
	name string
	dir  bool
}

func (b browseItem) Title() string       { return b.name }
func (b browseItem) Description() string { return "" }
func (b browseItem) FilterValue() string { return b.name }

func (m *Model) cleanup() {
	if m.daemon != nil {
		_ = m.daemon.Stop()
		m.daemon = nil
	}
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading...\n"
	}
	var b strings.Builder
	b.WriteString(renderHeader())

	switch m.screen {
	case ScreenMenu:
		b.WriteString(m.viewMenu())
	case ScreenNewDownload:
		b.WriteString("\n")
		b.WriteString(m.viewNewDownload())
	case ScreenActive:
		b.WriteString("\n")
		b.WriteString(m.viewActive())
	case ScreenSettings:
		b.WriteString("\n")
		b.WriteString(m.viewSettings())
	case ScreenHelp:
		b.WriteString("\n")
		b.WriteString(m.viewHelp())
	case ScreenDetail:
		b.WriteString("\n")
		b.WriteString(m.viewDetail())
	}

	if m.errMsg != "" {
		b.WriteString("\n" + errStyle.Render(m.errMsg))
	} else if m.status != "" {
		b.WriteString("\n" + okStyle.Render(m.status))
	}
	return b.String()
}

func (m Model) viewNewDownload() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("New Download") + "\n\n")
	if m.browsing {
		b.WriteString(m.browseList.View())
		b.WriteString("\n" + renderFooter("Enter", "open", "Space", "select", "Esc", "back"))
		return b.String()
	}
	if m.newStep == 0 {
		b.WriteString(m.urlInput.View())
		b.WriteString("\n\n" + renderFooter("Enter", "next", "Esc", "menu"))
	} else {
		b.WriteString(helpStyle.Render("URL: ") + m.urlInput.Value() + "\n\n")
		b.WriteString(m.destInput.View())
		b.WriteString("\n\n" + renderFooter("B", "browse", "Enter", "start", "Esc", "menu"))
	}
	return b.String()
}

func (m Model) viewActive() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("Active Downloads") + "\n\n")
	if len(m.downloads) == 0 {
		b.WriteString(helpStyle.Render("No active downloads.") + "\n")
	} else {
		for i, d := range m.downloads {
			prefix := "  "
			if i == m.selectedDL {
				prefix = cursorStyle.Render(">") + " "
			}
			name := filepath.Base(d.Name)
			if name == "" {
				name = d.GID
			}
			if i == m.selectedDL {
				name = menuActive.Render(name)
			}
			bar := aria2.ProgressBar(d.Completed, d.Total, 20)
			line := fmt.Sprintf("%s%s\n    %s  %s  ETA %s  %s\n",
				prefix, name, bar, aria2.FormatSpeed(d.Speed), aria2.FormatETA(d.ETA), menuDescStyle.Render("["+d.Status+"]"))
			b.WriteString(line)
		}
	}
	b.WriteString("\n" + renderFooter("J/K", "select", "P", "pause", "R", "resume", "X", "cancel", "Enter", "details", "Esc", "menu"))
	return b.String()
}

func (m Model) viewDetail() string {
	d := m.detail
	return boxStyle.Render(fmt.Sprintf(
		"GID: %s\nName: %s\nStatus: %s\nDir: %s\nProgress: %s / %s\nSpeed: %s\nETA: %s\n\nesc back",
		d.GID, d.Name, d.Status, d.Dir,
		aria2.FormatBytes(d.Completed), aria2.FormatBytes(d.Total),
		aria2.FormatSpeed(d.Speed), aria2.FormatETA(d.ETA),
	))
}

func (m Model) viewSettings() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("Settings") + "\n\n")
	b.WriteString(m.settingsDir.View() + "\n")
	b.WriteString(m.settingsAria2.View() + "\n")
	b.WriteString("\n" + renderFooter("Tab", "switch", "Enter", "save", "Esc", "menu"))
	return b.String()
}

func (m Model) viewHelp() string {
	body := sectionStyle.Render("Help") + `

  adl              Main menu
  adl download     New download
  adl list         Active downloads
  adl settings     Settings

  Config: ~/.config/adl/config.yaml
  Requires: aria2c (brew install aria2)

` + renderFooter("Esc", "back")
	return boxStyle.Render(body)
}

func Run(cfg config.Config, start StartView) error {
	daemon, err := aria2.Start(cfg.Aria2Path, cfg.RPCPort, cfg.DownloadDir)
	if err != nil {
		return err
	}

	m := NewModel(cfg, start)
	m.daemon = daemon

	finalModel, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		_ = daemon.Stop()
		return err
	}
	if fm, ok := finalModel.(Model); ok {
		fm.cleanup()
	} else {
		_ = daemon.Stop()
	}
	return nil
}
