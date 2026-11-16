package ui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

const (
	detailProseMax  = 80
	detailMetaWidth = 24
	detailGap       = 2
	detailMargin    = 1
	detailBodyMin   = 20
	maxToolPayload  = 4
	maxMessageParts = 2

	// detailSyncMessages is the largest conversation rendered inside Update.
	detailSyncMessages = 500
)

type detailLayout struct {
	body  int
	prose int
	meta  int
}

func newDetailLayout(width int) detailLayout {
	if width < detailBodyMin+detailGap+detailMetaWidth+2*detailMargin {
		body := max(width-2*detailMargin, 1)
		return detailLayout{body: body, prose: min(detailProseMax, body)}
	}
	body := width - 2*detailMargin - detailGap - detailMetaWidth
	return detailLayout{body: body, prose: min(detailProseMax, body), meta: detailMetaWidth}
}

type messageView struct {
	body []string
	meta [2]string
}

// detailBlock is one message, or one collapsed run, and the lines it occupies.
type detailBlock struct {
	start int
	msg   int
	run   bool
	lines []string
}

// detailRendering is a conversation rendered for one terminal width and one
// message type selection.
type detailRendering struct {
	width    int
	hidden   typeMask
	built    bool
	head     []string
	blocks   []detailBlock
	total    int
	messages int
}

func (r detailRendering) add(msg int, lines []string) detailRendering {
	return r.addBlock(msg, false, lines)
}

// addRun adds the block of a collapsed run of empty messages, which carries the
// index of the first message in the run but cannot be selected on its own.
func (r detailRendering) addRun(msg int, lines []string) detailRendering {
	return r.addBlock(msg, true, lines)
}

func (r detailRendering) addBlock(msg int, run bool, lines []string) detailRendering {
	if len(lines) == 0 {
		return r
	}
	r.blocks = append(r.blocks, detailBlock{start: r.total, msg: msg, run: run, lines: lines})
	r.total += len(lines)
	return r
}

// selectable returns the message indices that can be selected, in the order
// they appear. Collapsed runs of empty messages are not among them.
func (r detailRendering) selectable() []int {
	out := make([]int, 0, len(r.blocks))
	for _, block := range r.blocks {
		if !block.run {
			out = append(out, block.msg)
		}
	}
	return out
}

// lineInfo reports which message draws a rendered line and whether that line is
// the first one of its message.
func (r detailRendering) lineInfo(line int) (msg int, first bool) {
	if !r.built || line < 0 || line >= r.total {
		return -1, false
	}
	found := sort.Search(len(r.blocks), func(i int) bool {
		return r.blocks[i].start+len(r.blocks[i].lines) > line
	})
	if found >= len(r.blocks) || r.blocks[found].start > line {
		return -1, false
	}
	block := r.blocks[found]
	return block.msg, block.start == line
}

// messageSpan returns the first and the last line of a message block.
func (r detailRendering) messageSpan(msg int) (start, end int, ok bool) {
	found := sort.Search(len(r.blocks), func(i int) bool {
		return r.blocks[i].msg >= msg
	})
	if found >= len(r.blocks) || r.blocks[found].msg != msg || r.blocks[found].run {
		return 0, 0, false
	}
	block := r.blocks[found]
	return block.start, block.start + len(block.lines) - 1, true
}

// messageAt returns the index of the message at a rendered line, or -1 when
// that line belongs to the thread header.
func (r detailRendering) messageAt(line int) int {
	if !r.built || line < 0 || line >= r.total {
		return -1
	}
	found := sort.Search(len(r.blocks), func(i int) bool {
		return r.blocks[i].start+len(r.blocks[i].lines) > line
	})
	if found >= len(r.blocks) || r.blocks[found].start > line {
		return -1
	}
	return r.blocks[found].msg
}

// scrollToMessage returns the line of the first block that starts at the given
// message or at a later one.
func (r detailRendering) scrollToMessage(msg int) int {
	found := sort.Search(len(r.blocks), func(i int) bool {
		return r.blocks[i].msg >= msg
	})
	if found >= len(r.blocks) {
		return 0
	}
	return r.blocks[found].start
}

