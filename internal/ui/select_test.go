package ui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func stamp(hour, minute int) time.Time {
	return time.Date(2026, 11, 16, hour, minute, 0, 0, time.UTC)
}

func TestExportMarkdownHeadings(t *testing.T) {
	conv := atuin.Conversation{Messages: []atuin.Message{
		{Role: "user", Timestamp: stamp(11, 53), Content: []atuin.Content{{Type: "text", Text: "write the script"}}},
		{Role: "assistant", Timestamp: stamp(11, 54), Content: []atuin.Content{{Type: "text", Text: "done"}}},
	}}

	got := exportMarkdown(conv.Messages, 0)
	if !strings.Contains(got, "### 11:53 user") {
		t.Fatalf("missing user heading in\n%s", got)
	}
	if !strings.Contains(got, "### 11:54 asst") {
		t.Fatalf("missing assistant heading in\n%s", got)
	}
	if !strings.Contains(got, "write the script\n\n### 11:54 asst\n\ndone") {
		t.Fatalf("bodies are not separated by a blank line:\n%s", got)
	}
}

func TestExportSingleReplyHasNoHeading(t *testing.T) {
	msg := atuin.Message{Role: "assistant", Timestamp: stamp(11, 54),
		Content: []atuin.Content{{Type: "text", Text: "just this"}}}

	if got := exportMarkdown([]atuin.Message{msg}, 0); got != "just this\n\n" && got != "just this\n" {
		t.Fatalf("payload = %q, want the prose alone", got)
	}
}

func TestExportFollowsTheTypeFilter(t *testing.T) {
	msg := atuin.Message{Role: "assistant", Timestamp: stamp(11, 54), Content: []atuin.Content{
		{Type: "thinking", Text: "secret reasoning"},
		{Type: "text", Text: "visible reply"},
	}}

	if got := exportMarkdown([]atuin.Message{msg}, typeThinking); strings.Contains(got, "secret reasoning") {
		t.Fatalf("hidden thinking leaked into\n%s", got)
	}
	if got := exportMarkdown([]atuin.Message{msg}, typeThinking); !strings.Contains(got, "visible reply") {
		t.Fatalf("visible reply missing from\n%s", got)
	}
}

func TestExportToolCallKeepsOneLineCommand(t *testing.T) {
	long := strings.Repeat("segment ", 40)
	msg := atuin.Message{Role: "assistant", Timestamp: stamp(12, 0), Content: []atuin.Content{
		{Type: "tool_call", Name: "shell", Input: `{"command":"` + long + `","description":"run it"}`},
	}}

	got := exportMarkdown([]atuin.Message{msg}, 0)
	if !strings.Contains(got, "**shell** - run it") {
		t.Fatalf("missing head in\n%s", got)
	}
	if !strings.Contains(got, "$ "+long) {
		t.Fatal("the command is not on one line")
	}
	if strings.Count(got, "\n") < 4 {
		t.Fatalf("the command was not fenced:\n%s", got)
	}
}

func TestExportFenceGrowsPastContent(t *testing.T) {
	msg := atuin.Message{Role: "assistant", Timestamp: stamp(12, 0), Content: []atuin.Content{
		{Type: "tool_result", Content: "before\n```\ninside\n```\nafter"},
	}}

	got := exportMarkdown([]atuin.Message{msg}, 0)
	if !strings.Contains(got, "````") {
		t.Fatalf("fence did not grow:\n%s", got)
	}
	closed := 0
	for _, line := range strings.Split(got, "\n") {
		if line == "````" {
			closed++
		}
	}
	if closed != 2 {
		t.Fatalf("fence delimiters = %d, want an opening and a closing one:\n%s", closed, got)
	}
}

func TestExportPatchShowsCounts(t *testing.T) {
	msg := atuin.Message{Role: "assistant", Timestamp: stamp(12, 0), Content: []atuin.Content{
		{Type: "patch", Files: []atuin.PatchFile{{
			Path:      "flake.nix",
			Change:    "update",
			Additions: 2,
			Deletions: 1,
			Diff:      "@@ -1,1 +1,2 @@\n-old\n+new\n+more",
		}}},
	}}

	got := exportMarkdown([]atuin.Message{msg}, 0)
	if !strings.Contains(got, "**patch** flake.nix (+2 -1)") {
		t.Fatalf("missing patch head in\n%s", got)
	}
	if !strings.Contains(got, "-old") || !strings.Contains(got, "+new") {
		t.Fatalf("missing diff lines in\n%s", got)
	}
}

func TestExportErrorMarksTheFirstLine(t *testing.T) {
	msg := atuin.Message{Role: "tool", Timestamp: stamp(12, 0), Content: []atuin.Content{
		{Type: "tool_result", Content: "boom\nmore", IsError: true},
	}}

	got := exportMarkdown([]atuin.Message{msg}, 0)
	if !strings.Contains(got, "! boom") {
		t.Fatalf("error is not marked in\n%s", got)
	}
}

