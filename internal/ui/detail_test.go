package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func detailModel(conv atuin.Conversation, width int) Model {
	return Model{width: width, height: 36, mode: modeDetail, conversation: conv}
}

func frameLines(conv atuin.Conversation, width int) []string {
	return detailModel(conv, width).detailLines()
}

func message(role string, blocks ...atuin.Content) atuin.Message {
	return atuin.Message{
		Role:      role,
		Timestamp: time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC),
		Content:   blocks,
	}
}

func TestDetailLayoutWidths(t *testing.T) {
	cases := []struct {
		width int
		body  int
		prose int
		meta  int
	}{
		{170, 142, 80, 24},
		{120, 92, 80, 24},
		{80, 52, 52, 24},
		{48, 20, 20, 24},
		{47, 45, 45, 0},
		{40, 38, 38, 0},
		{2, 1, 1, 0},
	}

	for _, tc := range cases {
		got := newDetailLayout(tc.width)
		if got.body != tc.body || got.prose != tc.prose || got.meta != tc.meta {
			t.Errorf("layout(%d) = %+v, want body %d prose %d meta %d",
				tc.width, got, tc.body, tc.prose, tc.meta)
		}
	}
}

func TestDetailLayoutSurvivesTinyWidth(t *testing.T) {
	conv := atuin.Conversation{
		Messages: []atuin.Message{message("user", atuin.Content{Type: "text", Text: "hello"})},
	}
	if lines := frameLines(conv, 2); len(lines) == 0 {
		t.Fatal("no lines at width 2")
	}
}

func TestDetailFrameGeometry(t *testing.T) {
	conv := atuin.Conversation{
		Session: atuin.Thread{Title: "Geometry", CWD: "/home/pim/work"},
		Messages: []atuin.Message{
			message("user", atuin.Content{Type: "text", Text: "Please bootstrap this repository."}),
			message("assistant", atuin.Content{Type: "text", Text: "on it"}),
		},
	}

	cases := []struct {
		width int
		rail  int
	}{
		{170, 145},
		{120, 95},
		{80, 55},
		{48, 23},
		{47, -1},
	}

	for _, tc := range cases {
		lines := frameLines(conv, tc.width)
		layout := newDetailLayout(tc.width)

		for _, line := range lines {
			text := plain(line)
			if text == "" {
				continue
			}
			if !strings.HasPrefix(text, " ") {
				t.Errorf("width %d: content row without a left margin: %q", tc.width, text)
			}
			if got := ansi.StringWidth(text); got > tc.width-1 {
				t.Errorf("width %d: row uses %d cells, want at most %d: %q",
					tc.width, got, tc.width-1, text)
			}
		}

		if tc.rail < 0 {
			if frame := plain(strings.Join(lines, "\n")); strings.Contains(frame, "12:00") {
				t.Errorf("width %d: meta survived without a rail:\n%s", tc.width, frame)
			}
			continue
		}

		row := -1
		for i, line := range lines {
			if strings.Contains(plain(line), "12:00 - user") {
				row = i
				break
			}
		}
		if row < 0 {
			t.Fatalf("width %d: no meta row in frame:\n%s", tc.width, plain(strings.Join(lines, "\n")))
		}

		if at := strings.Index(plain(lines[row]), "12:00 - user"); at != tc.rail {
			t.Errorf("width %d: meta at %d, want %d: %q", tc.width, at, tc.rail, plain(lines[row]))
		}
		if at := strings.Index(plain(lines[row+1]), "text"); at != tc.rail {
			t.Errorf("width %d: kind at %d, want %d: %q", tc.width, at, tc.rail, plain(lines[row+1]))
		}
		if body := plain(lines[row]); !strings.HasPrefix(body, " Please bootstrap") {
			t.Errorf("width %d: body does not start after the margin: %q", tc.width, body)
		}
		if layout.meta != 0 && layout.body != tc.rail-1-detailGap {
			t.Errorf("width %d: body %d does not meet the rail at %d", tc.width, layout.body, tc.rail)
		}
	}
}

