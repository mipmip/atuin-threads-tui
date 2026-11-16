package ui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func TestTableColumnsLineUp(t *testing.T) {
	threads := sampleThreads()

	for _, width := range []int{120, 100, 99, 80, 56, 55, 42, 41, 38, 37, 24} {
		layout := newListLayout(width, widestMessages(threads))
		header := plain(layout.header())
		rule := plain(layout.rule())

		if got := lipgloss.Width(header); got != lipgloss.Width(rule) {
			t.Errorf("width %d: header %d cells, rule %d cells", width, got, lipgloss.Width(rule))
		}
		if got := lipgloss.Width(header); got > width {
			t.Errorf("width %d: header is %d cells", width, got)
		}

		for _, thread := range threads {
			row := plain(layout.row(thread, false))
			if got := lipgloss.Width(row); got != lipgloss.Width(header) {
				t.Errorf("width %d: row %d cells, header %d cells: %q", width, got, lipgloss.Width(header), row)
			}
		}
	}
}

func TestTableDropsColumns(t *testing.T) {
	cases := []struct {
		width int
		want  []string
		drop  []string
	}{
		{120, []string{"date", "ag", "msgs", "dir"}, nil},
		{100, []string{"date", "ag", "msgs", "dir"}, nil},
		{99, []string{"date", "ag", "msgs"}, []string{"dir"}},
		{80, []string{"date", "ag", "msgs"}, []string{"dir"}},
		{56, []string{"date", "ag", "msgs"}, []string{"dir"}},
		{55, []string{"ag", "msgs"}, []string{"date"}},
		{42, []string{"ag", "msgs"}, []string{"date"}},
		{41, []string{"msgs"}, []string{"date", "ag"}},
		{38, []string{"msgs"}, []string{"date", "ag"}},
		{37, nil, []string{"date", "ag", "msgs"}},
	}

	for _, tc := range cases {
		header := plain(newListLayout(tc.width, 969).header())
		for _, want := range tc.want {
			if !strings.Contains(header, want) {
				t.Errorf("width %d: header %q misses %q", tc.width, header, want)
			}
		}
		for _, drop := range tc.drop {
			if strings.Contains(header, drop) {
				t.Errorf("width %d: header %q still shows %q", tc.width, header, drop)
			}
		}
	}
}

func TestTableRowCells(t *testing.T) {
	thread := sampleThreads()[0]
	layout := newListLayout(120, 969)
	row := plain(layout.row(thread, true))
	header := plain(layout.header())

	if !strings.HasPrefix(row, "> 2cf53ec4  ") {
		t.Errorf("row prefix = %q", row[:14])
	}
	if got := strings.Index(row, "OpenSpec and Beans project bootstrap"); got != 30 {
		t.Errorf("title starts at %d, want 30: %q", got, row)
	}
	if got := strings.Index(row, "27"); got != 98 {
		t.Errorf("msgs starts at %d, want 98: %q", got, row)
	}
	if got := strings.Index(row, "atuin-threads-tui"); got != 102 {
		t.Errorf("dir starts at %d, want 102: %q", got, row)
	}
	if got := strings.Index(header, "title"); got != 30 {
		t.Errorf("title header starts at %d, want 30: %q", got, header)
	}
	if got := strings.Index(header, "msgs"); got != 96 {
		t.Errorf("msgs header starts at %d, want 96: %q", got, header)
	}
	if got := strings.Index(header, "dir"); got != 102 {
		t.Errorf("dir header starts at %d, want 102: %q", got, header)
	}
}

func TestTableWideTitleTruncates(t *testing.T) {
	thread := sampleThreads()[0]
	thread.Title = strings.Repeat("wide title ", 20)

	for _, width := range []int{120, 80} {
		layout := newListLayout(width, 27)
		row := plain(layout.row(thread, false))
		if lipgloss.Width(row) != width {
			t.Errorf("width %d: row is %d cells: %q", width, lipgloss.Width(row), row)
		}
		if !strings.Contains(row, "…") {
			t.Errorf("width %d: title was not truncated: %q", width, row)
		}
	}
}

func TestAgentTag(t *testing.T) {
	cases := []struct {
		harness string
		want    string
	}{
		{"opencode", "oc"},
		{"claude-code", "cc"},
		{"OpenCode", "oc"},
		{"codex", "co"},
		{"", "--"},
		{"  ", "--"},
		{"x", "x"},
	}
	for _, tc := range cases {
		if got := agentTag(tc.harness); got != tc.want {
			t.Errorf("agentTag(%q) = %q, want %q", tc.harness, got, tc.want)
		}
	}
}

