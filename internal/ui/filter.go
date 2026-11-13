package ui

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func (m Model) visible() []atuin.Thread {
	query := strings.ToLower(strings.TrimSpace(m.filter))
	if query == "" {
		return m.threads
	}

	out := make([]atuin.Thread, 0, len(m.threads))
	for _, thread := range m.threads {
		if matches(thread, query) {
			out = append(out, thread)
		}
	}
	return out
}

func matches(thread atuin.Thread, query string) bool {
	fields := []string{
		thread.Title,
		thread.Harness,
		thread.Model,
		thread.CWD,
		thread.Preview,
		thread.LastReply,
		thread.GitBranch,
		thread.AtuinID,
		thread.SessionID,
	}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}

func (m Model) searchThreads() tea.Cmd {
	query := strings.TrimSpace(m.filter)
	return func() tea.Msg {
		threads, err := m.source.Search(context.Background(), query, 100)
		return searchMsg{query: query, threads: threads, err: err}
	}
}

func (m Model) handleFilterKey(key string, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.filter = ""
		m.filterInput = ""
		m.mode = modeList
		m.cursor = 0
		m.offset = 0
		m.clampCursor()
		return m, nil

	case "enter":
		m.filter = strings.TrimSpace(m.filterInput)
		m.mode = modeList
		m.cursor = 0
		m.offset = 0
		m.clampCursor()
		return m, nil

	case "backspace":
		runes := []rune(m.filterInput)
		if len(runes) > 0 {
			m.filterInput = string(runes[:len(runes)-1])
		}

	case "ctrl+u":
		m.filterInput = ""

	case " ":
		m.filterInput += " "

	default:
		if len(msg.Runes) > 0 && !msg.Alt {
			m.filterInput += string(msg.Runes)
		}
	}

	m.filter = m.filterInput
	m.cursor = 0
	m.offset = 0
	m.clampCursor()
	return m, nil
}
