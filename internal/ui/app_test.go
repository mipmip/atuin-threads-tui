package ui

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

var sgr = regexp.MustCompile("\x1b\\[[0-9;]*m")

func plain(s string) string { return sgr.ReplaceAllString(s, "") }

type fakeSource struct {
	threads       []atuin.Thread
	conversation  atuin.Conversation
	listErr       error
	showErr       error
	showID        string
	searchResults []atuin.Thread
	searchErr     error
	searchQuery   string
	searchLimit   int
}

func (f *fakeSource) List(context.Context) ([]atuin.Thread, error) {
	return f.threads, f.listErr
}

func (f *fakeSource) Show(_ context.Context, id string) (atuin.Conversation, error) {
	f.showID = id
	return f.conversation, f.showErr
}

func (f *fakeSource) Search(_ context.Context, query string, limit int) ([]atuin.Thread, error) {
	f.searchQuery = query
	f.searchLimit = limit
	return f.searchResults, f.searchErr
}

func sampleThreads() []atuin.Thread {
	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	return []atuin.Thread{
		{
			Harness:      "opencode",
			AtuinID:      "01a1d43c1f237541a0e1c6c42cf53ec4",
			Title:        "OpenSpec and Beans project bootstrap",
			CWD:          "/home/pim/gh.mipmip/atuin-threads-tui",
			MessageCount: 27,
			UpdatedAt:    now.Add(-13 * time.Second),
			Preview:      "Bootstrap the project",
		},
		{
			Harness:      "claude-code",
			AtuinID:      "01a1acb62fa877ab83c0c981da8da47e",
			Title:        "Second thread",
			MessageCount: 604,
			UpdatedAt:    now.Add(-2 * time.Hour),
		},
		{
			Harness:      "codex",
			AtuinID:      "01a1d43c1f237541a0e1c6c42cf1111",
			Title:        "",
			MessageCount: 0,
			UpdatedAt:    now.Add(-72 * time.Hour),
		},
	}
}

func sampleConversation() atuin.Conversation {
	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	return atuin.Conversation{
		Session: atuin.Thread{
			Harness:      "opencode",
			AtuinID:      "01a1d43c1f237541a0e1c6c42cf53ec4",
			Title:        "OpenSpec and Beans project bootstrap",
			Model:        "deepseek-v4.1-flash",
			CWD:          "/home/pim/gh.mipmip/atuin-threads-tui",
			MessageCount: 2,
			UpdatedAt:    now,
		},
		Messages: []atuin.Message{
			{
				Role:      "user",
				Timestamp: now,
				Content: []atuin.Content{
					{Type: "text", Text: "Please bootstrap this repository."},
				},
			},
			{
				Role:      "assistant",
				Timestamp: now,
				Content: []atuin.Content{
					{Type: "thinking", Text: "Considering the layout."},
					{Type: "text", Text: "hello there"},
					{Type: "tool_call", Name: "shell", Input: `{"command":"ls"}`},
					{Type: "tool_result", Content: "README.md", IsError: false},
				},
			},
		},
	}
}

func feed(t *testing.T, m Model, msg tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	updated, cmd := m.Update(msg)
	model, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", updated)
	}
	return model, cmd
}

func press(t *testing.T, m Model, key tea.KeyType) (Model, tea.Cmd) {
	t.Helper()
	return feed(t, m, tea.KeyMsg{Type: key})
}

func pressRune(t *testing.T, m Model, r rune) (Model, tea.Cmd) {
	t.Helper()
	return feed(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
}

func loaded(t *testing.T, source *fakeSource) Model {
	t.Helper()
	m := New(source)
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init returned no command")
	}
	m, _ = feed(t, m, cmd())
	return m
}

func TestLoadAndNavigate(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})

	view := plain(m.View())
	for _, want := range []string{"3 threads", "OpenSpec and Beans", "(untitled)", "> "} {
		if !strings.Contains(view, want) {
			t.Errorf("view missing %q:\n%s", want, view)
		}
	}

	m, _ = pressRune(t, m, 'j')
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}
	m, _ = pressRune(t, m, 'k')
	m, _ = pressRune(t, m, 'k')
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after moving past the top", m.cursor)
	}
	m, _ = pressRune(t, m, 'G')
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want last thread", m.cursor)
	}
	m, _ = pressRune(t, m, 'j')
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want it to stay on the last thread", m.cursor)
	}
	m, _ = press(t, m, tea.KeyHome)
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 after home", m.cursor)
	}
}

