package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func TestCutWidth(t *testing.T) {
	cases := []struct {
		name  string
		in    string
		width int
		head  string
		rest  string
	}{
		{"fits", "hello", 10, "hello", ""},
		{"splits", "hello world", 5, "hello", " world"},
		{"zero width keeps one rune", "hello", 0, "h", "ello"},
		{"wide rune at width one", "日本", 1, "日", "本"},
		{"wide rune counts two cells", "日本語", 4, "日本", "語"},
		{"empty", "", 5, "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			head, rest := cutWidth(c.in, c.width)
			if head != c.head || rest != c.rest {
				t.Fatalf("cutWidth(%q, %d) = %q, %q; want %q, %q",
					c.in, c.width, head, rest, c.head, c.rest)
			}
		})
	}
}

func TestCutWidthAlwaysAdvances(t *testing.T) {
	for _, width := range []int{-5, 0, 1} {
		head, rest := cutWidth("abcdef", width)
		if head == "" || len(rest) != len("abcdef")-len(head) {
			t.Fatalf("width %d did not advance: %q / %q", width, head, rest)
		}
	}
}

func TestWrapPreservingLongLine(t *testing.T) {
	line := strings.Repeat("x", 100000)
	chunks := wrapPreserving(line, 80)

	want := 1 + (len(line)-80+77)/78
	if len(chunks) != want {
		t.Fatalf("got %d chunks, want %d", len(chunks), want)
	}
	if len(chunks[0]) != 80 {
		t.Fatalf("first chunk is %d cells, want 80", len(chunks[0]))
	}
	for i, chunk := range chunks[1:] {
		if !strings.HasPrefix(chunk, "  ") {
			t.Fatalf("chunk %d has no hang: %q", i+1, chunk)
		}
		if len(chunk) > 80 {
			t.Fatalf("chunk %d is %d cells, want at most 80", i+1, len(chunk))
		}
	}
	if last := chunks[len(chunks)-1]; len(last) > 80 {
		t.Fatalf("last chunk is %d cells", len(last))
	}
	var rebuilt strings.Builder
	for i, chunk := range chunks {
		if i > 0 {
			chunk = strings.TrimPrefix(chunk, "  ")
		}
		rebuilt.WriteString(chunk)
	}
	if rebuilt.String() != line {
		t.Fatalf("chunks do not reassemble the line")
	}
}

func TestWrapPreservingWideText(t *testing.T) {
	chunks := wrapPreserving(strings.Repeat("日本", 40), 10)
	for i, chunk := range chunks {
		if width := lipgloss.Width(chunk); width > 10 {
			t.Fatalf("chunk %d is %d cells: %q", i, width, chunk)
		}
		if chunk == "" {
			t.Fatalf("chunk %d is empty", i)
		}
	}
}

func TestSplitWidthLongLine(t *testing.T) {
	line := strings.Repeat("ab", 50000)
	chunks := splitWidth(line, 70)
	if len(chunks) != 1429 {
		t.Fatalf("got %d chunks, want 1429", len(chunks))
	}
	if strings.Join(chunks, "") != line {
		t.Fatalf("chunks do not reassemble the line")
	}
}

func TestResultWithHugeLineBuildsFast(t *testing.T) {
	conversation := atuin.Conversation{
		Session: atuin.Thread{Title: "big result", AtuinID: "abc"},
		Messages: []atuin.Message{{
			Role: "assistant",
			Content: []atuin.Content{{
				Type:    "tool_result",
				Content: strings.Repeat("y", 500000),
			}},
		}},
	}
	m := Model{mode: modeDetail, width: 120, height: 36, conversation: conversation, now: time.Now}

	start := time.Now()
	build := m.buildDetail()
	elapsed := time.Since(start)

	if !build.built {
		t.Fatal("conversation was not built")
	}
	if elapsed > 2*time.Second {
		t.Fatalf("building one huge line took %s", elapsed)
	}
	if build.total < 2 {
		t.Fatalf("got %d lines, want the title and one message", build.total)
	}
}
