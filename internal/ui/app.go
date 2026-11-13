package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

type mode int

const (
	modeList mode = iota
	modeDetail
)

type threadsMsg struct {
	threads []atuin.Thread
	err     error
}

type conversationMsg struct {
	conversation atuin.Conversation
	err          error
}

// Model is the bubbletea model for the browser.
type Model struct {
	source atuin.Source
	now    func() time.Time

	threads []atuin.Thread
	cursor  int
	offset  int

	mode         mode
	conversation atuin.Conversation
	detailScroll int

	width  int
	height int

	loading bool
	err     error
}

// New returns a model that reads threads from source.
func New(source atuin.Source) Model {
	return Model{
		source:  source,
		now:     time.Now,
		loading: true,
		width:   80,
		height:  24,
	}
}

// Init starts the first load.
func (m Model) Init() tea.Cmd {
	return m.listThreads()
}

func (m Model) listThreads() tea.Cmd {
	return func() tea.Msg {
		threads, err := m.source.List(context.Background())
		return threadsMsg{threads: threads, err: err}
	}
}

func (m Model) openThread() tea.Cmd {
	thread, ok := m.selected()
	if !ok {
		return nil
	}
	return func() tea.Msg {
		conversation, err := m.source.Show(context.Background(), thread.AtuinID)
		return conversationMsg{conversation: conversation, err: err}
	}
}

func (m Model) selected() (atuin.Thread, bool) {
	if m.cursor < 0 || m.cursor >= len(m.threads) {
		return atuin.Thread{}, false
	}
	return m.threads[m.cursor], true
}

// Update handles terminal input and load results.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clampCursor()
		m.clampScroll()
		return m, nil

	case threadsMsg:
		m.loading = false
		m.err = msg.err
		if msg.err != nil {
			return m, nil
		}
		m.threads = msg.threads
		m.cursor = 0
		m.offset = 0
		return m, nil

	case conversationMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.err = nil
		m.conversation = msg.conversation
		m.mode = modeDetail
		m.detailScroll = 0
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	if m.mode == modeDetail {
		return m.handleDetailKey(key)
	}

	switch key {
	case "q":
		return m, tea.Quit
	case "ctrl+r":
		m.loading = true
		m.err = nil
		return m, m.listThreads()
	case "j", "down":
		m.cursor++
	case "k", "up":
		m.cursor--
	case "g", "home":
		m.cursor = 0
	case "G", "end":
		m.cursor = len(m.threads) - 1
	case "enter":
		return m, m.openThread()
	}

	m.clampCursor()
	return m, nil
}

func (m Model) handleDetailKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc", "left", "h":
		m.mode = modeList
		m.err = nil
		return m, nil
	case "j", "down":
		m.detailScroll++
	case "k", "up":
		m.detailScroll--
	case "pgdown", " ":
		m.detailScroll += m.bodyHeight()
	case "pgup":
		m.detailScroll -= m.bodyHeight()
	case "g", "home":
		m.detailScroll = 0
	case "G", "end":
		m.detailScroll = len(m.detailLines())
	}

	m.clampScroll()
	return m, nil
}

func (m *Model) clampCursor() {
	if len(m.threads) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.threads) {
		m.cursor = len(m.threads) - 1
	}

	height := m.bodyHeight()
	if height <= 0 {
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+height {
		m.offset = m.cursor - height + 1
	}
	if maxOffset := len(m.threads) - height; m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *Model) clampScroll() {
	maxScroll := len(m.detailLines()) - m.bodyHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.detailScroll > maxScroll {
		m.detailScroll = maxScroll
	}
	if m.detailScroll < 0 {
		m.detailScroll = 0
	}
}

func (m Model) bodyHeight() int {
	if m.height-3 < 0 {
		return 0
	}
	return m.height - 3
}