func TestDetailFramePutsMetaRight(t *testing.T) {
	conv := atuin.Conversation{
		Session: atuin.Thread{
			Title:        "Frame",
			Harness:      "opencode",
			Model:        "deepseek-v4.1-flash",
			CWD:          "/home/pim/work",
			AtuinID:      "01a1d43c1f237541a0e1c6c42cf53ec4",
			MessageCount: 2,
		},
		Messages: []atuin.Message{
			message("user", atuin.Content{Type: "text", Text: "Please bootstrap this repository."}),
			message("assistant", atuin.Content{Type: "text", Text: "on it"}),
		},
	}

	lines := frameLines(conv, 120)

	row := -1
	for i, line := range lines {
		if strings.Contains(plain(line), "12:00 - user") {
			row = i
			break
		}
	}
	if row < 0 {
		t.Fatalf("no user meta in frame:\n%s", plain(strings.Join(lines, "\n")))
	}

	first := plain(lines[row])
	if at := strings.Index(first, "12:00 - user"); at != 95 {
		t.Fatalf("meta starts at %d, want 95: %q", at, first)
	}
	if !strings.HasPrefix(first, " Please bootstrap this repository.") {
		t.Errorf("body does not start after the margin: %q", first)
	}

	kind := plain(lines[row+1])
	if at := strings.Index(kind, "text"); at != 95 {
		t.Errorf("kind starts at %d, want 95: %q", at, kind)
	}

	if lines[row+2] != "" {
		t.Errorf("no blank line between messages: %q", plain(lines[row+2]))
	}

	second := plain(lines[row+3])
	if at := strings.Index(second, "12:00 - asst"); at != 95 {
		t.Errorf("second meta at %d, want 95: %q", at, second)
	}
}

func TestDetailDropsMetaWhenNarrow(t *testing.T) {
	conv := atuin.Conversation{
		Messages: []atuin.Message{message("user", atuin.Content{Type: "text", Text: "short"})},
	}

	for _, line := range frameLines(conv, 45) {
		if strings.Contains(plain(line), "12:00") {
			t.Fatalf("meta survived at 45 columns: %q", plain(line))
		}
	}
}

func TestDetailProseKeepsItsMeasure(t *testing.T) {
	text := strings.Repeat("word ", 40)
	block := atuin.Content{Type: "text", Text: text}

	for _, width := range []int{160, 120, 80} {
		layout := newDetailLayout(width)
		for _, line := range blockLines(block, layout.prose, layout.body) {
			if got := ansi.StringWidth(line); got > layout.prose {
				t.Fatalf("width %d: prose line is %d cells, want at most %d: %q",
					width, got, layout.prose, line)
			}
		}
	}

	if got := newDetailLayout(80).prose; got != 52 {
		t.Errorf("prose at 80 columns = %d, want 52", got)
	}
}

func TestDetailCodeUsesTheBodyColumn(t *testing.T) {
	command := "echo " + strings.Repeat("x", 75)
	conv := atuin.Conversation{Messages: []atuin.Message{message("assistant",
		atuin.Content{Type: "tool_call", Name: "shell", Input: `{"command":"` + command + `"}`})}}

	frame := plain(strings.Join(frameLines(conv, 120), "\n"))
	if strings.Contains(frame, "…") {
		t.Errorf("80-cell command was truncated at 120 columns:\n%s", frame)
	}
	if !strings.Contains(frame, "$ echo "+strings.Repeat("x", 75)) {
		t.Errorf("80-cell command did not stay on one row:\n%s", frame)
	}
}