// slice returns the lines between from and from+count, reading only the blocks
// that overlap that range.
func (r detailRendering) slice(from, count int) []string {
	if !r.built || count <= 0 || from < 0 || from >= r.total {
		return nil
	}
	end := from + count
	if end > r.total {
		end = r.total
	}

	out := make([]string, 0, end-from)
	if from < len(r.head) {
		stop := min(len(r.head), end)
		out = append(out, r.head[from:stop]...)
	}

	first := sort.Search(len(r.blocks), func(i int) bool {
		block := r.blocks[i]
		return block.start+len(block.lines) > from
	})
	for i := first; i < len(r.blocks) && len(out) < end-from; i++ {
		block := r.blocks[i]
		lo := max(0, from-block.start)
		hi := min(len(block.lines), end-block.start)
		if hi <= lo {
			continue
		}
		out = append(out, block.lines[lo:hi]...)
	}
	return out
}

// all joins every rendered line, which is only needed by tests and callers
// outside the draw path.
func (r detailRendering) all() []string {
	if !r.built {
		return nil
	}
	out := make([]string, 0, r.total)
	out = append(out, r.head...)
	for _, block := range r.blocks {
		out = append(out, block.lines...)
	}
	return out
}

func (m Model) detailLines() []string {
	return m.buildDetail().all()
}

func (m Model) buildDetail() detailRendering {
	if m.mode != modeDetail {
		return detailRendering{}
	}

	layout := newDetailLayout(m.width)
	session := m.conversation.Session
	now := m.now
	if now == nil {
		now = time.Now
	}

	r := detailRendering{width: m.width, hidden: m.hidden, built: true}

	title := strings.TrimSpace(session.Title)
	if title == "" {
		title = "(untitled)"
	}
	r.head = append(r.head, titleStyle.Render(" "+truncate(title, layout.body)))

	var header []string
	if session.Harness != "" {
		header = append(header, session.Harness)
	}
	if session.Model != "" {
		header = append(header, session.Model)
	}
	if session.AtuinID != "" {
		header = append(header, shortID(session.AtuinID))
	}
	if session.MessageCount > 0 {
		header = append(header, fmt.Sprintf("%d msgs", session.MessageCount))
	}
	if !session.UpdatedAt.IsZero() {
		header = append(header, age(now(), session.UpdatedAt)+" ago")
	}
	if len(header) > 0 {
		r.head = append(r.head, mutedStyle.Render(" "+truncate(strings.Join(header, " · "), layout.body)))
	}
	if session.CWD != "" {
		r.head = append(r.head, mutedStyle.Render(" "+truncate(session.CWD, layout.body)))
	}
	r.head = append(r.head, "")
	r.total = len(r.head)

	views := make([]messageView, len(m.conversation.Messages))
	for i, msg := range m.conversation.Messages {
		views[i] = newMessageView(msg, layout.prose, layout.body, m.hidden)
	}

	for i := 0; i < len(m.conversation.Messages); {
		msg := m.conversation.Messages[i]
		if len(msg.Content) == 0 {
			if !m.hidden.shows(typeNotes) {
				for i < len(m.conversation.Messages) && len(m.conversation.Messages[i].Content) == 0 {
					i++
				}
				continue
			}
			end := i + 1
			for end < len(views) && len(m.conversation.Messages[end].Content) == 0 {
				end++
			}
			if end-i >= 2 {
				lines := append(layout.rows(collapsedView(m.conversation.Messages[i:end], views[i])), "")
				r = r.addRun(i, lines)
				i = end
				continue
			}
		}
		if len(msg.Content) > 0 && !hasVisibleBlock(msg, m.hidden) {
			i++
			continue
		}
		r.messages++
		r = r.add(i, append(layout.rows(views[i]), ""))
		i++
	}

	if r.messages == 0 && len(m.conversation.Messages) > 0 {
		r.head = append(r.head, mutedStyle.Render(" no messages of the selected types"))
		r.total = len(r.head)
	}

	return r
}

