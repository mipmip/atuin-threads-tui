package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func openDetail(t *testing.T, conv atuin.Conversation, width, height int) (Model, tea.Cmd) {
	t.Helper()
	m := New(nil)
	m.width, m.height = width, height
	return feed(t, m, conversationMsg{conversation: conv})
}

func bulkConversation(n int) atuin.Conversation {
	conv := sampleConversation()
	for i := 0; i < n; i++ {
		conv.Messages = append(conv.Messages, message("assistant",
			atuin.Content{Type: "text", Text: fmt.Sprintf("message %d", i)}))
	}
	return conv
}

func TestDetailRendersOnce(t *testing.T) {
	m, _ := openDetail(t, longConversation(), 120, 24)
	if !m.detail.built {
		t.Fatal("the conversation was not rendered on open")
	}
	if m.builds != 1 {
		t.Fatalf("builds = %d after opening, want 1", m.builds)
	}

	total := m.detail.total
	for i := 0; i < 20; i++ {
		m, _ = pressRune(t, m, 'j')
	}

	if m.builds != 1 {
		t.Errorf("builds = %d after scrolling, want 1", m.builds)
	}
	if m.detail.total != total {
		t.Errorf("total = %d after scrolling, want %d", m.detail.total, total)
	}
	if m.detailScroll != 20 {
		t.Errorf("scroll = %d after 20 keys, want 20", m.detailScroll)
	}
}

func TestDetailSliceMatchesFullRendering(t *testing.T) {
	m, _ := openDetail(t, longConversation(), 120, 24)
	rendering := m.detail
	all := rendering.all()

	if len(all) != rendering.total {
		t.Fatalf("all() = %d lines, total = %d", len(all), rendering.total)
	}

	for _, from := range []int{0, 1, 2, 3, 7, 40, len(all) - 1, len(all)} {
		if from < 0 || from > len(all) {
			continue
		}
		for _, count := range []int{1, 5, 23, 200, len(all)} {
			got := rendering.slice(from, count)
			end := min(from+count, len(all))
			want := all[from:end]
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Fatalf("slice(%d, %d) differs from the full rendering", from, count)
			}
		}
	}
}

func TestDetailSliceEdges(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	rendering := m.detail

	if lines := rendering.slice(-1, 10); lines != nil {
		t.Errorf("slice with a negative start returned %d lines", len(lines))
	}
	if lines := rendering.slice(0, 0); lines != nil {
		t.Errorf("slice with no count returned %d lines", len(lines))
	}
	if lines := rendering.slice(rendering.total, 10); lines != nil {
		t.Errorf("slice past the end returned %d lines", len(lines))
	}
	if lines := rendering.slice(2, 3); len(lines) != 3 {
		t.Errorf("slice inside the view returned %d lines, want 3", len(lines))
	}

	first := rendering.blocks[0]
	if first.start != len(rendering.head) {
		t.Errorf("first block starts at %d, want %d", first.start, len(rendering.head))
	}
	for i := 1; i < len(rendering.blocks); i++ {
		prev := rendering.blocks[i-1]
		if rendering.blocks[i].start != prev.start+len(prev.lines) {
			t.Fatalf("block %d starts at %d, want %d", i, rendering.blocks[i].start, prev.start+len(prev.lines))
		}
	}
}

func TestDetailResizeRebuilds(t *testing.T) {
	m, _ := openDetail(t, longConversation(), 120, 24)
	if m.builds != 1 {
		t.Fatalf("builds = %d after opening, want 1", m.builds)
	}

	m, _ = feed(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	if m.builds != 2 {
		t.Errorf("builds = %d after resize, want 2", m.builds)
	}
	if m.detail.width != 80 {
		t.Errorf("rendering width = %d, want 80", m.detail.width)
	}

	m, _ = feed(t, m, tea.WindowSizeMsg{Width: 80, Height: 30})
	if m.builds != 2 {
		t.Errorf("builds = %d after a height-only resize, want 2", m.builds)
	}

	body := plain(strings.Join(m.detail.slice(0, 8), "\n"))
	if !strings.Contains(body, "Please bootstrap this repository.") {
		t.Errorf("body after resize does not show the message:\n%s", body)
	}
}

func TestDetailLargeThreadRendersInBackground(t *testing.T) {
	conv := bulkConversation(detailSyncMessages)
	m, cmd := openDetail(t, conv, 120, 24)

	if !m.rendering {
		t.Fatal("the model is not rendering a large thread")
	}
	if m.detail.built {
		t.Fatal("the large thread was rendered on the open path")
	}
	if cmd == nil {
		t.Fatal("no command was returned for the background render")
	}

	view := plain(m.View())
	if !strings.Contains(view, fmt.Sprintf("rendering %d messages", len(conv.Messages))) {
		t.Errorf("view does not show a rendering state:\n%s", view)
	}
	if !strings.Contains(view, "· rendering") {
		t.Errorf("status does not show a rendering state:\n%s", view)
	}

	m, _ = feed(t, m, cmd())
	if m.rendering {
		t.Error("the model is still rendering after the result arrived")
	}
	if !m.detail.built {
		t.Fatal("the rendering was not stored")
	}

	view = plain(m.View())
	if strings.Contains(view, "rendering") {
		t.Errorf("view still shows a rendering state:\n%s", view)
	}
	if !strings.Contains(view, fmt.Sprintf("line 1 of %d", m.detail.total)) {
		t.Errorf("status does not report the cached total:\n%s", view)
	}
}

func TestDetailStaleRenderingIsDropped(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	m.detail = detailRendering{}
	m.rendering = true

	stale := Model{mode: modeDetail, width: 100, conversation: m.conversation}.buildDetail()
	m, _ = feed(t, m, detailRenderedMsg{rendering: stale})

	if m.detail.built {
		t.Error("a rendering for another width was stored")
	}
	if !m.rendering {
		t.Error("the pending render was cleared by a stale result")
	}
}

func TestDetailRenderingNoteWithoutCache(t *testing.T) {
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: sampleConversation()}
	view := plain(m.View())
	if !strings.Contains(view, "rendering 2 messages") {
		t.Errorf("view does not show a rendering state without a cache:\n%s", view)
	}
}

func TestDetailCloseDropsRendering(t *testing.T) {
	m, _ := openDetail(t, sampleConversation(), 120, 24)
	if !m.detail.built {
		t.Fatal("the conversation was not rendered")
	}

	m, _ = press(t, m, tea.KeyEsc)
	if m.detail.built {
		t.Error("the rendering survived closing the thread")
	}
}

func TestDetailBatchedKeysScroll(t *testing.T) {
	m, _ := openDetail(t, longConversation(), 120, 24)
	m, _ = feed(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("jjjjj")})
	if m.detailScroll != 5 {
		t.Errorf("scroll = %d after five batched keys, want 5", m.detailScroll)
	}
}