func TestDateLabel(t *testing.T) {
	updated := time.Date(2026, 11, 16, 11, 6, 0, 0, time.UTC)
	thread := atuin.Thread{UpdatedAt: updated}
	if got, want := dateLabel(thread), updated.Local().Format("02 Jan 15:04"); got != want {
		t.Errorf("dateLabel = %q, want %q", got, want)
	}
	if got := dateLabel(atuin.Thread{}); got != "-" {
		t.Errorf("dateLabel(zero) = %q, want %q", got, "-")
	}
	if got := lipgloss.Width(dateLabel(thread)); got != dateWidth {
		t.Errorf("date width = %d, want %d", got, dateWidth)
	}
}

func TestDirLabel(t *testing.T) {
	cases := []struct {
		cwd  string
		want string
	}{
		{"/home/pim/gh.mipmip/atuin-threads-tui", "atuin-threads-tui"},
		{"/home/pim/gh.mipmip/atuin-threads-tui/", "atuin-threads-tui"},
		{"relative/path", "path"},
		{"/", "/"},
		{"", "-"},
	}
	for _, tc := range cases {
		if got := dirLabel(tc.cwd); got != tc.want {
			t.Errorf("dirLabel(%q) = %q, want %q", tc.cwd, got, tc.want)
		}
	}
}

func TestTailTruncate(t *testing.T) {
	cases := []struct {
		s     string
		width int
		want  string
	}{
		{"atuin-threads-tui", 18, "atuin-threads-tui"},
		{"a-very-long-project-name", 10, "…ject-name"},
		{"日本語のプロジェクト", 7, "…ェクト"},
		{"abc", 1, "…"},
		{"abc", 0, ""},
	}
	for _, tc := range cases {
		got := tailTruncate(tc.s, tc.width)
		if got != tc.want {
			t.Errorf("tailTruncate(%q, %d) = %q, want %q", tc.s, tc.width, got, tc.want)
		}
		if tc.width > 0 && lipgloss.Width(got) > tc.width {
			t.Errorf("tailTruncate(%q, %d) = %q is %d cells", tc.s, tc.width, got, lipgloss.Width(got))
		}
	}
}

func TestTableScrollingLeavesRoomForHeader(t *testing.T) {
	threads := sampleThreads()
	for i := 0; i < 30; i++ {
		threads = append(threads, atuin.Thread{
			AtuinID:      "01a1d43c1f237541a0e1c6c42cf9999",
			Title:        "extra thread",
			MessageCount: i,
		})
	}

	m := loaded(t, &fakeSource{threads: threads})
	m, _ = feed(t, m, tea.WindowSizeMsg{Width: 120, Height: 24})

	if m.bodyHeight() != 21 || m.listHeight() != 19 {
		t.Fatalf("body height = %d, list height = %d, want 21 and 19", m.bodyHeight(), m.listHeight())
	}
	if lines := strings.Split(plain(m.View()), "\n"); len(lines) != 24 {
		t.Errorf("view has %d lines, want 24", len(lines))
	}

	for i := 0; i < m.listHeight(); i++ {
		m, _ = pressRune(t, m, 'j')
	}
	if m.offset == 0 {
		t.Errorf("list did not scroll: cursor=%d offset=%d", m.cursor, m.offset)
	}
	if m.cursor-m.offset >= m.listHeight() {
		t.Errorf("selection is off screen: cursor=%d offset=%d listHeight=%d", m.cursor, m.offset, m.listHeight())
	}

	body := strings.Split(plain(m.View()), "\n")
	if !strings.Contains(body[1], "id") || !strings.Contains(body[2], "─") {
		t.Errorf("header is not at the top of the body:\n%s", strings.Join(body[:4], "\n"))
	}
}

func TestListViewShowsTableHeader(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})
	lines := strings.Split(plain(m.View()), "\n")

	if !strings.Contains(lines[1], "title") || !strings.Contains(lines[1], "msgs") {
		t.Errorf("header row missing:\n%s", plain(m.View()))
	}
	if !strings.Contains(lines[2], "────") {
		t.Errorf("rule row missing:\n%s", plain(m.View()))
	}
	if !strings.Contains(lines[3], "OpenSpec and Beans project") {
		t.Errorf("first thread row missing:\n%s", plain(m.View()))
	}
}

func TestListEmptyStateHasNoTable(t *testing.T) {
	m := loaded(t, &fakeSource{})
	if view := plain(m.View()); strings.Contains(view, "title") {
		t.Errorf("empty state shows a table header:\n%s", view)
	}
}
