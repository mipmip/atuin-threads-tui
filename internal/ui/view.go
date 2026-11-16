package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

const maxBlockLines = 8

const (
	markerWidth   = 2
	listGutter    = 2
	idWidth       = 8
	dateWidth     = 12
	agentMinWidth = 2
	agentMaxWidth = 6
	msgsWidth     = 4
	titleMin      = 20
	titleMax      = 48
	dirMin        = 8
	dirMax        = 24
)

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
		lines := strings.Split(body, "\n")
		if len(lines) > height {
			lines = lines[:height]
		}
		for len(lines) < height {
			lines = append(lines, "")
		}
		if m.mode == modeDetail {
			lines = m.withScrollbar(lines)
		}
		body = strings.Join(lines, "\n")
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
	if m.mode == modeTypes {
		return m.typesBody()
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

	layout := newListLayout(m.width, measureThreads(threads))
	lines := []string{mutedStyle.Render(layout.header())}
	if m.bodyHeight() >= 2 {
		lines = append(lines, mutedStyle.Render(layout.rule()))
	}

	height := m.listHeight()
	end := m.offset + height
	if end > len(threads) {
		end = len(threads)
	}

	for i := m.offset; i < end; i++ {
		lines = append(lines, layout.row(threads[i], i == m.cursor))
	}
	return strings.Join(lines, "\n")
}

func (m Model) detailBody() string {
	if m.rendering || !m.detail.built {
		return mutedStyle.Render(fmt.Sprintf("rendering %d messages...", len(m.conversation.Messages)))
	}
	lines := m.detail.slice(m.detailScroll, m.bodyHeight())
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}

func (m Model) messageTotal() int {
	if count := m.conversation.Session.MessageCount; count > 0 {
		return count
	}
	return len(m.conversation.Messages)
}

func (m Model) messageCountLabel() string {
	if m.hidden == 0 {
		return fmt.Sprintf("%d msgs", m.messageTotal())
	}
	return fmt.Sprintf("%d of %d msgs", m.detail.messages, m.messageTotal())
}

