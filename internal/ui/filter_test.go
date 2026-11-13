package ui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func typeText(t *testing.T, m Model, text string) Model {
	t.Helper()
	for _, r := range text {
		m, _ = feed(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	return m
}

func TestFilterNarrowsList(t *testing.T) {
	source := &fakeSource{threads: sampleThreads(), conversation: sampleConversation()}
	m := loaded(t, source)

	m, _ = pressRune(t, m, '/')
	if m.mode != modeFilter {
		t.Fatalf("mode = %v, want filter", m.mode)
	}

	m = typeText(t, m, "second")
	if len(m.visible()) != 1 {
		t.Fatalf("visible = %d, want 1", len(m.visible()))
	}

	view := plain(m.View())
	for _, want := range []string{"1 of 3 threads", "Second thread", "filter: second"} {
		if !strings.Contains(view, want) {
			t.Errorf("view missing %q:\n%s", want, view)
		}
	}
	if strings.Contains(view, "OpenSpec and Beans project bootstrap") {
		t.Errorf("filtered view still shows another thread:\n%s", view)
	}

	m, _ = press(t, m, tea.KeyEnter)
	if m.mode != modeList {
		t.Fatalf("mode = %v, want list after enter", m.mode)
	}
	if m.filter != "second" {
		t.Errorf("filter = %q, want it applied", m.filter)
	}

	m, cmd := press(t, m, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("enter on a filtered list returned no command")
	}
	m, _ = feed(t, m, cmd())
	if source.showID != source.threads[1].AtuinID {
		t.Errorf("opened %q, want the filtered thread %q", source.showID, source.threads[1].AtuinID)
	}
}

func TestFilterMatchesOtherFields(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})
	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "GH.MIPMIP")

	if len(m.visible()) != 1 || m.visible()[0].Title != "OpenSpec and Beans project bootstrap" {
		t.Errorf("cwd match failed: %d visible", len(m.visible()))
	}
}

func TestFilterNoMatch(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})
	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "zzzz")

	view := plain(m.View())
	if !strings.Contains(view, "no thread matches the filter") || !strings.Contains(view, "0 of 3 threads") {
		t.Errorf("no-match view unexpected:\n%s", view)
	}
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
}

func TestFilterClearAndBackspace(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})
	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "abc")

	m, _ = press(t, m, tea.KeyBackspace)
	if m.filterInput != "ab" {
		t.Errorf("filter input = %q, want ab", m.filterInput)
	}

	m, _ = press(t, m, tea.KeyCtrlU)
	if m.filterInput != "" || m.filter != "" {
		t.Errorf("ctrl+u left %q / %q", m.filterInput, m.filter)
	}

	m = typeText(t, m, "zzz")
	m, _ = press(t, m, tea.KeyEsc)
	if m.mode != modeList {
		t.Errorf("mode = %v, want list", m.mode)
	}
	if m.filter != "" || m.filterInput != "" {
		t.Errorf("esc left filter %q", m.filter)
	}
	if !strings.Contains(plain(m.View()), "3 threads") {
		t.Errorf("esc did not restore the list:\n%s", plain(m.View()))
	}
}

func TestFilterNavigationBounds(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})
	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "thread")
	m, _ = press(t, m, tea.KeyEnter)

	rows := len(m.visible())
	if rows != 2 {
		t.Fatalf("visible = %d, want 2", rows)
	}

	m, _ = pressRune(t, m, 'G')
	if m.cursor != rows-1 {
		t.Errorf("cursor = %d, want %d", m.cursor, rows-1)
	}
	m, _ = pressRune(t, m, 'j')
	if m.cursor != rows-1 {
		t.Errorf("cursor = %d, want it clamped", m.cursor)
	}
}

func TestFullTextSearch(t *testing.T) {
	result := atuin.Thread{Title: "Found by full text", AtuinID: "ffffffff"}
	source := &fakeSource{
		threads:       sampleThreads(),
		searchResults: []atuin.Thread{result},
	}
	m := loaded(t, source)

	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "bootstrap")
	m, _ = press(t, m, tea.KeyEnter)

	m, cmd := press(t, m, tea.KeyCtrlS)
	if cmd == nil {
		t.Fatal("ctrl+s returned no command")
	}
	if !m.searching {
		t.Error("model is not in a searching state")
	}

	m, _ = feed(t, m, cmd())
	if source.searchQuery != "bootstrap" || source.searchLimit != 100 {
		t.Errorf("search used %q / %d, want bootstrap / 100", source.searchQuery, source.searchLimit)
	}
	if len(m.threads) != 1 || m.threads[0].Title != "Found by full text" {
		t.Fatalf("threads = %+v, want the search result", m.threads)
	}
	if m.searchQuery != "bootstrap" || m.searching {
		t.Errorf("search state: query=%q searching=%v", m.searchQuery, m.searching)
	}
	if view := plain(m.View()); !strings.Contains(view, "1 results for \"bootstrap\"") {
		t.Errorf("view missing result count:\n%s", view)
	}
}

func TestFullTextSearchNeedsQuery(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})

	m, cmd := press(t, m, tea.KeyCtrlS)
	if cmd != nil {
		t.Error("ctrl+s with an empty filter returned a command")
	}
	if m.searching {
		t.Error("model is searching without a query")
	}
}

func TestFullTextSearchFailure(t *testing.T) {
	source := &fakeSource{
		threads:   sampleThreads(),
		searchErr: errors.New("search unavailable"),
	}
	m := loaded(t, source)

	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "x")
	m, _ = press(t, m, tea.KeyEnter)
	m, cmd := press(t, m, tea.KeyCtrlS)
	m, _ = feed(t, m, cmd())

	if m.err == nil || !strings.Contains(plain(m.View()), "search unavailable") {
		t.Errorf("error not shown: %v\n%s", m.err, plain(m.View()))
	}
	if len(m.threads) != 3 {
		t.Errorf("threads replaced after a failed search: %d", len(m.threads))
	}
}

func TestClearSearchReloadsList(t *testing.T) {
	source := &fakeSource{
		threads:       sampleThreads(),
		searchResults: []atuin.Thread{{Title: "hit"}},
	}
	m := loaded(t, source)

	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "hit")
	m, _ = press(t, m, tea.KeyEnter)
	m, cmd := press(t, m, tea.KeyCtrlS)
	m, _ = feed(t, m, cmd())

	m, cmd = press(t, m, tea.KeyEsc)
	if cmd == nil {
		t.Fatal("esc did not issue a reload command")
	}
	if m.searchQuery != "" {
		t.Errorf("search query = %q, want empty", m.searchQuery)
	}
	if !m.loading {
		t.Error("model is not loading after clearing the search")
	}

	m, _ = feed(t, m, cmd())
	if len(m.threads) != 3 || strings.Contains(plain(m.View()), "results for") {
		t.Errorf("list not restored: %d threads\n%s", len(m.threads), plain(m.View()))
	}
}

func TestMatches(t *testing.T) {
	thread := sampleThreads()[0]

	for _, query := range []string{"openspec", "opencode", "bootstrap", "2cf53ec4", "/home/pim"} {
		if !matches(thread, strings.ToLower(query)) {
			t.Errorf("matches(%q) = false", query)
		}
	}
	if matches(thread, "no-such-text") {
		t.Error("matches reported a false positive")
	}
}
