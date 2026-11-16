package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func TestCategoryOf(t *testing.T) {
	cases := []struct {
		name  string
		block atuin.Content
		role  string
		want  typeMask
	}{
		{"user text", atuin.Content{Type: "text"}, "user", typePrompts},
		{"assistant text", atuin.Content{Type: "text"}, "assistant", typeReplies},
		{"tool text", atuin.Content{Type: "text"}, "tool", typeReplies},
		{"thinking", atuin.Content{Type: "thinking"}, "assistant", typeThinking},
		{"tool call", atuin.Content{Type: "tool_call"}, "assistant", typeCalls},
		{"tool result", atuin.Content{Type: "tool_result"}, "tool", typeResults},
		{"error result", atuin.Content{Type: "error"}, "assistant", typeResults},
		{"patch", atuin.Content{Type: "patch"}, "assistant", typePatches},
		{"summary", atuin.Content{Type: "summary"}, "away_summary", typeSummaries},
		{"unknown", atuin.Content{Type: "image"}, "assistant", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := categoryOf(tc.block, tc.role); got != tc.want {
				t.Errorf("categoryOf = %b, want %b", got, tc.want)
			}
		})
	}
}

func TestUnknownCategoryStaysVisible(t *testing.T) {
	if !allTypes.shows(0) {
		t.Error("a block of unknown kind is hidden")
	}
	if allTypes.shows(typeThinking) {
		t.Error("hiding everything still shows thinking")
	}
}

func typeConversation() atuin.Conversation {
	now := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	return atuin.Conversation{
		Session: atuin.Thread{
			AtuinID:      "01a1d43c1f237541a0e1c6c42cf53ec4",
			Title:        "Type filter",
			MessageCount: 5,
			UpdatedAt:    now,
		},
		Messages: []atuin.Message{
			{Role: "user", Timestamp: now, Content: []atuin.Content{{Type: "text", Text: "Please bootstrap"}}},
			{Role: "assistant", Timestamp: now, Content: []atuin.Content{
				{Type: "thinking", Text: "Considering the layout"},
				{Type: "text", Text: "hello there"},
				{Type: "tool_call", Name: "shell", Input: `{"command":"ls"}`},
				{Type: "tool_result", Content: "README.md"},
			}},
			{Role: "assistant", Timestamp: now, Content: []atuin.Content{
				{Type: "tool_call", Name: "shell", Input: `{"command":"pwd"}`},
			}},
			{Role: "assistant", Timestamp: now, Content: []atuin.Content{
				{Type: "thinking", Text: "Reasoning · 12 tokens"},
			}},
			{Role: "attachment", Timestamp: now},
			{Role: "attachment", Timestamp: now},
			{Role: "attachment", Timestamp: now},
		},
	}
}

func rebuilt(t *testing.T, m Model) Model {
	t.Helper()
	updated, cmd := m.startDetail()
	model := updated.(Model)
	for cmd != nil {
		model, cmd = feed(t, model, cmd())
	}
	return model
}

func renderDetail(m Model) string {
	lines := m.buildDetail().all()
	return plain(strings.Join(lines, "\n"))
}

func TestFilteringDropsHiddenBlocks(t *testing.T) {
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: typeConversation(), hidden: typeThinking, now: time.Now}

	body := renderDetail(m)
	for _, hidden := range []string{"Considering the layout", "Reasoning · 12 tokens"} {
		if strings.Contains(body, hidden) {
			t.Errorf("body still holds hidden thinking text %q", hidden)
		}
	}
	for _, shown := range []string{"Please bootstrap", "hello there", "ls", "README.md"} {
		if !strings.Contains(body, shown) {
			t.Errorf("body lost visible content %q", shown)
		}
	}
}

func TestFilteringKeepsTheVisiblePartOfAMessage(t *testing.T) {
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: typeConversation(), hidden: typeThinking, now: time.Now}

	body := renderDetail(m)
	if !strings.Contains(body, "hello there") {
		t.Fatal("the reply of a message with hidden thinking disappeared")
	}
	if !strings.Contains(body, "text+tool+out") {
		t.Errorf("meta does not name the visible kinds:\n%s", body)
	}
}

func TestFilteringSkipsMessagesWithoutVisibleBlocks(t *testing.T) {
	conv := typeConversation()
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: conv, hidden: typeCalls, now: time.Now}

	rendered := m.buildDetail()
	if rendered.messages != 3 {
		t.Errorf("rendered messages = %d, want 3", rendered.messages)
	}
	body := plain(strings.Join(rendered.all(), "\n"))
	if strings.Contains(body, "pwd") {
		t.Error("a tool-only message is still rendered")
	}
}

func TestNotesHideEmptyMessages(t *testing.T) {
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: typeConversation(), hidden: typeNotes, now: time.Now}

	body := renderDetail(m)
	if strings.Contains(body, "attachments") {
		t.Errorf("empty messages are still shown:\n%s", body)
	}
}