func (l detailLayout) rows(view messageView) []string {
	total := detailMargin + l.body
	if l.meta > 0 {
		total += detailGap + l.meta
	}

	count := len(view.body)
	if l.meta > 0 && count < len(view.meta) {
		count = len(view.meta)
	}

	out := make([]string, count)
	for i := 0; i < count; i++ {
		var body string
		if i < len(view.body) {
			body = view.body[i]
		}
		if l.meta == 0 {
			out[i] = strings.TrimRight(strings.Repeat(" ", detailMargin)+truncate(body, l.body), " ")
			continue
		}
		var meta string
		if i < len(view.meta) {
			meta = view.meta[i]
		}
		if meta == "" {
			out[i] = strings.TrimRight(strings.Repeat(" ", detailMargin)+truncate(body, l.body), " ")
			continue
		}
		line := strings.Repeat(" ", detailMargin) + padRight(body, l.body) +
			strings.Repeat(" ", detailGap) + meta
		out[i] = ansi.Truncate(line, total, "…")
	}
	return out
}

func newMessageView(msg atuin.Message, prose, code int, hidden typeMask) messageView {
	body := messageBody(msg, prose, code, hidden)
	if onlyBlank(body) {
		body = nil
	}

	return messageView{
		body: body,
		meta: [2]string{
			mutedStyle.Render(clock(msg.Timestamp)+" - ") +
				authorStyle(msg.Role).Render(authorLabel(msg.Role)),
			mutedStyle.Render(secondMeta(msg, hidden)),
		},
	}
}

func collapsedView(messages []atuin.Message, first messageView) messageView {
	label := fmt.Sprintf("(%d empty messages)", len(messages))
	attachments := 0
	for _, msg := range messages {
		if strings.EqualFold(strings.TrimSpace(msg.Role), "attachment") {
			attachments++
		}
	}
	if attachments == len(messages) {
		label = fmt.Sprintf("(%d attachments)", attachments)
	}

	return messageView{
		body: []string{mutedStyle.Render(label)},
		meta: first.meta,
	}
}

func messageBody(msg atuin.Message, prose, code int, hidden typeMask) []string {
	style := roleProseStyle(msg.Role)
	var out []string
	for _, block := range visibleBlocks(msg, hidden) {
		out = append(out, blockLines(block, prose, code, style)...)
	}
	return out
}

func roleProseStyle(role string) func(...string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "user":
		return promptStyle.Render
	case "assistant":
		return replyStyle.Render
	}
	return nil
}

func styled(style func(...string) string, line string) string {
	if style == nil {
		return line
	}
	return style(line)
}

func blockLines(block atuin.Content, prose, code int, proseStyle func(...string) string) []string {
	if prose < 1 {
		prose = 1
	}
	if code < 1 {
		code = 1
	}

	switch block.Type {
	case "text":
		return styleAll(wrap(block.Text, prose), proseStyle)

	case "thinking":
		if placeholderThinking(block.Text) {
			return nil
		}
		return styleAll(wrap(block.Text, prose), mutedStyle.Render)

	case "tool_call":
		return toolCallLines(block, code)

	case "tool_result":
		return resultLines(block, code)

	case "patch":
		return patchLines(block, code)

	case "summary":
		return wrap(block.Text, prose)
	}

	name := strings.TrimSpace(block.Type)
	if name == "" {
		return nil
	}
	return []string{mutedStyle.Render("(" + name + ")")}
}

func toolCallLines(block atuin.Content, width int) []string {
	input := parseInput(block.Input)
	name := strings.TrimSpace(block.Name)
	if name == "" {
		name = "tool"
	}

	if input == nil {
		lines := []string{truncate(name, width)}
		if raw := strings.TrimSpace(block.Input); raw != "" {
			for _, part := range preserve(raw, width-2) {
				lines = append(lines, "  "+part)
			}
		}
		return capped(lines)
	}

	switch strings.ToLower(name) {
	case "bash", "shell":
		return shellLines(name, input, width)
	case "write":
		return writeLines(name, input, width)
	case "edit":
		return editLines(name, input, width)
	case "read":
		return readLines(name, input, width)
	case "skill":
		return skillLines(name, input, width)
	case "compress":
		return compressLines(name, input, width)
	}

	return genericLines(name, input, width)
}