func (m Model) typesBody() string {
	label := "show message types"
	count := m.messageCountLabel()
	gap := m.width - 2 - len(label) - len(count)
	if gap < 1 {
		gap = 1
	}

	lines := make([]string, 0, len(typeRows)+2)
	lines = append(lines, " "+titleStyle.Render(label)+strings.Repeat(" ", gap)+mutedStyle.Render(count))
	lines = append(lines, "")

	for i, row := range typeRows {
		marker := "  "
		if i == m.typesCursor {
			marker = "> "
		}
		box := "[ ] "
		if m.typesDraft.shows(row.category) {
			box = "[x] "
		}
		name := padRight(row.label, 12)
		if i == m.typesCursor {
			name = selectedStyle.Render(name)
		}
		lines = append(lines, " "+marker+box+name+
			padLeft(strconv.Itoa(m.typesCounts[row.category]), 5)+
			"   "+mutedStyle.Render(row.note))
	}
	return strings.Join(lines, "\n")
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
	case modeDetail, modeTypes:
		if m.rendering || !m.detail.built {
			return mutedStyle.Render(fmt.Sprintf(
				"%s · %d messages · rendering",
				shortID(m.conversation.Session.AtuinID),
				len(m.conversation.Messages),
			))
		}
		count := m.messageCountLabel()
		if m.mode == modeTypes {
			return mutedStyle.Render(shortID(m.conversation.Session.AtuinID) + " · " + count)
		}
		return mutedStyle.Render(fmt.Sprintf(
			"%s · %s · line %d of %d",
			shortID(m.conversation.Session.AtuinID),
			count,
			m.detailScroll+1,
			m.detail.total,
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
		return mutedStyle.Render("j/k scroll · pgup/pgdn page · g/G top/bottom · t types · esc back · q quit")
	case modeTypes:
		return mutedStyle.Render("space toggle · a all · n conversation · 1-4 preset · enter apply · esc cancel")
	case modeFilter:
		return mutedStyle.Render("type to filter · enter apply · ctrl+u reset · esc clear")
	default:
		return mutedStyle.Render("j/k move · enter read · / filter · ctrl+s full-text · ctrl+r reload · q quit")
	}
}

type listLayout struct {
	width int
	id    int
	date  int
	agent int
	msgs  int
	title int
	dir   int
}

type listMetrics struct {
	msgs  int
	agent int
	title int
	dir   int
}

func measureThreads(threads []atuin.Thread) listMetrics {
	var m listMetrics
	for _, t := range threads {
		m.msgs = max(m.msgs, t.MessageCount)
		m.agent = max(m.agent, ansi.StringWidth(agentTag(t.Harness)))
		m.title = max(m.title, ansi.StringWidth(titleLabel(t)))
		m.dir = max(m.dir, ansi.StringWidth(dirLabel(t.CWD)))
	}
	return m
}

func newListLayout(width int, m listMetrics) listLayout {
	l := listLayout{
		width: width,
		id:    idWidth,
		date:  dateWidth,
		agent: clamp(m.agent, agentMinWidth, agentMaxWidth),
		msgs:  max(msgsWidth, len(strconv.Itoa(m.msgs))),
		title: clamp(m.title, titleMin, titleMax),
		dir:   clamp(m.dir, dirMin, dirMax),
	}

	for l.total() > width && l.title > titleMin {
		l.title--
	}
	for l.total() > width && l.dir > dirMin {
		l.dir--
	}
	for _, column := range []*int{&l.dir, &l.date, &l.agent, &l.msgs} {
		if l.total() <= width {
			break
		}
		*column = 0
	}
	l.title = max(1, min(l.title, width-l.fixed()))
	return l
}

func (l listLayout) fixed() int {
	width := markerWidth + l.id
	for _, column := range []int{l.date, l.agent, l.msgs} {
		if column > 0 {
			width += listGutter + column
		}
	}
	return width
}

func (l listLayout) total() int {
	width := l.fixed() + listGutter + l.title
	if l.dir > 0 {
		width += listGutter + l.dir
	}
	return width
}

func (l listLayout) header() string {
	return l.line("  ", "id", "date", agentHeader(l.agent), "msgs", "title", "dir")
}

func (l listLayout) rule() string {
	rule := func(n int) string { return strings.Repeat("─", n) }
	return l.line("  ", rule(l.id), rule(l.date), rule(l.agent), rule(l.msgs), rule(l.title), rule(l.dir))
}

func (l listLayout) row(t atuin.Thread, selected bool) string {
	marker := "  "
	if selected {
		marker = "> "
	}

	line := l.line(
		marker,
		shortID(t.AtuinID),
		dateLabel(t),
		agentTag(t.Harness),
		strconv.Itoa(t.MessageCount),
		titleLabel(t),
		tailTruncate(dirLabel(t.CWD), l.dir),
	)
	if selected {
		return selectedStyle.Render(line)
	}
	return line
}

func (l listLayout) line(marker, id, date, agent, msgs, title, dir string) string {
	gutter := strings.Repeat(" ", listGutter)
	var b strings.Builder
	b.WriteString(marker)
	b.WriteString(padRight(id, l.id))
	if l.date > 0 {
		b.WriteString(gutter)
		b.WriteString(padRight(date, l.date))
	}
	if l.agent > 0 {
		b.WriteString(gutter)
		b.WriteString(padRight(agent, l.agent))
	}
	if l.msgs > 0 {
		b.WriteString(gutter)
		b.WriteString(padLeft(msgs, l.msgs))
	}
	b.WriteString(gutter)
	b.WriteString(padRight(title, l.title))
	if l.dir > 0 {
		b.WriteString(gutter)
		b.WriteString(padRight(dir, l.dir))
	}
	return ansi.Truncate(b.String(), l.width, "")
}

func agentHeader(width int) string {
	if width >= agentMaxWidth {
		return "agent"
	}
	return "ag"
}

func clamp(v, low, high int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func dateLabel(t atuin.Thread) string {
	updated := t.UpdatedAt
	if updated.IsZero() {
		return "-"
	}
	return updated.Local().Format("02 Jan 15:04")
}

func agentTag(harness string) string {
	name := strings.ToLower(strings.TrimSpace(harness))
	switch name {
	case "":
		return "--"
	case "claude-code":
		return "claude"
	case "opencode":
		return "ocode"
	case "codex":
		return "codex"
	case "pi":
		return "pi"
	}
	runes := []rune(name)
	if len(runes) > agentMaxWidth {
		return string(runes[:agentMaxWidth])
	}
	return name
}

func titleLabel(t atuin.Thread) string {
	title := strings.TrimSpace(t.Title)
	if title == "" {
		return "(untitled)"
	}
	return title
}

func dirLabel(cwd string) string {
	cwd = strings.TrimSpace(cwd)
	if cwd == "" {
		return "-"
	}
	cwd = strings.TrimRight(cwd, "/")
	if cwd == "" {
		return "/"
	}
	return cwd[strings.LastIndex(cwd, "/")+1:]
}

func padRight(s string, width int) string {
	s = ansi.Truncate(s, width, "…")
	if pad := width - ansi.StringWidth(s); pad > 0 {
		s += strings.Repeat(" ", pad)
	}
	return s
}

func padLeft(s string, width int) string {
	s = ansi.Truncate(s, width, "…")
	if pad := width - ansi.StringWidth(s); pad > 0 {
		s = strings.Repeat(" ", pad) + s
	}
	return s
}

func tailTruncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	w := ansi.StringWidth(s)
	if w <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return ansi.TruncateLeft(s, w-width+1, "…")
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if ansi.StringWidth(s) <= width {
		return s
	}
	return ansi.Truncate(s, width, "…")
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