func TestEmptyAndErrorStates(t *testing.T) {
	m := loaded(t, &fakeSource{})
	view := plain(m.View())
	if !strings.Contains(view, "no threads captured") || !strings.Contains(view, "0 threads") {
		t.Errorf("empty view unexpected:\n%s", view)
	}

	m = New(&fakeSource{listErr: errors.New("atuin: no such command")})
	m, _ = feed(t, m, m.Init()())
	view = plain(m.View())
	if !strings.Contains(view, "error: atuin: no such command") || !strings.Contains(view, "ctrl+r to retry") {
		t.Errorf("error view unexpected:\n%s", view)
	}

	m, cmd := feed(t, m, tea.KeyMsg{Type: tea.KeyCtrlR})
	if cmd == nil {
		t.Error("ctrl+r returned no reload command")
	}
	if m.err != nil || !m.loading {
		t.Errorf("reload did not reset state: err=%v loading=%v", m.err, m.loading)
	}
}

func TestOpenAndCloseThread(t *testing.T) {
	source := &fakeSource{threads: sampleThreads(), conversation: sampleConversation()}
	m := loaded(t, source)

	m, cmd := press(t, m, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("enter returned no command")
	}

	m, _ = feed(t, m, cmd())
	if source.showID != source.threads[0].AtuinID {
		t.Errorf("Show id = %q, want %q", source.showID, source.threads[0].AtuinID)
	}
	if m.mode != modeDetail {
		t.Fatalf("mode = %v, want detail", m.mode)
	}

	view := plain(m.View())
	for _, want := range []string{
		"OpenSpec and Beans project bootstrap",
		"- user",
		"- asst",
		"hello there",
		"shell",
		"$ ls",
		"README.md",
		"/home/pim/gh.mipmip/atuin-threads-tui",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("detail view missing %q:\n%s", want, view)
		}
	}

	m, _ = press(t, m, tea.KeyEsc)
	if m.mode != modeList {
		t.Errorf("mode = %v, want list after esc", m.mode)
	}
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want selection preserved", m.cursor)
	}
}

func TestNoThreadSelectedOpensNothing(t *testing.T) {
	m := loaded(t, &fakeSource{})
	if _, cmd := press(t, m, tea.KeyEnter); cmd != nil {
		t.Error("enter on an empty list returned a command")
	}
}

func TestOpenFailureStaysInList(t *testing.T) {
	source := &fakeSource{threads: sampleThreads(), showErr: errors.New("thread gone")}
	m := loaded(t, source)

	m, cmd := press(t, m, tea.KeyEnter)
	m, _ = feed(t, m, cmd())

	if m.mode != modeList {
		t.Errorf("mode = %v, want list", m.mode)
	}
	if m.err == nil || !strings.Contains(plain(m.View()), "thread gone") {
		t.Errorf("error not shown: %v\n%s", m.err, plain(m.View()))
	}
}

func longConversation() atuin.Conversation {
	conversation := sampleConversation()
	for i := 0; i < 40; i++ {
		conversation.Messages = append(conversation.Messages, atuin.Message{
			Role:    "assistant",
			Content: []atuin.Content{{Type: "text", Text: fmt.Sprintf("message %d with enough words to wrap onto another line", i)}},
		})
	}
	return conversation
}

func TestDetailScrolling(t *testing.T) {
	source := &fakeSource{threads: sampleThreads(), conversation: longConversation()}
	m := loaded(t, source)
	m, cmd := press(t, m, tea.KeyEnter)
	m, _ = feed(t, m, cmd())

	maxScroll := len(m.detailLines()) - m.bodyHeight()
	if maxScroll <= 0 {
		t.Fatalf("test conversation is not long enough: maxScroll=%d", maxScroll)
	}

	m, _ = pressRune(t, m, 'j')
	if m.detailScroll != 1 {
		t.Errorf("scroll = %d, want 1", m.detailScroll)
	}
	m, _ = press(t, m, tea.KeyPgDown)
	if m.detailScroll != 1+m.bodyHeight() {
		t.Errorf("scroll = %d, want one page", m.detailScroll)
	}
	m, _ = pressRune(t, m, 'G')
	if m.detailScroll != maxScroll {
		t.Errorf("scroll = %d, want max %d", m.detailScroll, maxScroll)
	}
	m, _ = pressRune(t, m, 'j')
	if m.detailScroll != maxScroll {
		t.Errorf("scroll = %d, want it clamped at %d", m.detailScroll, maxScroll)
	}
	m, _ = pressRune(t, m, 'g')
	if m.detailScroll != 0 {
		t.Errorf("scroll = %d, want 0", m.detailScroll)
	}
	m, _ = pressRune(t, m, 'k')
	if m.detailScroll != 0 {
		t.Errorf("scroll = %d, want it clamped at 0", m.detailScroll)
	}
}

