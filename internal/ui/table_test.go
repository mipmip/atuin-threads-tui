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
		layout := newListLayout(width, measureThreads(threads))
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
	// Long content on every axis, so the table needs the full width and the
	// columns have to give way.
	wide := listMetrics{msgs: 969, agent: 6, title: 120, dir: 40}

	cases := []struct {
		width int
		want  []string
		drop  []string
	}{
		{120, []string{"date", "agent", "msgs", "dir"}, nil},
		{100, []string{"date", "agent", "msgs", "dir"}, nil},
		{70, []string{"date", "agent", "msgs", "dir"}, nil},
		{69, []string{"date", "agent", "msgs"}, []string{"dir"}},
		{60, []string{"date", "agent", "msgs"}, []string{"dir"}},
		{59, []string{"agent", "msgs"}, []string{"date", "dir"}},
		{46, []string{"agent", "msgs"}, []string{"date", "dir"}},
		{45, []string{"msgs"}, []string{"date", "agent", "dir"}},
		{38, []string{"msgs"}, []string{"date", "agent", "dir"}},
		{37, nil, []string{"date", "agent", "msgs", "dir"}},
	}

	for _, tc := range cases {
		header := plain(newListLayout(tc.width, wide).header())
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
		if got := lipgloss.Width(header); got > tc.width {
			t.Errorf("width %d: header is %d cells: %q", tc.width, got, header)
		}
	}
}

func TestTableTitleFollowsContent(t *testing.T) {
	threads := sampleThreads()
	layout := newListLayout(120, measureThreads(threads))

	if layout.title != lipgloss.Width("OpenSpec and Beans project bootstrap") {
		t.Errorf("title column = %d cells, want the longest title", layout.title)
	}
	if layout.total() >= 120 {
		t.Errorf("table is %d cells on a 120 cell screen, want room to spare", layout.total())
	}
	if row := plain(layout.row(threads[0], false)); !strings.Contains(row, "atuin-threads-tui") {
		t.Errorf("directory is cut although the title left room: %q", row)
	}
}

func TestTableTitleIsCapped(t *testing.T) {
	long := atuin.Thread{Title: strings.Repeat("x", 200), CWD: "/tmp/some-project-name"}
	layout := newListLayout(200, measureThreads([]atuin.Thread{long}))

	if layout.title != titleMax {
		t.Errorf("title column = %d cells, want the cap of %d", layout.title, titleMax)
	}
	if layout.dir != lipgloss.Width("some-project-name") {
		t.Errorf("dir column = %d cells, want the basename width", layout.dir)
	}
}

func TestTableAgentColumnFollowsCode(t *testing.T) {
	threads := []atuin.Thread{
		{Harness: "claude-code", Title: "a"},
		{Harness: "opencode", Title: "b"},
		{Harness: "opencode-fork", Title: "c"},
	}
	layout := newListLayout(120, measureThreads(threads))

	if layout.agent != agentMaxWidth {
		t.Errorf("agent column = %d cells, want %d", layout.agent, agentMaxWidth)
	}
	if header := plain(layout.header()); !strings.Contains(header, "agent") {
		t.Errorf("header does not grow with the column: %q", header)
	}
	if row := plain(layout.row(threads[0], false)); !strings.Contains(row, "claude") {
		t.Errorf("agent code is not shown in full: %q", row)
	}
	if row := plain(layout.row(threads[2], false)); !strings.Contains(row, "openc") ||
		strings.Contains(row, "opencode-fork") {
		t.Errorf("agent code is not cut to the column: %q", row)
	}
}

func TestTableRowCells(t *testing.T) {
	thread := sampleThreads()[0]
	layout := newListLayout(120, measureThreads(sampleThreads()))
	row := plain(layout.row(thread, true))
	header := plain(layout.header())

	if !strings.HasPrefix(row, "> 2cf53ec4  ") {
		t.Errorf("row prefix = %q", row[:14])
	}
	if got := strings.Index(row, "27"); got != 36 {
		t.Errorf("msgs starts at %d, want 36: %q", got, row)
	}
	if got := strings.Index(row, "OpenSpec and Beans project bootstrap"); got != 40 {
		t.Errorf("title starts at %d, want 40: %q", got, row)
	}
	if got := strings.Index(row, "atuin-threads-tui"); got != 78 {
		t.Errorf("dir starts at %d, want 78: %q", got, row)
	}
	if got := strings.Index(header, "msgs"); got != 34 {
		t.Errorf("msgs header starts at %d, want 34: %q", got, header)
	}
	if got := strings.Index(header, "title"); got != 40 {
		t.Errorf("title header starts at %d, want 40: %q", got, header)
	}
	if got := strings.Index(header, "dir"); got != 78 {
		t.Errorf("dir header starts at %d, want 78: %q", got, header)
	}
}

func TestTableWideTitleTruncates(t *testing.T) {
	thread := sampleThreads()[0]
	thread.Title = strings.Repeat("wide title ", 20)

	for _, width := range []int{120, 80} {
		layout := newListLayout(width, measureThreads([]atuin.Thread{thread}))
		row := plain(layout.row(thread, false))
		if got := lipgloss.Width(row); got != layout.total() {
			t.Errorf("width %d: row is %d cells, table is %d", width, got, layout.total())
		}
		if lipgloss.Width(row) > width {
			t.Errorf("width %d: row is %d cells", width, lipgloss.Width(row))
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
		{"claude-code", "claude"},
		{"opencode", "ocode"},
		{"OpenCode", "ocode"},
		{"codex", "codex"},
		{"pi", "pi"},
		{"opencode-fork", "openco"},
		{"claude-code-enterprise", "claude"},
		{"cursor-agent-with-a-long-name", "cursor"},
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
	if !strings.Contains(lines[3], "OpenSpec and Beans") {
		t.Errorf("first thread row missing:\n%s", plain(m.View()))
	}
}

func TestListEmptyStateHasNoTable(t *testing.T) {
	m := loaded(t, &fakeSource{})
	if view := plain(m.View()); strings.Contains(view, "title") {
		t.Errorf("empty state shows a table header:\n%s", view)
	}
}