func TestExportNameUsesOneBasedMessages(t *testing.T) {
	if got := exportName("fd025611", 11, 13); got != "att-fd025611-12-14.md" {
		t.Fatalf("name = %q", got)
	}
}

func TestWriteExportRefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	if _, err := writeExport(dir, "att-x-1-2.md", "one"); err != nil {
		t.Fatal(err)
	}
	if _, err := writeExport(dir, "att-x-1-2.md", "two"); err == nil {
		t.Fatal("the second write replaced the file")
	}
	data, err := os.ReadFile(filepath.Join(dir, "att-x-1-2.md"))
	if err != nil || string(data) != "one" {
		t.Fatalf("file = %q, %v", data, err)
	}
}

func selectMode(t *testing.T, conv atuin.Conversation) Model {
	t.Helper()
	m, _ := openDetail(t, conv, 120, 24)
	m, _ = pressRune(t, m, 'v')
	return m
}

func TestSelectEntersOnTheTopMessage(t *testing.T) {
	m := selectMode(t, sampleConversation())
	if m.mode != modeSelect {
		t.Fatal("v did not enter select mode")
	}
	if len(m.selection) != 2 {
		t.Fatalf("selectable = %v, want two messages", m.selection)
	}
	if m.selectCursor != 0 {
		t.Fatalf("cursor = %d, want the top message", m.selectCursor)
	}
}

func TestSelectMovesByMessage(t *testing.T) {
	m := selectMode(t, longConversation())
	first := m.selection[m.selectCursor]

	m, _ = pressRune(t, m, 'j')
	if m.selection[m.selectCursor] != first+1 {
		t.Fatalf("j did not move one message: %d", m.selection[m.selectCursor])
	}

	m, _ = pressRune(t, m, 'J')
	if m.selection[m.selectCursor] != first+11 {
		t.Fatalf("J did not move ten messages: %d", m.selection[m.selectCursor])
	}

	m, _ = pressRune(t, m, 'G')
	if m.selectCursor != len(m.selection)-1 {
		t.Fatalf("G did not reach the last message: %d", m.selectCursor)
	}

	m, _ = pressRune(t, m, 'g')
	if m.selectCursor != 0 {
		t.Fatalf("g did not reach the first message: %d", m.selectCursor)
	}
}

func TestSpaceMarksMessages(t *testing.T) {
	m := selectMode(t, longConversation())
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')
	m, _ = pressRune(t, m, ' ')

	if got := m.markedCount(); got != 3 {
		t.Fatalf("marked = %d, want 3", got)
	}
	if !m.marked[0] || !m.marked[1] || !m.marked[2] {
		t.Fatalf("marked = %v, want messages 0, 1 and 2", m.marked)
	}
}

func TestSpaceUnmarksAMessage(t *testing.T) {
	m := selectMode(t, longConversation())
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, ' ')

	if m.markedCount() != 0 {
		t.Fatalf("marked = %v, want nothing after toggling twice", m.marked)
	}
}

func TestMarkingSkipsAroundTheCursor(t *testing.T) {
	m := selectMode(t, longConversation())
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')
	m, _ = pressRune(t, m, 'j')
	m, _ = pressRune(t, m, ' ')

	payload, first, last, count := m.markdownSelection()
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
	if first != 0 || last != 2 {
		t.Fatalf("span = %d..%d, want 0..2", first, last)
	}
	if strings.Count(payload, "###") != 2 {
		t.Fatalf("payload headings = %q", payload)
	}
}

func TestSelectMarksTheMargin(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')

	stripes := map[rune]int{}
	for _, line := range strings.Split(plain(m.View()), "\n") {
		runes := []rune(line)
		if len(runes) == 0 {
			continue
		}
		if runes[0] == '#' || runes[0] == '>' {
			stripes[runes[0]]++
		}
	}
	if stripes['#'] < 1 {
		t.Fatalf("no marked stripe: %v", stripes)
	}
	if stripes['>'] < 1 {
		t.Fatalf("no cursor stripe: %v", stripes)
	}
}

func TestSelectEscLeavesWithoutCopying(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, 'j')
	m, _ = pressRune(t, m, 'k')
	m, _ = press(t, m, tea.KeyEsc)

	if m.mode != modeDetail {
		t.Fatalf("mode = %v, want detail", m.mode)
	}
	if m.selection != nil {
		t.Fatal("the selection survived escaping")
	}
	if m.marked != nil {
		t.Fatal("the marks survived escaping")
	}
	if m.notice != "" {
		t.Fatalf("notice = %q, want none", m.notice)
	}
}