func TestMessageMeta(t *testing.T) {
	cases := []struct {
		name string
		msg  atuin.Message
		want string
	}{
		{
			"thinking uses reasoning",
			atuin.Message{Role: "assistant", Tokens: &atuin.Tokens{Output: 265, Reasoning: 83},
				Content: []atuin.Content{{Type: "thinking", Text: "hmm"}}},
			"think - 83 tokens",
		},
		{
			"text uses output",
			atuin.Message{Role: "assistant", Tokens: &atuin.Tokens{Output: 100},
				Content: []atuin.Content{{Type: "text", Text: "hi"}}},
			"text - 100 tokens",
		},
		{
			"several kinds join",
			atuin.Message{Role: "assistant", Tokens: &atuin.Tokens{Output: 265},
				Content: []atuin.Content{{Type: "thinking", Text: "a"}, {Type: "tool_call", Name: "shell", Input: `{"command":"ls"}`}}},
			"think+tool - 265 tokens",
		},
		{
			"no usage",
			atuin.Message{Role: "tool", Content: []atuin.Content{{Type: "tool_result", Content: "ok"}}},
			"out",
		},
		{
			"no blocks",
			atuin.Message{Role: "attachment"},
			"empty",
		},
		{
			"repeated kinds collapse",
			atuin.Message{Role: "assistant", Content: []atuin.Content{{Type: "text", Text: "a"}, {Type: "text", Text: "b"}}},
			"text",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := secondMeta(tc.msg); got != tc.want {
				t.Errorf("secondMeta = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAuthorLabels(t *testing.T) {
	cases := map[string]string{
		"user":             "user",
		"assistant":        "asst",
		"tool":             "tool",
		"system":           "sys",
		"attachment":       "att",
		"ai-title":         "title",
		"away_summary":     "away",
		"":                 "?",
		"custom-very-long": "custo",
	}

	for role, want := range cases {
		if got := authorLabel(role); got != want {
			t.Errorf("authorLabel(%q) = %q, want %q", role, got, want)
		}
	}
}

func TestCollapsedRuns(t *testing.T) {
	conv := atuin.Conversation{
		Messages: []atuin.Message{
			message("attachment"),
			message("attachment"),
			message("attachment"),
			message("assistant", atuin.Content{Type: "text", Text: "after the noise"}),
		},
	}

	lines := frameLines(conv, 120)
	body := plain(strings.Join(lines, "\n"))
	if !strings.Contains(body, "(3 attachments)") {
		t.Fatalf("run not collapsed:\n%s", body)
	}
	if strings.Contains(body, "(1 attachment)") {
		t.Errorf("run counted wrong:\n%s", body)
	}
	if strings.Count(body, "12:00 - att") != 1 {
		t.Errorf("collapsed run should print one meta line:\n%s", body)
	}
	if !strings.Contains(body, "after the noise") {
		t.Errorf("message after the run missing:\n%s", body)
	}
}

func TestCollapsedMixedRun(t *testing.T) {
	conv := atuin.Conversation{
		Messages: []atuin.Message{
			message("attachment"),
			message("ai-title"),
			message("system"),
		},
	}

	body := plain(strings.Join(frameLines(conv, 120), "\n"))
	if !strings.Contains(body, "(3 empty messages)") {
		t.Fatalf("mixed run not collapsed:\n%s", body)
	}
}

func TestSingleEmptyMessageStays(t *testing.T) {
	conv := atuin.Conversation{
		Messages: []atuin.Message{
			message("attachment"),
			message("user", atuin.Content{Type: "text", Text: "hello"}),
		},
	}

	body := plain(strings.Join(frameLines(conv, 120), "\n"))
	if strings.Contains(body, "(1 attachments)") {
		t.Errorf("single empty message collapsed:\n%s", body)
	}
	if !strings.Contains(body, "12:00 - att") {
		t.Errorf("single empty message missing:\n%s", body)
	}
}

func TestToolCallHasNoJSON(t *testing.T) {
	input := `{"command":"cd /repo; ls scripts\ncat dev-fe.sh","description":"Inspect dev scripts"}`
	lines := toolCallLines(atuin.Content{Type: "tool_call", Name: "Bash", Input: input}, 80)

	if lines[0] != "Bash  Inspect dev scripts" {
		t.Errorf("head = %q", lines[0])
	}
	if lines[1] != "  $ cd /repo; ls scripts" {
		t.Errorf("command line = %q", lines[1])
	}

	joined := strings.Join(lines, "\n")
	for _, forbidden := range []string{"{", "}", `\"`, `"`} {
		if strings.Contains(joined, forbidden) {
			t.Errorf("rendered tool call contains %q:\n%s", forbidden, joined)
		}
	}
}

func TestToolCallShapes(t *testing.T) {
	cases := []struct {
		name  string
		block atuin.Content
		want  []string
	}{
		{
			"write",
			atuin.Content{Name: "write", Input: `{"path":"a.txt","content":"one\ntwo\nthree"}`},
			[]string{"write a.txt  +3 lines"},
		},
		{
			"edit",
			atuin.Content{Name: "edit", Input: `{"path":"b.go","oldString":"one","newString":"two"}`},
			[]string{"edit b.go  +1 -1", "  - one", "  + two"},
		},
		{
			"read",
			atuin.Content{Name: "Read", Input: `{"file_path":"c.md","offset":40,"limit":40}`},
			[]string{"Read c.md  40..80"},
		},
		{
			"skill",
			atuin.Content{Name: "Skill", Input: `{"skill":"docs","args":"guide"}`},
			[]string{"Skill docs"},
		},
		{
			"compress",
			atuin.Content{Name: "compress", Input: `{"topic":"Old context"}`},
			[]string{"compress Old context"},
		},
		{
			"unknown",
			atuin.Content{Name: "webfetch", Input: `{"url":"https://example.com"}`},
			[]string{"webfetch  url=https://example.com"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.block.Type = "tool_call"
			lines := toolCallLines(tc.block, 80)
			for _, want := range tc.want {
				if !strings.Contains(strings.Join(lines, "\n"), want) {
					t.Errorf("lines %q missing %q", lines, want)
				}
			}
		})
	}
}

func TestToolPayloadIsCapped(t *testing.T) {
	input := `{"command":"` + strings.Repeat("line\n", 20) + `"}`
	lines := toolCallLines(atuin.Content{Type: "tool_call", Name: "shell", Input: input}, 80)
	if len(lines) != maxBlockLines+1 {
		t.Fatalf("lines = %d, want %d", len(lines), maxBlockLines+1)
	}
	if !strings.Contains(plain(lines[len(lines)-1]), "more lines") {
		t.Errorf("missing truncation note: %q", plain(lines[len(lines)-1]))
	}
}

func TestResultKeepsIndentation(t *testing.T) {
	content := "context: |\n  Quiqr ChatCMS is a monorepo\n    nested: true"
	lines := resultLines(atuin.Content{Type: "tool_result", Content: content}, 40)

	want := []string{"context: |", "  Quiqr ChatCMS is a monorepo", "    nested: true"}
	if len(lines) != len(want) {
		t.Fatalf("lines = %q", lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d = %q, want %q", i, lines[i], w)
		}
	}
}

func TestResultWrapsAtBodyWidth(t *testing.T) {
	content := strings.Repeat("x", 60)
	lines := resultLines(atuin.Content{Type: "tool_result", Content: content}, 20)

	for _, line := range lines {
		if width := ansi.StringWidth(line); width > 20 {
			t.Errorf("line is %d cells wide, want at most 20: %q", width, line)
		}
	}
	if len(lines) < 2 {
		t.Errorf("long line not wrapped: %q", lines)
	}
}

func TestPatchBlock(t *testing.T) {
	block := atuin.Content{
		Type: "patch",
		Files: []atuin.PatchFile{{
			Path:      "flake.nix",
			Change:    "update",
			Additions: 2,
			Deletions: 2,
			Diff:      "@@ -11,13 +11,13 @@\n-pname = \"atuin-threads-tui\";\n+pname = \"att\";",
		}},
	}

	lines := patchLines(block, 80)
	body := strings.Join(lines, "\n")
	if !strings.Contains(body, "flake.nix  update  +2 -2") {
		t.Errorf("head missing:\n%s", body)
	}
	if !strings.Contains(body, "  -pname") {
		t.Errorf("diff missing:\n%s", body)
	}
}

func TestSummaryBlockRenders(t *testing.T) {
	lines := blockLines(atuin.Content{Type: "summary", Text: "Summarizing 18 earlier messages"}, 80, 92)
	if len(lines) != 1 || lines[0] != "Summarizing 18 earlier messages" {
		t.Errorf("summary = %q", lines)
	}
}

func TestProseWrapsByDisplayWidth(t *testing.T) {
	text := strings.Repeat("日本語のプロジェクト ", 6)
	for _, width := range []int{20, 33, 54} {
		for _, line := range wrap(text, width) {
			if got := ansi.StringWidth(line); got > width {
				t.Errorf("width %d: line is %d cells: %q", width, got, line)
			}
		}
	}
}

func TestProseCombiningMarks(t *testing.T) {
	text := strings.Repeat("é", 10) + " e\u0301"
	for _, line := range wrap(text, 8) {
		if got := ansi.StringWidth(line); got > 8 {
			t.Errorf("line is %d cells: %q", got, line)
		}
	}
}

func TestPlaceholderThinkingIsHidden(t *testing.T) {
	if lines := blockLines(atuin.Content{Type: "thinking", Text: "Reasoning · 83 tokens"}, 80, 92); lines != nil {
		t.Errorf("placeholder thinking rendered: %q", lines)
	}
	if lines := blockLines(atuin.Content{Type: "thinking", Text: "Reasoning about the tokens"}, 80, 92); len(lines) == 0 {
		t.Error("real thinking text was hidden")
	}
}

func TestThinkingOnlyMessagesDoNotCollapse(t *testing.T) {
	usage := &atuin.Tokens{Output: 265, Reasoning: 83}
	conv := atuin.Conversation{
		Messages: []atuin.Message{
			{Role: "assistant", Tokens: usage, Content: []atuin.Content{{Type: "thinking", Text: "Reasoning · 83 tokens"}}},
			{Role: "assistant", Tokens: usage, Content: []atuin.Content{{Type: "thinking", Text: "Reasoning · 83 tokens"}}},
			{Role: "assistant", Tokens: usage, Content: []atuin.Content{{Type: "thinking", Text: "Reasoning · 83 tokens"}}},
		},
	}

	body := plain(strings.Join(frameLines(conv, 120), "\n"))
	if strings.Contains(body, "empty messages") {
		t.Errorf("thinking messages collapsed:\n%s", body)
	}
	if got := strings.Count(body, "think - 83 tokens"); got != 3 {
		t.Errorf("meta lines = %d, want 3:\n%s", got, body)
	}
}