func shellLines(name string, input map[string]any, width int) []string {
	head := titleStyle.Render(name)
	if description := field(input, "description"); description != "" {
		head += "  " + description
	}
	lines := []string{truncate(head, width)}

	command := strings.TrimSpace(field(input, "command"))
	if command == "" {
		return lines
	}

	parts := preserveLines(command, width-4)
	if len(parts) > maxToolPayload {
		parts = parts[:maxToolPayload]
	}
	for i, part := range parts {
		if i == 0 {
			lines = append(lines, "  $ "+commandStyle.Render(part))
			continue
		}
		lines = append(lines, "    "+commandStyle.Render(part))
	}
	return capped(lines)
}

func headSuffix(head, suffix string, width int) string {
	if room := width - ansi.StringWidth(suffix); room >= 4 {
		return truncate(head, room) + suffix
	}
	return truncate(head, width)
}

func countsSuffix(additions, deletions int) string {
	return "  " + additionStyle.Render(fmt.Sprintf("+%d", additions)) +
		" " + deletionStyle.Render(fmt.Sprintf("-%d", deletions))
}

func writeLines(name string, input map[string]any, width int) []string {
	path := firstField(input, "file_path", "path")
	content := field(input, "content")
	head := titleStyle.Render(name)
	if path != "" {
		head += " " + pathStyle.Render(path)
	}
	if count := countLines(content); count > 0 {
		head = headSuffix(head, "  "+additionStyle.Render(fmt.Sprintf("+%d", count))+" lines", width)
	} else {
		head = truncate(head, width)
	}
	lines := []string{head}

	for i, part := range preserve(content, width-2) {
		if i >= 3 {
			break
		}
		lines = append(lines, "  "+part)
	}
	return capped(lines)
}

func editLines(name string, input map[string]any, width int) []string {
	path := firstField(input, "file_path", "path")
	oldText := strings.TrimSpace(firstField(input, "old_string", "oldString"))
	newText := strings.TrimSpace(firstField(input, "new_string", "newString"))

	head := titleStyle.Render(name)
	if path != "" {
		head += " " + pathStyle.Render(path)
	}
	if oldText != "" || newText != "" {
		head = headSuffix(head, countsSuffix(countLines(newText), countLines(oldText)), width)
	} else {
		head = truncate(head, width)
	}
	lines := []string{head}

	lines = append(lines, prefixedLines(oldText, "- ", width, deletionStyle.Render)...)
	lines = append(lines, prefixedLines(newText, "+ ", width, additionStyle.Render)...)
	return capped(lines)
}

func prefixedLines(text, prefix string, width int, style func(...string) string) []string {
	if text == "" {
		return nil
	}
	parts := preserve(text, width-4)
	if len(parts) > maxMessageParts {
		parts = parts[:maxMessageParts]
	}
	var out []string
	for i, part := range parts {
		if i == 0 {
			out = append(out, "  "+prefix+styled(style, part))
			continue
		}
		out = append(out, "    "+styled(style, part))
	}
	return out
}

func readLines(name string, input map[string]any, width int) []string {
	path := firstField(input, "file_path", "path")
	head := titleStyle.Render(name)
	if path != "" {
		head += " " + pathStyle.Render(path)
	}
	offset := intField(input, "offset")
	limit := intField(input, "limit")
	switch {
	case offset > 0 && limit > 0:
		head += fmt.Sprintf("  %d..%d", offset, offset+limit)
	case offset > 0:
		head += fmt.Sprintf("  from %d", offset)
	}
	return []string{truncate(head, width)}
}

func skillLines(name string, input map[string]any, width int) []string {
	head := titleStyle.Render(name)
	if skill := field(input, "skill"); skill != "" {
		head += " " + pathStyle.Render(skill)
	}
	line := truncate(head, width)
	if args := field(input, "args"); args != "" {
		room := width - ansi.StringWidth(line) - 2
		if room > 0 {
			line += "  " + mutedStyle.Render(truncate(args, room))
		}
	}
	return []string{line}
}

func compressLines(name string, input map[string]any, width int) []string {
	head := titleStyle.Render(name)
	if topic := field(input, "topic"); topic != "" {
		head += " " + topic
	}
	return []string{truncate(head, width)}
}

