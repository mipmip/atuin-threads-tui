package ui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

const maxBlockLines = 8

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	mutedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	selectedStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	userStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	assistantStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	systemStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("241"))
)

// View renders the whole screen.
func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("att"))
	b.WriteString(mutedStyle.Render(" · atuin threads"))
	b.WriteString("\n")

	body := m.body()
	if height := m.bodyHeight(); height > 0 {
		if lines := strings.Count(body, "\n") + 1; lines < height {
			body += strings.Repeat("\n", height-lines)
		}
	}
	b.WriteString(body)
	b.WriteString("\n")
	b.WriteString(m.status())
	b.WriteString("\n")
	b.WriteString(m.hints())

	return b.String()
}

func (m Model) body() string {
	if m.mode == modeDetail {
		return m.detailBody()
	}

	switch {
	case m.loading:
		return mutedStyle.Render("loading threads...")
	case m.err != nil && len(m.threads) == 0:
		return errorStyle.Render("error: "+m.err.Error()) + "\n" +
			mutedStyle.Render("ctrl+r to retry")
	case len(m.threads) == 0:
		return mutedStyle.Render("no threads captured yet") + "\n" +
			mutedStyle.Render("enable capture with: atuin config enable threads")
	}

	threads := m.visible()
	if len(threads) == 0 {
		return mutedStyle.Render("no thread matches the filter")
	}

	height := m.bodyHeight()
	end := m.offset + height
	if end > len(threads) {
		end = len(threads)
	}

	rows := make([]string, 0, height)
	now := m.now()
	for i := m.offset; i < end; i++ {
		rows = append(rows, threadRow(m.width, threads[i], i == m.cursor, now))
	}
	return strings.Join(rows, "\n")
}

func (m Model) detailBody() string {
	lines := m.detailLines()
	start := m.detailScroll
	if start > len(lines) {
		start = len(lines)
	}
	end := start + m.bodyHeight()
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n")
}

func (m Model) status() string {
	if m.err != nil {
		return errorStyle.Render("error: " + m.err.Error())
	}

	switch m.mode {
	case modeFilter:
		counts := fmt.Sprintf("%d threads", len(m.threads))
		if strings.TrimSpace(m.filterInput) != "" {
			counts = fmt.Sprintf("%d of %d threads", len(m.visible()), len(m.threads))
		}
		return mutedStyle.Render("filter: ") + m.filterInput + "▌" + mutedStyle.Render("  "+counts)
	case modeDetail:
		return mutedStyle.Render(fmt.Sprintf(
			"%s · %d messages · line %d of %d",
			shortID(m.conversation.Session.AtuinID),
			len(m.conversation.Messages),
			m.detailScroll+1,
			len(m.detailLines()),
		))
	}

	visible := len(m.visible())
	switch {
	case m.searching:
		return mutedStyle.Render("searching...")
	case m.searchQuery != "":
		return mutedStyle.Render(fmt.Sprintf("%d results for %q", visible, m.searchQuery))
	case m.filter != "":
		return mutedStyle.Render(fmt.Sprintf("%d of %d threads", visible, len(m.threads)))
	case m.loading:
		return mutedStyle.Render("loading")
	default:
		return mutedStyle.Render(fmt.Sprintf("%d threads", len(m.threads)))
	}
}

func (m Model) hints() string {
	switch m.mode {
	case modeDetail:
		return mutedStyle.Render("j/k scroll · pgup/pgdn page · g/G top/bottom · esc back · q quit")
	case modeFilter:
		return mutedStyle.Render("type to filter · enter apply · ctrl+u reset · esc clear")
	default:
		return mutedStyle.Render("j/k move · enter read · / filter · ctrl+s full-text · ctrl+r reload · q quit")
	}
}

func threadRow(width int, t atuin.Thread, selected bool, now time.Time) string {
	marker := "  "
	if selected {
		marker = "> "
	}

	title := strings.TrimSpace(t.Title)
	if title == "" {
		title = "(untitled)"
	}

	line := fmt.Sprintf(
		"%s%s  %s  ·  %s  ·  %d msgs  ·  %s",
		marker, shortID(t.AtuinID), title, t.Harness, t.MessageCount, age(now, t.UpdatedAt),
	)
	line = truncate(line, width)

	if selected {
		return selectedStyle.Render(line)
	}
	return line
}

