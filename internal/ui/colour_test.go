package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func withColour(t *testing.T) {
	t.Helper()
	previous := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.ANSI256)
	t.Cleanup(func() { lipgloss.SetColorProfile(previous) })
}

func fg(code string) string {
	return "\x1b[38;5;" + code + "m"
}

func boldFg(code string) string {
	return "\x1b[1;38;5;" + code + "m"
}

func TestRoleProseColours(t *testing.T) {
	withColour(t)

	cases := []struct {
		role string
		want string
	}{
		{"user", fg("45")},
		{"assistant", fg("42")},
		{"tool", ""},
		{"system", ""},
		{"", ""},
	}

	for _, tc := range cases {
		got := styled(roleProseStyle(tc.role), "text")
		if tc.want == "" {
			if strings.Contains(got, "\x1b") {
				t.Errorf("role %q is styled: %q", tc.role, got)
			}
			continue
		}
		if !strings.HasPrefix(got, tc.want) {
			t.Errorf("role %q = %q, want prefix %q", tc.role, got, tc.want)
		}
	}
}

func TestProseColourCoversEveryWrappedLine(t *testing.T) {
	withColour(t)

	text := strings.Repeat("word ", 40)
	lines := blockLines(atuin.Content{Type: "text", Text: text}, 40, 40, roleProseStyle("user"))
	if len(lines) < 2 {
		t.Fatalf("expected wrapped lines, got %d", len(lines))
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, fg("45")) {
			t.Errorf("line is not prompt coloured: %q", line)
		}
		if width := ansi.StringWidth(plain(line)); width > 40 {
			t.Errorf("line is %d cells wide: %q", width, line)
		}
	}
}

func TestTextOfOtherRolesStaysUncoloured(t *testing.T) {
	withColour(t)

	lines := blockLines(atuin.Content{Type: "text", Text: "a result"}, 40, 40, roleProseStyle("tool"))
	for _, line := range lines {
		if strings.Contains(line, "\x1b") {
			t.Errorf("tool text is styled: %q", line)
		}
	}
}

func TestThinkingAndSummaryKeepTheirOwnColour(t *testing.T) {
	withColour(t)

	thinking := blockLines(atuin.Content{Type: "thinking", Text: "weighing options"}, 80, 80, roleProseStyle("assistant"))
	if len(thinking) == 0 {
		t.Fatal("no thinking lines")
	}
	if !strings.HasPrefix(thinking[0], fg("241")) {
		t.Errorf("thinking = %q, want muted", thinking[0])
	}
	if strings.HasPrefix(thinking[0], fg("42")) {
		t.Errorf("thinking took the reply colour: %q", thinking[0])
	}

	summary := blockLines(atuin.Content{Type: "summary", Text: "earlier messages"}, 80, 80, roleProseStyle("assistant"))
	for _, line := range summary {
		if strings.Contains(line, "\x1b") {
			t.Errorf("summary is styled: %q", line)
		}
	}
}

func TestStyleDiffLine(t *testing.T) {
	withColour(t)

	header := func(parts ...string) string { return "HDR:" + strings.Join(parts, "") }

	cases := []struct {
		line string
		want string
	}{
		{"@@ -1,2 +1,3 @@", fg("45")},
		{"+++ b/main.go", "HDR:"},
		{"--- a/main.go", "HDR:"},
		{"+added", fg("42")},
		{"-removed", fg("196")},
		{" context", " context"},
		{"", ""},
	}

	for _, tc := range cases {
		got := styleDiffLine(tc.line, header)
		if !strings.HasPrefix(got, tc.want) {
			t.Errorf("styleDiffLine(%q) = %q, want prefix %q", tc.line, got, tc.want)
		}
	}
}

func TestPatchLinesColourTheDiff(t *testing.T) {
	withColour(t)

	block := atuin.Content{Type: "patch", Files: []atuin.PatchFile{{
		Path:      "internal/ui/detail.go",
		Change:    "modified",
		Additions: 1,
		Deletions: 1,
		Diff:      "@@ -1 +1 @@\n-old line\n+new line",
	}}}
	lines := patchLines(block, 80)

	want := map[string]string{
		"  @@ -1 +1 @@": fg("45"),
		"  -old line":   fg("196"),
		"  +new line":   fg("42"),
	}
	seen := map[string]bool{}
	for _, line := range lines {
		text := plain(line)
		if prefix, ok := want[text]; ok {
			seen[text] = true
			if !strings.Contains(line, prefix) {
				t.Errorf("line %q lacks %q", line, prefix)
			}
		}
	}
	for text := range want {
		if !seen[text] {
			t.Errorf("missing diff line %q in %q", text, lines)
		}
	}
	if !strings.Contains(lines[0], fg("42")) || !strings.Contains(lines[0], fg("196")) {
		t.Errorf("head misses the change counts: %q", lines[0])
	}
}

func TestEditLinesColourOldAndNew(t *testing.T) {
	withColour(t)

	input := map[string]any{
		"path":      "internal/ui/detail.go",
		"oldString": "one",
		"newString": "two",
	}
	lines := editLines("edit", input, 80)

	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "  - "+fg("196")+"one") {
		t.Errorf("old text is not red: %q", lines)
	}
	if !strings.Contains(joined, "  + "+fg("42")+"two") {
		t.Errorf("new text is not green: %q", lines)
	}
	if !strings.HasPrefix(lines[0], boldFg("205")) {
		t.Errorf("verb is not the tool colour: %q", lines[0])
	}
}

func TestShellCommandIsColoured(t *testing.T) {
	withColour(t)

	input := map[string]any{"command": "ls\npwd", "description": "List files"}
	lines := shellLines("Bash", input, 80)

	if !strings.HasPrefix(lines[0], boldFg("205")) {
		t.Errorf("head is not the tool colour: %q", lines[0])
	}
	if lines[1] != "  $ "+fg("222")+"ls\x1b[0m" {
		t.Errorf("command line = %q", lines[1])
	}
	if lines[2] != "    "+fg("222")+"pwd\x1b[0m" {
		t.Errorf("continuation = %q", lines[2])
	}
}

func TestWriteHeadColours(t *testing.T) {
	withColour(t)

	input := map[string]any{"path": "a.txt", "content": "one\ntwo\nthree"}
	head := writeLines("write", input, 80)[0]

	if !strings.HasPrefix(head, boldFg("205")) {
		t.Errorf("verb is not the tool colour: %q", head)
	}
	if !strings.Contains(head, "\x1b[1ma.txt") {
		t.Errorf("path is not bold: %q", head)
	}
	if !strings.Contains(head, fg("42")+"+3") {
		t.Errorf("line count is not green: %q", head)
	}
}

func TestResultStaysPlainAndErrorsRed(t *testing.T) {
	withColour(t)

	plainResult := resultLines(atuin.Content{Type: "tool_result", Content: "ok\nmore"}, 80)
	for _, line := range plainResult {
		if strings.Contains(line, "\x1b") {
			t.Errorf("result is styled: %q", line)
		}
	}

	errorResult := resultLines(atuin.Content{Type: "tool_result", Content: "boom", IsError: true}, 80)
	if !strings.HasPrefix(errorResult[0], fg("196")+"! ") {
		t.Errorf("error prefix = %q", errorResult[0])
	}
	for _, line := range errorResult {
		if !strings.HasPrefix(line, fg("196")) {
			t.Errorf("error line is not red: %q", line)
		}
	}
}