func genericLines(name string, input map[string]any, width int) []string {
	pairs := make([]string, 0, len(input))
	for _, key := range sortedKeys(input) {
		value := strings.Join(strings.Fields(field(input, key)), " ")
		if value == "" {
			continue
		}
		pairs = append(pairs, key+"="+truncate(value, 40))
	}

	line := titleStyle.Render(name)
	used := 0
	for i, pair := range pairs {
		candidate := line + "  " + pair
		if ansi.StringWidth(candidate) > width {
			rest := len(pairs) - i
			line += fmt.Sprintf("  +%d more", rest)
			break
		}
		line = candidate
		used++
	}
	if used == 0 && len(pairs) > 0 && ansi.StringWidth(line) <= width {
		line += fmt.Sprintf("  +%d more", len(pairs))
	}
	return []string{truncate(line, width)}
}

func resultLines(block atuin.Content, width int) []string {
	content := strings.TrimSpace(block.Content)
	if content == "" {
		return nil
	}

	lines := capped(preserve(content, width))
	if !block.IsError {
		return lines
	}

	lines = styleAll(lines, errorStyle.Render)
	lines[0] = errorStyle.Render("! ") + lines[0]
	return lines
}

func patchLines(block atuin.Content, width int) []string {
	if len(block.Files) == 0 {
		return nil
	}

	var lines []string
	for _, file := range block.Files {
		head := pathStyle.Render(file.Path)
		if file.Change != "" {
			head += "  " + file.Change
		}
		if file.Additions > 0 || file.Deletions > 0 {
			head = headSuffix(head, countsSuffix(file.Additions, file.Deletions), width)
		} else {
			head = truncate(head, width)
		}
		lines = append(lines, head)

		diff := strings.TrimSpace(file.Diff)
		if diff == "" {
			continue
		}
		for _, part := range preserve(diff, width-2) {
			lines = append(lines, "  "+styleDiffLine(part, mutedStyle.Render))
			if len(lines) > maxBlockLines {
				break
			}
		}
		if len(lines) > maxBlockLines {
			break
		}
	}
	return capped(lines)
}

func styleDiffLine(line string, header func(...string) string) string {
	switch {
	case strings.HasPrefix(line, "@@"):
		return hunkStyle.Render(line)
	case strings.HasPrefix(line, "+++"), strings.HasPrefix(line, "---"):
		return styled(header, line)
	case strings.HasPrefix(line, "+"):
		return additionStyle.Render(line)
	case strings.HasPrefix(line, "-"):
		return deletionStyle.Render(line)
	}
	return line
}

func secondMeta(msg atuin.Message, hidden typeMask) string {
	kinds := messageKinds(msg, hidden)
	label := strings.Join(kinds, "+")
	if label == "" {
		label = "empty"
	}
	if msg.Tokens == nil {
		return label
	}

	tokens := msg.Tokens.Output
	if len(kinds) == 1 && kinds[0] == "think" {
		tokens = msg.Tokens.Reasoning
	}
	return fmt.Sprintf("%s - %d tokens", label, tokens)
}

func messageKinds(msg atuin.Message, hidden typeMask) []string {
	seen := make(map[string]bool, len(msg.Content))
	var kinds []string
	for _, block := range msg.Content {
		if !hidden.shows(categoryOf(block, msg.Role)) {
			continue
		}
		kind := kindName(block.Type)
		if kind == "" || seen[kind] {
			continue
		}
		seen[kind] = true
		kinds = append(kinds, kind)
	}
	return kinds
}

func kindName(blockType string) string {
	switch blockType {
	case "text":
		return "text"
	case "thinking":
		return "think"
	case "tool_call":
		return "tool"
	case "tool_result":
		return "out"
	case "patch":
		return "patch"
	case "summary":
		return "summ"
	case "error":
		return "err"
	}
	return truncate(strings.TrimSpace(blockType), 5)
}

func authorLabel(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "":
		return "?"
	case "user":
		return "user"
	case "assistant":
		return "asst"
	case "tool":
		return "tool"
	case "system":
		return "sys"
	case "attachment":
		return "att"
	case "ai-title":
		return "title"
	case "away_summary", "away-summary":
		return "away"
	}
	name := []rune(strings.ToLower(strings.TrimSpace(role)))
	if len(name) > 5 {
		name = name[:5]
	}
	return string(name)
}

func authorStyle(role string) lipgloss.Style {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "user":
		return userStyle
	case "assistant":
		return assistantStyle
	}
	return systemStyle
}

func clock(t time.Time) string {
	if t.IsZero() {
		return "--:--"
	}
	return t.Format("15:04")
}

