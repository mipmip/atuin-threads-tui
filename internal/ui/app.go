package ui

import (
	"context"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

type mode int

const (
	modeList mode = iota
	modeDetail
	modeFilter
)

type threadsMsg struct {
	threads []atuin.Thread
	err     error
}

type conversationMsg struct {
	conversation atuin.Conversation
	err          error
}

type searchMsg struct {
	query   string
	threads []atuin.Thread
	err     error
}

type detailRenderedMsg struct {
	rendering detailRendering
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
	detail       detailRendering
	rendering    bool
	builds       int

	filter      string
	filterInput string
	searchQuery string
	searching   bool

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
	threads := m.visible()
	if m.cursor < 0 || m.cursor >= len(threads) {
		return atuin.Thread{}, false
	}
	return threads[m.cursor], true
}

// Update handles terminal input and load results.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		widthChanged := msg.Width != m.width
		m.width = msg.Width
		m.height = msg.Height
		m.clampCursor()
		if m.mode == modeDetail && (widthChanged || !m.detail.built) {
			return m.startDetail()
		}
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
		return m.startDetail()

	case detailRenderedMsg:
		if m.mode != modeDetail || msg.rendering.width != m.width {
			return m, nil
		}
		m.rendering = false
		m.detail = msg.rendering
		m.clampScroll()
		return m, nil

	case searchMsg:
		m.searching = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.err = nil
		m.threads = msg.threads
		m.searchQuery = msg.query
		m.filter = ""
		m.filterInput = ""
		m.cursor = 0
		m.offset = 0
		return m, nil

	case tea.KeyMsg:
		if m.mode != modeFilter && len(msg.Runes) > 1 {
			for _, r := range msg.Runes {
				updated, cmd := m.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}, Alt: msg.Alt})
				m = updated.(Model)
				if cmd != nil {
					return m, cmd
				}
			}
			return m, nil
		}
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.mode {
	case modeDetail:
		return m.handleDetailKey(key)
	case modeFilter:
		return m.handleFilterKey(key, msg)
	}

	return m.handleListKey(key)
}

func (m Model) handleListKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q":
		return m, tea.Quit
	case "ctrl+r":
		m.loading = true
		m.err = nil
		m.searchQuery = ""
		return m, m.listThreads()
	case "/":
		m.mode = modeFilter
		m.filterInput = m.filter
		return m, nil
	case "esc":
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.filter = ""
			m.filterInput = ""
			m.loading = true
			m.err = nil
			return m, m.listThreads()
		}
		m.filter = ""
		m.filterInput = ""
		m.cursor = 0
		m.offset = 0
	case "ctrl+s":
		if query := strings.TrimSpace(m.filter); query != "" {
			m.searching = true
			m.err = nil
			return m, m.searchThreads()
		}
	case "j", "down":
		m.cursor++
	case "k", "up":
		m.cursor--
	case "g", "home":
		m.cursor = 0
	case "G", "end":
		m.cursor = len(m.visible()) - 1
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
		m.detail = detailRendering{}
		m.rendering = false
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
		m.detailScroll = m.detailTotal()
	}

	m.clampScroll()
	return m, nil
}

func (m *Model) clampCursor() {
	threads := m.visible()
	if len(threads) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(threads) {
		m.cursor = len(threads) - 1
	}

	height := m.listHeight()
	if height <= 0 {
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+height {
		m.offset = m.cursor - height + 1
	}
	if maxOffset := len(threads) - height; m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *Model) clampScroll() {
	maxScroll := m.detailTotal() - m.bodyHeight()
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

func (m Model) detailTotal() int {
	if m.detail.built {
		return m.detail.total
	}
	return 0
}

// startDetail renders the open conversation, in the background when it is too
// large to do without blocking the interface.
func (m Model) startDetail() (tea.Model, tea.Cmd) {
	m.builds++
	m.detail = detailRendering{}
	m.rendering = false

	if len(m.conversation.Messages) <= detailSyncMessages {
		m.detail = m.buildDetail()
		m.clampScroll()
		return m, nil
	}

	m.rendering = true
	return m, m.renderDetail()
}

func (m Model) renderDetail() tea.Cmd {
	model := Model{
		mode:         modeDetail,
		width:        m.width,
		conversation: m.conversation,
		now:          m.now,
	}
	return func() tea.Msg {
		return detailRenderedMsg{rendering: model.buildDetail()}
	}
}

func (m Model) bodyHeight() int {
	if m.height-3 < 0 {
		return 0
	}
	return m.height - 3
}

func (m Model) listHeight() int {
	if height := m.bodyHeight() - 2; height > 0 {
		return height
	}
	return 0
}