func TestResizeKeepsChrome(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})

	m, _ = feed(t, m, tea.WindowSizeMsg{Width: 120, Height: 12})
	if m.width != 120 || m.height != 12 {
		t.Errorf("size = %dx%d", m.width, m.height)
	}
	if lines := strings.Split(m.View(), "\n"); len(lines) != 12 {
		t.Errorf("view has %d lines, want 12:\n%s", len(lines), plain(m.View()))
	}

	m, _ = feed(t, m, tea.WindowSizeMsg{Width: 8, Height: 2})
	if m.bodyHeight() != 0 {
		t.Errorf("body height = %d, want 0", m.bodyHeight())
	}
	_ = m.View()
}

func TestQuitKeys(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})

	if _, cmd := pressRune(t, m, 'q'); cmd == nil {
		t.Error("q in the list returned no quit command")
	}

	m, cmd := press(t, m, tea.KeyEnter)
	m, _ = feed(t, m, cmd())
	if _, cmd := press(t, m, tea.KeyCtrlC); cmd == nil {
		t.Error("ctrl+c in detail returned no quit command")
	}
}

func TestThreadRowFitsWidth(t *testing.T) {
	thread := sampleThreads()[0]

	for _, width := range []int{10, 40, 120} {
		row := plain(newListLayout(width, measureThreads(sampleThreads())).row(thread, true))
		if got := lipgloss.Width(row); got > width {
			t.Errorf("row width = %d, want <= %d: %q", got, width, row)
		}
	}
	if !strings.HasPrefix(plain(newListLayout(80, measureThreads(sampleThreads())).row(thread, true)), "> ") {
		t.Error("selected row has no marker")
	}
}

func TestFormattingHelpers(t *testing.T) {
	if got := shortID("01a1d43c1f237541a0e1c6c42cf53ec4"); got != "2cf53ec4" {
		t.Errorf("shortID = %q", got)
	}
	if got := shortID("abc"); got != "abc" {
		t.Errorf("shortID = %q", got)
	}

	now := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		then time.Time
		want string
	}{
		{now.Add(-5 * time.Second), "5s"},
		{now.Add(-5 * time.Minute), "5m"},
		{now.Add(-5 * time.Hour), "5h"},
		{now.Add(-5 * 24 * time.Hour), "5d"},
		{now.Add(time.Minute), "now"},
		{time.Time{}, "?"},
	}
	for _, tc := range cases {
		if got := age(now, tc.then); got != tc.want {
			t.Errorf("age(%v) = %q, want %q", tc.then, got, tc.want)
		}
	}

	if got := truncate("hello", 10); got != "hello" {
		t.Errorf("truncate = %q", got)
	}
	if got := truncate("hello world", 6); got != "hello…" {
		t.Errorf("truncate = %q", got)
	}
	if got := truncate("hello", 0); got != "" {
		t.Errorf("truncate = %q", got)
	}
	if got := truncate("hello", 1); got != "…" {
		t.Errorf("truncate = %q", got)
	}
}

func TestWrap(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		width int
		want  []string
	}{
		{"words", "one two three", 7, []string{"one two", "three"}},
		{"paragraphs", "one\n\ntwo", 10, []string{"one", "", "two"}},
		{"hard split", "abcdefghij", 4, []string{"abcd", "efgh", "ij"}},
		{"zero width", "a b", 0, []string{"a", "b"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := wrap(tc.text, tc.width)
			if strings.Join(got, "|") != strings.Join(tc.want, "|") {
				t.Errorf("wrap = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBlockLines(t *testing.T) {
	if got := blockLines(atuin.Content{Type: "text", Text: "hello"}, 40, 40); len(got) != 1 || got[0] != "hello" {
		t.Errorf("text block = %q", got)
	}

	thinking := blockLines(atuin.Content{Type: "thinking", Text: "hmm"}, 40, 40)
	if len(thinking) != 1 || !strings.Contains(plain(thinking[0]), "hmm") {
		t.Errorf("thinking block = %q", thinking)
	}

	result := blockLines(atuin.Content{Type: "tool_result", Content: "boom", IsError: true}, 40, 40)
	if !strings.Contains(plain(strings.Join(result, " ")), "! boom") {
		t.Errorf("error result = %q", result)
	}

	var input strings.Builder
	for i := 0; i < 20; i++ {
		fmt.Fprintf(&input, "line%d\n", i)
	}
	cappedLines := blockLines(atuin.Content{Type: "tool_call", Name: "shell", Input: input.String()}, 40, 40)
	if len(cappedLines) != maxBlockLines+1 {
		t.Fatalf("capped block has %d lines, want %d", len(cappedLines), maxBlockLines+1)
	}
	if !strings.Contains(plain(cappedLines[maxBlockLines]), "more lines") {
		t.Errorf("missing truncation note: %q", plain(cappedLines[maxBlockLines]))
	}

	if got := blockLines(atuin.Content{Type: "unknown"}, 40, 40); len(got) != 1 || !strings.Contains(plain(got[0]), "unknown") {
		t.Errorf("unknown block = %q, want a label", got)
	}
}