func (m Model) detailLines() []string {
	if m.mode != modeDetail {
		return nil
	}

	session := m.conversation.Session
	width := m.width - 2
	if width < 10 {
		width = 10
	}

	title := strings.TrimSpace(session.Title)
	if title == "" {
		title = "(untitled)"
	}
	lines := []string{titleStyle.Render(truncate(title, width))}

	var meta []string
	if session.Harness != "" {
		meta = append(meta, session.Harness)
	}
	if session.Model != "" {
		meta = append(meta, session.Model)
	}
	if session.AtuinID != "" {
		meta = append(meta, shortID(session.AtuinID))
	}
	if session.MessageCount > 0 {
		meta = append(meta, fmt.Sprintf("%d msgs", session.MessageCount))
	}
	if !session.UpdatedAt.IsZero() {
		meta = append(meta, age(m.now(), session.UpdatedAt)+" ago")
	}
	if len(meta) > 0 {
		lines = append(lines, mutedStyle.Render(truncate(strings.Join(meta, " · "), width)))
	}
	if session.CWD != "" {
		lines = append(lines, mutedStyle.Render(truncate(session.CWD, width)))
	}
	lines = append(lines, "")

	for _, msg := range m.conversation.Messages {
		lines = append(lines, roleLine(msg))
		for _, block := range msg.Content {
			lines = append(lines, blockLines(block, width)...)
		}
		lines = append(lines, "")
	}

	return lines
}

func roleLine(msg atuin.Message) string {
	role := msg.Role
	if role == "" {
		role = "unknown"
	}

	style := systemStyle
	switch role {
	case "user":
		style = userStyle
	case "assistant":
		style = assistantStyle
	}

	line := style.Render("[" + role + "]")
	if !msg.Timestamp.IsZero() {
		line += mutedStyle.Render(" " + msg.Timestamp.Format("15:04"))
	}
	return line
}

func blockLines(block atuin.Content, width int) []string {
	const indent = "  "
	const bodyIndent = "    "

	switch block.Type {
	case "text":
		return indentLines(wrap(block.Text, width-len(indent)), indent)

	case "thinking":
		lines := indentLines(wrap(block.Text, width-len(indent)), indent)
		for i := range lines {
			lines[i] = mutedStyle.Render(lines[i])
		}
		return lines

	case "tool_call":
		lines := []string{indent + "tool " + block.Name}
		if input := strings.TrimSpace(block.Input); input != "" {
			lines = append(lines, indentLines(wrap(input, width-len(bodyIndent)), bodyIndent)...)
		}
		return capped(lines)

	case "tool_result":
		label := indent + "result"
		if block.IsError {
			label += " (error)"
		}
		lines := []string{label}
		if content := strings.TrimSpace(block.Content); content != "" {
			lines = append(lines, indentLines(wrap(content, width-len(bodyIndent)), bodyIndent)...)
		}
		return capped(lines)
	}

	return nil
}

func capped(lines []string) []string {
	if len(lines) <= maxBlockLines {
		return lines
	}
	out := append([]string{}, lines[:maxBlockLines]...)
	return append(out, fmt.Sprintf("    ... (%d more lines)", len(lines)-maxBlockLines))
}

func indentLines(lines []string, indent string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = indent + line
	}
	return out
}

func wrap(text string, width int) []string {
	if width < 1 {
		width = 1
	}

	var out []string
	for _, paragraph := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}

		line := ""
		for _, word := range words {
			for _, chunk := range hardSplit(word, width) {
				switch {
				case line == "":
					line = chunk
				case utf8.RuneCountInString(line)+1+utf8.RuneCountInString(chunk) <= width:
					line += " " + chunk
				default:
					out = append(out, line)
					line = chunk
				}
			}
		}
		out = append(out, line)
	}

	return out
}

func hardSplit(word string, width int) []string {
	runes := []rune(word)
	if len(runes) <= width {
		return []string{word}
	}

	var chunks []string
	for len(runes) > width {
		chunks = append(chunks, string(runes[:width]))
		runes = runes[width:]
	}
	if len(runes) > 0 {
		chunks = append(chunks, string(runes))
	}
	return chunks
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return string(runes[:width-1]) + "…"
}

func shortID(id string) string {
	runes := []rune(id)
	if len(runes) <= 8 {
		return id
	}
	return string(runes[len(runes)-8:])
}

func age(now, then time.Time) string {
	if then.IsZero() {
		return "?"
	}
	d := now.Sub(then)
	if d < 0 {
		return "now"
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