func TestSelectSkipsCollapsedRuns(t *testing.T) {
	conv := atuin.Conversation{Messages: []atuin.Message{
		message("user", atuin.Content{Type: "text", Text: "start"}),
		message("attachment"),
		message("attachment"),
		message("attachment"),
		message("assistant", atuin.Content{Type: "text", Text: "done"}),
	}}
	m := selectMode(t, conv)

	if len(m.selection) != 2 {
		t.Fatalf("selectable = %v, want the two messages with content", m.selection)
	}
	if m.selection[1] != 4 {
		t.Fatalf("last selectable = %d, want message 4", m.selection[1])
	}
}

func TestCopyUsesTheNativeTool(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, ' ')

	var copied string
	m.lookPath = func(string) (string, error) { return "/bin/wl-copy", nil }
	m.copyText = func(tool clipboardTool, text string) error {
		copied = text
		return nil
	}

	m, cmd := pressRune(t, m, 'y')
	m = feedCopy(t, m, cmd)
	if !strings.Contains(copied, "Please bootstrap this repository.") {
		t.Fatalf("payload = %q", copied)
	}
	if m.notice != "copied 1 messages" {
		t.Fatalf("notice = %q", m.notice)
	}
}

func TestCopyFallsBackToTheTerminal(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, ' ')

	var sequence string
	m.lookPath = func(string) (string, error) { return "", errors.New("missing") }
	m.clipboard = func(text string) string {
		sequence = text
		return ""
	}

	m, cmd := pressRune(t, m, 'y')
	m = feedCopy(t, m, cmd)
	if sequence == "" {
		t.Fatal("the terminal clipboard was not called")
	}
	if m.notice != "copied 1 messages (terminal)" {
		t.Fatalf("notice = %q", m.notice)
	}
}

func feedCopy(t *testing.T, m Model, cmd tea.Cmd) Model {
	t.Helper()
	if cmd == nil {
		t.Fatal("no copy command")
	}
	next, _ := m.Update(cmd())
	return next.(Model)
}

func TestCopyDoesNotBlockOnADaemonisingTool(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "att-fake-clip")
	body := "#!/bin/sh\nsleep 5 &\nexit 0\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	done := make(chan error, 1)
	go func() {
		done <- runClipboard(clipboardTool{name: "att-fake-clip"}, "hello")
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the clipboard call waited for a background child holding stderr")
	}
}

func TestCopyTooLargePointsAtTheFile(t *testing.T) {
	conv := atuin.Conversation{Messages: []atuin.Message{
		message("user", atuin.Content{Type: "text", Text: strings.Repeat("x", clipboardLimit+1)}),
	}}
	m := selectMode(t, conv)
	m, _ = pressRune(t, m, ' ')
	m.lookPath = func(string) (string, error) { return "/bin/wl-copy", nil }
	m.copyText = func(clipboardTool, string) error {
		t.Fatal("the clipboard was used for an oversized payload")
		return nil
	}

	m, _ = pressRune(t, m, 'y')
	if !strings.Contains(m.notice, "too large") {
		t.Fatalf("notice = %q", m.notice)
	}
}

func TestWriteSelectionWritesAFile(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, ' ')

	var name, payload string
	m.writeFile = func(dir, file, text string) (string, error) {
		name, payload = file, text
		return file, nil
	}

	m, _ = pressRune(t, m, 'w')
	if name != "att-2cf53ec4-1-1.md" {
		t.Fatalf("name = %q", name)
	}
	if !strings.Contains(payload, "Please bootstrap this repository.") {
		t.Fatalf("payload = %q", payload)
	}
	if m.notice != "wrote "+name {
		t.Fatalf("notice = %q", m.notice)
	}
}

func TestExportThreadCoversEveryVisibleMessage(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)

	var payload string
	m.writeFile = func(dir, file, text string) (string, error) {
		payload = text
		return file, nil
	}

	m, _ = pressRune(t, m, 'e')
	if !strings.Contains(payload, "Please bootstrap this repository.") {
		t.Fatalf("thread export misses the prompt:\n%s", payload)
	}
	if !strings.Contains(payload, "hello there") || !strings.Contains(payload, "$ ls") {
		t.Fatalf("thread export misses the reply:\n%s", payload)
	}
	if m.mode != modeDetail {
		t.Fatalf("mode = %v, want detail", m.mode)
	}
}

func TestSelectStatusAndHints(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')

	status := plain(m.status())
	if !strings.Contains(status, "1 marked of") {
		t.Fatalf("status = %q", status)
	}
	if !strings.Contains(status, "line ") {
		t.Fatalf("status lacks the line position: %q", status)
	}

	hints := plain(m.hints())
	for _, key := range []string{"space mark", "y copy", "w file", "e thread", "esc back"} {
		if !strings.Contains(hints, key) {
			t.Fatalf("hints %q lacks %q", hints, key)
		}
	}
}

func TestDetailHintsMentionSelection(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	hints := plain(m.hints())
	for _, key := range []string{"v select", "e export", "t types"} {
		if !strings.Contains(hints, key) {
			t.Fatalf("hints %q lacks %q", hints, key)
		}
	}
}
