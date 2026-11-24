package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
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
	m.selectCursor = anchor
	m.marked = make(map[int]bool)
	m.notice = ""
	m = m.revealSelection()
	return m
}

func (m Model) cursorMessage() int {
	if m.mode != modeSelect || len(m.selection) == 0 {
		return -1
	}
	return m.selection[m.selectCursor]
}

func (m Model) markedCount() int {
	return len(m.marked)
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
	switch {
	case height > 0 && end-start+1 >= height:
		m.detailScroll = start
	case start < m.detailScroll:
		m.detailScroll = start
	case height > 0 && end >= m.detailScroll+height:
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
		m.marked = nil
		return m, nil
	case " ":
		if msg := m.cursorMessage(); msg >= 0 {
			if m.marked[msg] {
				delete(m.marked, msg)
			} else {
				m.marked[msg] = true
			}
		}
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

const (
	backgroundEven   = "\x1b[48;5;234m"
	backgroundOdd    = "\x1b[48;5;236m"
	backgroundCursor = "\x1b[48;5;240m"
	backgroundMarked = "\x1b[48;5;238m"
)

func paintBackground(line, background string, width int) string {
	if width <= 0 {
		return line
	}
	head := ansi.Truncate(line, width, "")
	rest := ansi.TruncateLeft(line, width, "")
	if used := ansi.StringWidth(head); used < width {
		head += strings.Repeat(" ", width-used)
	}
	head = strings.ReplaceAll(head, "\x1b[0m", "\x1b[0m"+background)
	return background + head + "\x1b[0m" + rest
}

func replaceMarginCell(line, cell, background string) string {
	return strings.Replace(line, background+" ", background+cell, 1)
}

func (m Model) decorateDetail(lines []string) []string {
	layout := newDetailLayout(m.width)
	width := detailMargin + layout.body
	if width <= 0 {
		return lines
	}

	out := make([]string, len(lines))
	cursor := m.cursorMessage()
	for i, line := range lines {
		msg, _ := m.detail.lineInfo(m.detailScroll + i)
		if msg < 0 {
			out[i] = line
			continue
		}

		background := backgroundOdd
		if msg%2 == 0 {
			background = backgroundEven
		}
		cell := " "
		if m.mode == modeSelect {
			switch {
			case msg == cursor:
				background = backgroundCursor
				cell = ">"
			case m.marked[msg]:
				background = backgroundMarked
				cell = "#"
			}
		}

		out[i] = paintBackground(line, background, width)
		out[i] = replaceMarginCell(out[i], cell, background)
	}
	return out
}

func (m Model) markedMessages() []atuin.Message {
	if len(m.marked) == 0 {
		return nil
	}
	messages := make([]atuin.Message, 0, len(m.marked))
	for i, msg := range m.conversation.Messages {
		if m.marked[i] {
			messages = append(messages, msg)
		}
	}
	return messages
}

func (m Model) markdownSelection() (string, int, int, int) {
	messages := m.markedMessages()
	if len(messages) == 0 {
		return "", 0, 0, 0
	}

	first, last := -1, -1
	for i := range m.conversation.Messages {
		if !m.marked[i] {
			continue
		}
		if first < 0 {
			first = i
		}
		last = i
	}
	return exportMarkdown(messages, m.hidden), first, last, len(messages)
}

func (m Model) copySelection() (tea.Model, tea.Cmd) {
	payload, _, _, count := m.markdownSelection()
	if count == 0 {
		m.notice = "nothing marked, press space to mark a message"
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
	sequence := m.clipboard
	if sequence == nil {
		sequence = ansi.SetSystemClipboard
	}

	return m, func() tea.Msg {
		if tool, ok := detectClipboard(lookPath); ok {
			if err := run(tool, payload); err == nil {
				return copiedMsg{count: count}
			}
		}
		writeStdout(sequence(payload))
		return copiedMsg{count: count, terminal: true}
	}
}

func (m Model) writeSelection() (tea.Model, tea.Cmd) {
	payload, lo, hi, count := m.markdownSelection()
	if count == 0 {
		m.notice = "nothing marked, press space to mark a message"
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

func writeStdout(sequence string) {
	fmt.Fprint(os.Stdout, sequence)
}