func onlyBlank(lines []string) bool {
	for _, line := range lines {
		if strings.TrimSpace(ansi.Strip(line)) != "" {
			return false
		}
	}
	return true
}

func styleAll(lines []string, style func(...string) string) []string {
	if style == nil {
		return lines
	}
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = style(line)
	}
	return out
}

func capped(lines []string) []string {
	if len(lines) <= maxBlockLines {
		return lines
	}
	out := append([]string{}, lines[:maxBlockLines]...)
	return append(out, "  ... ("+fmt.Sprintf("%d", len(lines)-maxBlockLines)+" more lines)")
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
		used := 0
		for _, word := range words {
			for _, chunk := range splitWidth(word, width) {
				size := ansi.StringWidth(chunk)
				switch {
				case line == "":
					line, used = chunk, size
				case used+1+size <= width:
					line += " " + chunk
					used += 1 + size
				default:
					out = append(out, line)
					line, used = chunk, size
				}
			}
		}
		out = append(out, line)
	}
	return out
}

func preserveLines(text string, width int) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		line = strings.TrimRight(expandTabs(line), " ")
		if line == "" {
			out = append(out, "")
			continue
		}
		out = append(out, splitWidth(line, width)...)
	}
	return out
}

func preserve(text string, width int) []string {
	if width < 2 {
		width = 2
	}

	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		line = strings.TrimRight(expandTabs(line), " ")
		if line == "" {
			out = append(out, "")
			continue
		}
		if ansi.StringWidth(line) <= width {
			out = append(out, line)
			continue
		}
		out = append(out, wrapPreserving(line, width)...)
	}
	return out
}

func wrapPreserving(line string, width int) []string {
	var out []string
	for rest := line; rest != ""; {
		limit := width
		prefix := ""
		if len(out) > 0 {
			prefix = "  "
			limit = width - 2
		}
		head, next := cutWidth(rest, limit)
		out = append(out, prefix+head)
		rest = next
	}
	return out
}

func cutWidth(s string, width int) (string, string) {
	if width < 1 {
		width = 1
	}
	used := 0
	for i, r := range s {
		size := runewidth.RuneWidth(r)
		if used > 0 && used+size > width {
			return s[:i], s[i:]
		}
		used += size
	}
	return s, ""
}

func placeholderThinking(text string) bool {
	trimmed := strings.TrimSpace(text)
	if !strings.HasPrefix(trimmed, "Reasoning") || !strings.HasSuffix(trimmed, "tokens") {
		return false
	}
	digits := strings.TrimSuffix(strings.TrimPrefix(trimmed, "Reasoning · "), " tokens")
	if digits == trimmed {
		return false
	}
	if digits == "" {
		return false
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func expandTabs(line string) string {
	if !strings.ContainsRune(line, '\t') {
		return line
	}

	var b strings.Builder
	width := 0
	for _, r := range line {
		if r != '\t' {
			b.WriteRune(r)
			width += runewidth.RuneWidth(r)
			continue
		}
		spaces := 4 - width%4
		b.WriteString(strings.Repeat(" ", spaces))
		width += spaces
	}
	return b.String()
}

func splitWidth(s string, width int) []string {
	if width < 1 {
		width = 1
	}

	var out []string
	for rest := s; rest != ""; {
		head, next := cutWidth(rest, width)
		out = append(out, head)
		rest = next
	}
	return out
}

func parseInput(raw string) map[string]any {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed[0] != '{' {
		return nil
	}

	var input map[string]any
	if err := json.Unmarshal([]byte(trimmed), &input); err != nil {
		return nil
	}
	return input
}

func field(input map[string]any, key string) string {
	value, ok := input[key]
	if !ok {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%g", typed)
	case bool:
		return fmt.Sprintf("%t", typed)
	case nil:
		return ""
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(encoded)
}

func firstField(input map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := field(input, key); value != "" {
			return value
		}
	}
	return ""
}

func intField(input map[string]any, key string) int {
	value, ok := input[key].(float64)
	if !ok {
		return 0
	}
	return int(value)
}

func countLines(text string) int {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(strings.ReplaceAll(trimmed, "\r\n", "\n"), "\n"))
}

func sortedKeys(input map[string]any) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
