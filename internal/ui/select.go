package ui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

const clipboardLimit = 50 * 1024

func (m Model) startSelection() Model {
	list := m.detail.selectable()
	if len(list) == 0 {
		return m
	}

	anchor := 0
	if top := m.detail.messageAt(m.detailScroll); top >= 0 {
		anchor = len(list) - 1
		for i, msg := range list {
			if msg >= top {
				anchor = i
				break
			}
		}
	}

	m.mode = modeSelect
	m.selection = list
	m.selectAnchor = anchor
	m.selectCursor = anchor
	m.notice = ""
	m = m.revealSelection()
	return m
}

func (m Model) selectedRange() (int, int, bool) {
	if m.mode != modeSelect || len(m.selection) == 0 {
		return 0, 0, false
	}
	lo, hi := m.selectAnchor, m.selectCursor
	if lo > hi {
		lo, hi = hi, lo
	}
	return m.selection[lo], m.selection[hi], true
}

func (m Model) selectedCount() int {
	lo, hi, ok := m.selectedRange()
	if !ok {
		return 0
	}
	count := 0
	for _, msg := range m.selection {
		if msg >= lo && msg <= hi {
			count++
		}
	}
	return count
}

func (m Model) revealSelection() Model {
	if len(m.selection) == 0 {
		return m
	}
	start, end, ok := m.detail.messageSpan(m.selection[m.selectCursor])
	if !ok {
		return m
	}
	height := m.bodyHeight()
	if start < m.detailScroll {
		m.detailScroll = start
	}
	if height > 0 && end >= m.detailScroll+height {
		m.detailScroll = end - height + 1
	}
	m.clampScroll()
	return m
}

func (m Model) handleSelectKey(key string) (tea.Model, tea.Cmd) {
	m.notice = ""

	switch key {
	case "esc", "q", "v":
		m.mode = modeDetail
		m.selection = nil
		return m, nil
	case "j", "down":
		m.selectCursor++
	case "k", "up":
		m.selectCursor--
	case "J":
		m.selectCursor += 10
	case "K":
		m.selectCursor -= 10
	case "g", "home":
		m.selectCursor = 0
	case "G", "end":
		m.selectCursor = len(m.selection) - 1
	case "y":
		return m.copySelection()
	case "w":
		return m.writeSelection()
	}

	if m.selectCursor < 0 {
		m.selectCursor = 0
	}
	if m.selectCursor >= len(m.selection) {
		m.selectCursor = len(m.selection) - 1
	}
	m = m.revealSelection()
	return m, nil
}

func (m Model) markSelection(lines []string) []string {
	lo, hi, ok := m.selectedRange()
	if !ok {
		return lines
	}
	cursor := m.selection[m.selectCursor]
	for i, line := range lines {
		msg, first := m.detail.lineInfo(m.detailScroll + i)
		if msg < 0 || !first {
			continue
		}
		marker := ""
		switch {
		case msg == cursor:
			marker = ">"
		case msg >= lo && msg <= hi:
			marker = "*"
		}
		if marker == "" || len(line) == 0 {
			continue
		}
		if line[0] == ' ' {
			lines[i] = marker + line[1:]
		}
	}
	return lines
}

func (m Model) markdownSelection() (string, int, int, int) {
	lo, hi, ok := m.selectedRange()
	if !ok {
		return "", 0, 0, 0
	}
	payload := exportMarkdown(m.conversation.Messages[lo:hi+1], m.hidden)
	return payload, lo, hi, m.selectedCount()
}

func (m Model) copySelection() (tea.Model, tea.Cmd) {
	payload, _, _, count := m.markdownSelection()
	if count == 0 {
		return m, nil
	}

	if len(payload) > clipboardLimit {
		m.notice = "selection too large for the clipboard, press w to write a file"
		return m, nil
	}

	lookPath := m.lookPath
	if lookPath == nil {
		lookPath = execLookPath
	}
	run := m.copyText
	if run == nil {
		run = runClipboard
	}

	if tool, ok := detectClipboard(lookPath); ok {
		if err := run(tool, payload); err == nil {
			m.notice = fmt.Sprintf("copied %d messages", count)
			return m, nil
		}
	}

	sequence := m.clipboard
	if sequence == nil {
		sequence = ansi.SetSystemClipboard
	}
	m.notice = fmt.Sprintf("copied %d messages (terminal)", count)
	return m, writeStdout(sequence(payload))
}

func (m Model) writeSelection() (tea.Model, tea.Cmd) {
	payload, lo, hi, count := m.markdownSelection()
	if count == 0 {
		return m, nil
	}
	return m.writeExportFile(payload, lo, hi)
}

func (m Model) exportThread() (tea.Model, tea.Cmd) {
	messages := m.conversation.Messages
	if len(messages) == 0 {
		return m, nil
	}

	first, last := -1, -1
	for i, msg := range messages {
		if !hasVisibleBlock(msg, m.hidden) {
			continue
		}
		if first < 0 {
			first = i
		}
		last = i
	}
	if first < 0 {
		return m, nil
	}

	payload := exportMarkdown(messages[first:last+1], m.hidden)
	return m.writeExportFile(payload, first, last)
}

func (m Model) writeExportFile(payload string, first, last int) (tea.Model, tea.Cmd) {
	dir := m.exportDir
	if dir == "" {
		dir = "."
	}
	name := exportName(shortID(m.conversation.Session.AtuinID), first, last)

	write := m.writeFile
	if write == nil {
		write = writeExport
	}
	written, err := write(dir, name, payload)
	if err != nil {
		m.notice = "could not write " + name
		return m, nil
	}
	path := written
	if dir != "." {
		path = filepath.Join(dir, written)
	}
	m.notice = "wrote " + path
	return m, nil
}

func writeStdout(sequence string) tea.Cmd {
	return func() tea.Msg {
		fmt.Fprint(os.Stdout, sequence)
		return nil
	}
}