func TestNotesShowAsOneRow(t *testing.T) {
	conv := typeConversation()
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: conv, now: time.Now}

	body := renderDetail(m)
	if !strings.Contains(body, "(3 attachments)") {
		t.Errorf("the collapsed run is missing:\n%s", body)
	}
}

func TestEverythingHiddenShowsANotice(t *testing.T) {
	m := Model{mode: modeDetail, width: 100, height: 24, conversation: typeConversation(), hidden: allTypes, now: time.Now}

	body := renderDetail(m)
	if !strings.Contains(body, "no messages of the selected types") {
		t.Errorf("no empty notice:\n%s", body)
	}
	if m.buildDetail().messages != 0 {
		t.Error("messages were counted while everything is hidden")
	}
}

func TestStatusCountsShownAndTotal(t *testing.T) {
	m := Model{mode: modeDetail, width: 120, height: 24, conversation: typeConversation(), now: time.Now}
	m, _ = feed(t, m, conversationMsg{conversation: typeConversation()})

	status := plain(m.status())
	if !strings.Contains(status, "5 msgs") || strings.Contains(status, " of 5 msgs") {
		t.Errorf("unfiltered status = %q", status)
	}

	m.hidden = typePrompts
	m = rebuilt(t, m)
	status = plain(m.status())
	if !strings.Contains(status, "3 of 5 msgs") {
		t.Errorf("filtered status = %q", status)
	}
}

func TestTypePanelOpensWithCounts(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	m, cmd := pressRune(t, m, 't')
	if cmd != nil {
		t.Fatal("opening the panel returned a command")
	}
	if m.mode != modeTypes {
		t.Fatalf("mode = %v, want modeTypes", m.mode)
	}
	if m.typesDraft != m.hidden {
		t.Error("the draft does not start from the selection in use")
	}
	if got := m.typesCounts[typePrompts]; got != 1 {
		t.Errorf("prompt count = %d, want 1", got)
	}
	if got := m.typesCounts[typeNotes]; got != 3 {
		t.Errorf("notes count = %d, want 3", got)
	}

	body := plain(m.View())
	for _, want := range []string{"show message types", "[x] prompts", "[x] thinking", "[x] notes", "your messages", "5 msgs"} {
		if !strings.Contains(body, want) {
			t.Errorf("panel missing %q:\n%s", want, body)
		}
	}
}

func TestTypePanelToggleAndApply(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	m, _ = pressRune(t, m, 't')

	m, _ = press(t, m, tea.KeyDown)
	m, _ = press(t, m, tea.KeyDown)
	if m.typesCursor != 2 {
		t.Fatalf("cursor = %d, want 2", m.typesCursor)
	}
	m, _ = press(t, m, tea.KeySpace)
	if m.typesDraft.shows(typeThinking) {
		t.Error("space did not hide the marked row")
	}
	if m.hidden != 0 {
		t.Error("the selection was applied before enter")
	}

	builds := m.builds
	m, _ = press(t, m, tea.KeyEnter)
	if m.mode != modeDetail {
		t.Errorf("mode = %v, want modeDetail", m.mode)
	}
	if m.hidden != typeThinking {
		t.Errorf("hidden = %b, want the thinking bit", m.hidden)
	}
	if m.builds != builds+1 {
		t.Errorf("builds = %d, want one rebuild", m.builds-builds)
	}
}

func TestTypePanelCancelKeepsSelection(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	m, _ = pressRune(t, m, 't')
	m, _ = press(t, m, tea.KeySpace)
	m, _ = press(t, m, tea.KeyEsc)

	if m.mode != modeDetail {
		t.Errorf("mode = %v, want modeDetail", m.mode)
	}
	if m.hidden != 0 {
		t.Error("esc applied the draft")
	}
}

func TestTypePanelPresets(t *testing.T) {
	cases := []struct {
		key  rune
		want typeMask
	}{
		{'1', 0},
		{'2', allTypes &^ conversationTypes},
		{'3', allTypes &^ toolTypes},
		{'4', allTypes &^ promptTypes},
		{'a', 0},
		{'n', allTypes &^ conversationTypes},
	}

	for _, tc := range cases {
		t.Run(string(tc.key), func(t *testing.T) {
			m, _ := openDetail(t, typeConversation(), 120, 24)
			m, _ = pressRune(t, m, 't')
			m, _ = press(t, m, tea.KeySpace)
			m, _ = pressRune(t, m, tc.key)
			if m.typesDraft != tc.want {
				t.Errorf("draft = %b, want %b", m.typesDraft, tc.want)
			}
		})
	}
}

