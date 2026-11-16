package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func TestScrollbarHiddenWhenThreadFits(t *testing.T) {
	if bar := newScrollbar(0, 20, 20); bar.shown {
		t.Errorf("bar shown for a thread that exactly fits: %+v", bar)
	}
	if bar := newScrollbar(0, 5, 20); bar.shown {
		t.Errorf("bar shown for a short thread: %+v", bar)
	}
	if bar := newScrollbar(0, 20, 0); bar.shown {
		t.Errorf("bar shown without a track: %+v", bar)
	}
}

func TestScrollbarThumbIsClamped(t *testing.T) {
	bar := newScrollbar(0, 10000, 20)
	if bar.thumb != 1 {
		t.Errorf("thumb = %d for a very long thread, want 1", bar.thumb)
	}
	bar = newScrollbar(0, 21, 20)
	if bar.thumb < 1 || bar.thumb > 20 {
		t.Errorf("thumb = %d, want between 1 and the track", bar.thumb)
	}
}

func TestScrollbarThumbFollowsTheHeightRatio(t *testing.T) {
	bar := newScrollbar(0, 200, 20)
	if bar.thumb != 2 {
		t.Errorf("thumb = %d for a thread ten times the track, want 2", bar.thumb)
	}
}

func TestScrollbarStartsAtTheTop(t *testing.T) {
	bar := newScrollbar(0, 200, 20)
	if bar.top != 0 {
		t.Errorf("top = %d at scroll 0, want 0", bar.top)
	}
}

func TestScrollbarEndsAtTheBottom(t *testing.T) {
	bar := newScrollbar(180, 200, 20)
	if got := bar.top + bar.thumb; got != 20 {
		t.Errorf("thumb ends at row %d, want the last row 20", got)
	}
}

func TestScrollbarMovesWithTheScroll(t *testing.T) {
	bar := newScrollbar(90, 200, 20)
	if bar.top == 0 {
		t.Errorf("top = %d halfway down a long thread, want it moved", bar.top)
	}
}

func TestScrollbarStaysInsideTheTrack(t *testing.T) {
	for scroll := -5; scroll <= 300; scroll += 5 {
		bar := newScrollbar(scroll, 200, 20)
		if bar.top < 0 || bar.top+bar.thumb > 20 {
			t.Fatalf("scroll %d puts the thumb at %d..%d", scroll, bar.top, bar.top+bar.thumb)
		}
	}
}

func lastCell(line string) string {
	runes := []rune(plain(line))
	if len(runes) == 0 {
		return ""
	}
	return string(runes[len(runes)-1])
}

func TestDetailShowsTheBarForALongThread(t *testing.T) {
	m, _ := openDetail(t, longConversation(), 120, 24)
	view := m.View()
	rows := strings.Split(view, "\n")
	if len(rows) != 24 {
		t.Fatalf("view has %d rows, want 24", len(rows))
	}

	body := rows[1 : 1+m.bodyHeight()]
	if got := lastCell(body[0]); got != scrollThumbCell {
		t.Errorf("first body row ends in %q, want the thumb %q", got, scrollThumbCell)
	}
	if got := lastCell(body[len(body)-1]); got != scrollTrackCell {
		t.Errorf("last body row ends in %q, want the track %q", got, scrollTrackCell)
	}
}

func TestDetailBarReachesTheBottomAtTheEnd(t *testing.T) {
	m, _ := openDetail(t, longConversation(), 120, 24)
	m, _ = pressRune(t, m, 'G')

	body := strings.Split(m.View(), "\n")[1 : 1+m.bodyHeight()]
	if got := lastCell(body[len(body)-1]); got != scrollThumbCell {
		t.Errorf("last body row ends in %q after G, want the thumb", got)
	}
	if got := lastCell(body[0]); got != scrollTrackCell {
		t.Errorf("first body row ends in %q after G, want the track", got)
	}
}

func TestShortThreadKeepsTheColumnReserved(t *testing.T) {
	long, _ := openDetail(t, longConversation(), 120, 24)
	short, _ := openDetail(t, sampleConversation(), 120, 24)
	if long.detail.built == false || short.detail.built == false {
		t.Fatal("a conversation was not rendered")
	}

	longBody := strings.Split(long.View(), "\n")[1]
	shortBody := strings.Split(short.View(), "\n")[1]
	if lastCell(shortBody) != "" && lastCell(shortBody) != " " {
		t.Errorf("short thread shows %q in the margin, want it empty", lastCell(shortBody))
	}
	if lipgloss.Width(longBody) != lipgloss.Width(shortBody) {
		t.Errorf("body width differs: long %d, short %d", lipgloss.Width(longBody), lipgloss.Width(shortBody))
	}
}

func TestWideCharacterCannotPushTheBarOffTheEdge(t *testing.T) {
	conv := sampleConversation()
	conv.Messages = nil
	for i := 0; i < 40; i++ {
		conv.Messages = append(conv.Messages, message("assistant",
			atuin.Content{Type: "text", Text: strings.Repeat("日本語のプロジェクト", 3)}))
	}

	m, _ := openDetail(t, conv, 80, 24)
	for _, line := range strings.Split(m.View(), "\n") {
		if lipgloss.Width(line) > 80 {
			t.Fatalf("line is %d cells wide, want 80: %q", lipgloss.Width(line), plain(line))
		}
	}
	if got := lastCell(strings.Split(m.View(), "\n")[1]); got != scrollThumbCell {
		t.Errorf("first body row ends in %q, want the thumb", got)
	}
}

func TestBarIsAbsentInTheListAndWhileRendering(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads()})
	for _, line := range strings.Split(m.View(), "\n") {
		if strings.HasSuffix(plain(line), scrollTrackCell) {
			t.Fatalf("list view shows a track: %q", plain(line))
		}
	}

	m, _ = openDetail(t, bulkConversation(600), 120, 24)
	if !m.rendering {
		t.Fatal("the large thread did not render in the background")
	}
	for _, line := range strings.Split(m.View(), "\n") {
		if strings.HasSuffix(plain(line), scrollTrackCell) || strings.HasSuffix(plain(line), scrollThumbCell) {
			t.Fatalf("bar drawn while rendering: %q", plain(line))
		}
	}
}
