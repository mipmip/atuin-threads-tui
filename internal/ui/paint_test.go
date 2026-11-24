package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestPaintBackgroundPadsToTheWidth(t *testing.T) {
	painted := paintBackground(" a", backgroundEven, 5)

	if !strings.HasPrefix(painted, backgroundEven) {
		t.Fatalf("painted = %q, want the background first", painted)
	}
	if !strings.HasSuffix(painted, "\x1b[0m") {
		t.Fatalf("painted = %q, want a reset at the end", painted)
	}
	head := ansi.Truncate(painted, 5, "")
	if got := ansi.StringWidth(head); got != 5 {
		t.Fatalf("painted width = %d, want 5", got)
	}
}

func TestPaintBackgroundKeepsTheRest(t *testing.T) {
	painted := paintBackground(" abcdmeta", backgroundOdd, 3)

	if !strings.HasSuffix(painted, "\x1b[0mcdmeta") {
		t.Fatalf("the meta column was lost: %q", painted)
	}
	if got := ansi.StringWidth(painted); got != 9 {
		t.Fatalf("painted width = %d, want 9", got)
	}
}

func TestPaintBackgroundReappliesAfterInnerColour(t *testing.T) {
	line := " \x1b[38;5;45mhello\x1b[0m"
	painted := paintBackground(line, backgroundCursor, 7)

	if !strings.Contains(painted, "\x1b[0m"+backgroundCursor) {
		t.Fatalf("the background was not restored after the inner reset: %q", painted)
	}
}

func TestPaintBackgroundIgnoresNonPositiveWidth(t *testing.T) {
	if got := paintBackground(" a", backgroundEven, 0); got != " a" {
		t.Fatalf("painted = %q, want the line untouched", got)
	}
}

func TestDetailAlternatesMessageBackgrounds(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	view := m.View()

	if !strings.Contains(view, backgroundEven) {
		t.Fatal("the even message is not tinted")
	}
	if !strings.Contains(view, backgroundOdd) {
		t.Fatal("the odd message is not tinted")
	}
}

func TestDetailHeadStaysUntinted(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	lines := strings.Split(m.View(), "\n")

	if strings.Contains(lines[0], "\x1b[48;5;") {
		t.Fatalf("the thread title is tinted: %q", lines[0])
	}
}

func TestSelectTintsTheCursorAndTheMarks(t *testing.T) {
	m := selectMode(t, sampleConversation())
	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')
	view := m.View()

	if !strings.Contains(view, backgroundCursor) {
		t.Fatal("the cursor message is not tinted")
	}
	if !strings.Contains(view, backgroundMarked) {
		t.Fatal("the marked message is not tinted")
	}
}

func TestMarkingDoesNotRebuildTheThread(t *testing.T) {
	m := selectMode(t, longConversation())
	builds := m.builds

	m, _ = pressRune(t, m, ' ')
	m, _ = pressRune(t, m, 'j')
	m, _ = pressRune(t, m, 'j')

	if m.builds != builds {
		t.Fatalf("builds = %d, want %d", m.builds, builds)
	}
}

func TestZebraKeepsTheFrameWidth(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	for i, line := range strings.Split(m.View(), "\n") {
		if got := ansi.StringWidth(line); got > 120 {
			t.Fatalf("line %d is %d cells wide", i, got)
		}
	}
}