func TestSelectionOutlivesTheThread(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	m, _ = pressRune(t, m, 't')
	m, _ = press(t, m, tea.KeySpace)
	m, _ = press(t, m, tea.KeyEnter)

	m, _ = feed(t, m, conversationMsg{conversation: longConversation()})
	if m.hidden != typePrompts {
		t.Errorf("hidden = %b, want the prompts bit after opening another thread", m.hidden)
	}
}

func TestPanelCursorStaysInBounds(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	m, _ = pressRune(t, m, 't')

	for i := 0; i < 20; i++ {
		m, _ = press(t, m, tea.KeyDown)
	}
	if m.typesCursor != len(typeRows)-1 {
		t.Errorf("cursor = %d, want %d", m.typesCursor, len(typeRows)-1)
	}
	m, _ = press(t, m, tea.KeyHome)
	if m.typesCursor != 0 {
		t.Errorf("cursor = %d, want 0", m.typesCursor)
	}
	m, _ = press(t, m, tea.KeyUp)
	if m.typesCursor != 0 {
		t.Errorf("cursor = %d, want 0 above the first row", m.typesCursor)
	}
}

func TestHintsMentionTypes(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	if hints := plain(m.hints()); !strings.Contains(hints, "t types") {
		t.Errorf("detail hints = %q", hints)
	}

	m, _ = pressRune(t, m, 't')
	for _, want := range []string{"space toggle", "enter apply", "esc cancel"} {
		if hints := plain(m.hints()); !strings.Contains(hints, want) {
			t.Errorf("panel hints missing %q: %q", want, hints)
		}
	}
}

func TestPanelIsNotShownInTheList(t *testing.T) {
	m := loaded(t, &fakeSource{threads: sampleThreads(), conversation: typeConversation()})
	if view := plain(m.View()); strings.Contains(view, "show message types") {
		t.Error("the panel is drawn in the list")
	}
}

func TestApplyingASelectionKeepsTheTopMessage(t *testing.T) {
	conv := sampleConversation()
	for i := 0; i < 40; i++ {
		conv.Messages = append(conv.Messages, message("assistant",
			atuin.Content{Type: "thinking", Text: fmt.Sprintf("thinking %d", i)},
			atuin.Content{Type: "text", Text: fmt.Sprintf("reply %d", i)}))
	}

	m, _ := openDetail(t, conv, 120, 24)
	m.detailScroll = 25
	anchor := m.detail.messageAt(m.detailScroll)
	if anchor <= 0 {
		t.Fatalf("anchor = %d, want a message", anchor)
	}

	m.hidden = typeThinking
	m = rebuilt(t, m)

	if m.detailScroll > 25 {
		t.Errorf("scroll = %d, want it to move up as thinking lines are gone", m.detailScroll)
	}
	if got := m.detail.messageAt(m.detailScroll); got != anchor {
		t.Errorf("top message = %d, want %d", got, anchor)
	}
}

func TestHiddenAnchorFallsThrough(t *testing.T) {
	now := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	conv := atuin.Conversation{Session: atuin.Thread{AtuinID: "01a1d43c1f237541a0e1c6c42cf53ec4", MessageCount: 4}}
	for i := 0; i < 3; i++ {
		conv.Messages = append(conv.Messages, atuin.Message{Role: "assistant", Timestamp: now,
			Content: []atuin.Content{{Type: "text", Text: fmt.Sprintf("reply %d", i)}}})
	}
	conv.Messages = append(conv.Messages, atuin.Message{Role: "assistant", Timestamp: now,
		Content: []atuin.Content{{Type: "tool_call", Name: "shell", Input: `{"command":"ls"}`}}})

	m, _ := openDetail(t, conv, 120, 24)
	m.detailScroll = m.detail.blocks[3].start
	anchor := m.detail.messageAt(m.detailScroll)
	if anchor != 3 {
		t.Fatalf("anchor = %d, want the tool message", anchor)
	}

	m.hidden = typeCalls
	m = rebuilt(t, m)

	if m.detailScroll != 0 {
		t.Errorf("scroll = %d, want the top of the thread when nothing follows the anchor", m.detailScroll)
	}
}

func TestStaleMaskRenderingIsDropped(t *testing.T) {
	m, _ := openDetail(t, typeConversation(), 120, 24)
	fresh := m.detail

	m.rendering = true
	m.detail = detailRendering{}
	m.hidden = typeThinking

	stale := fresh
	m, _ = feed(t, m, detailRenderedMsg{rendering: stale})
	if !m.rendering || m.detail.built {
		t.Error("a rendering made for another selection was accepted")
	}

	m, _ = feed(t, m, detailRenderedMsg{rendering: stale, anchor: -1})
	if !m.rendering {
		t.Error("the stale rendering was accepted twice")
	}

	stale.hidden = typeThinking
	m, _ = feed(t, m, detailRenderedMsg{rendering: stale})
	if m.rendering || !m.detail.built || m.detail.hidden != typeThinking {
		t.Error("the rendering for the current selection was not accepted")
	}
}
