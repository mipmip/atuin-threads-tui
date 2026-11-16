package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

type exportMessage struct {
	index   int
	heading string
	body    []string
}

func exportMarkdown(messages []atuin.Message, hidden typeMask) string {
	parts := make([]exportMessage, 0, len(messages))
	for i, msg := range messages {
		blocks := visibleBlocks(msg, hidden)
		if len(blocks) == 0 {
			continue
		}
		parts = append(parts, exportMessage{
			index:   i,
			heading: fmt.Sprintf("### %s %s", clock(msg.Timestamp), authorLabel(msg.Role)),
			body:    messageMarkdown(blocks),
		})
	}

	keepHeadings := len(parts) > 1 || (len(parts) == 1 && len(parts[0].body) > 1)

	var b strings.Builder
	for i, part := range parts {
		if i > 0 {
			b.WriteString("\n")
		}
		if keepHeadings {
			b.WriteString(part.heading)
			b.WriteString("\n\n")
		}
		b.WriteString(strings.Join(part.body, "\n"))
		b.WriteString("\n")
	}
	return b.String()
}

func messageMarkdown(blocks []atuin.Content) []string {
	var out []string
	for _, block := range blocks {
		if len(out) > 0 {
			out = append(out, "")
		}
		out = append(out, blockMarkdown(block)...)
	}
	return out
}

func blockMarkdown(block atuin.Content) []string {
	switch block.Type {
	case "text", "summary":
		return strings.Split(strings.TrimRight(block.Text, "\n"), "\n")
	case "thinking":
		return quote(strings.Split(strings.TrimRight(block.Text, "\n"), "\n"))
	case "tool_call":
		return toolCallMarkdown(block)
	case "tool_result":
		lines := resultMarkdown(block)
		if block.IsError && len(lines) > 0 {
			lines[0] = "! " + lines[0]
		}
		return fenced(lines)
	case "error":
		lines := strings.Split(strings.TrimRight(block.Text, "\n"), "\n")
		if len(lines) > 0 {
			lines[0] = "! " + lines[0]
		}
		return fenced(lines)
	case "patch":
		return patchMarkdown(block)
	default:
		return []string{"(" + block.Type + ")"}
	}
}

func resultMarkdown(block atuin.Content) []string {
	text := block.Content
	if strings.TrimSpace(text) == "" {
		text = block.Text
	}
	return strings.Split(strings.TrimRight(text, "\n"), "\n")
}

func toolCallMarkdown(block atuin.Content) []string {
	input := parseInput(block.Input)
	name := strings.TrimSpace(block.Name)
	kind := strings.ToLower(name)

	if input == nil {
		head := "**" + name + "**"
		if strings.TrimSpace(block.Input) == "" {
			return []string{head}
		}
		return append([]string{head}, fenced(strings.Split(block.Input, "\n"))...)
	}

	switch kind {
	case "bash", "shell":
		head := "**" + name + "**"
		if desc := field(input, "description"); desc != "" {
			head += " - " + desc
		}
		if command := field(input, "command"); command != "" {
			return append([]string{head}, fenced([]string{"$ " + command})...)
		}
		return []string{head}
	case "write":
		path := firstField(input, "file_path", "path")
		head := "**" + name + "** " + path
		if content := firstField(input, "content"); content != "" {
			lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
			head += fmt.Sprintf(" (+%d lines)", len(lines))
			return append([]string{head}, fenced(lines)...)
		}
		return []string{head}
	case "edit":
		path := firstField(input, "file_path", "path")
		old := firstField(input, "old_string", "oldString")
		add := firstField(input, "new_string", "newString")
		head := "**" + name + "** " + path
		var out []string
		if old != "" || add != "" {
			out = []string{head, ""}
			out = append(out, fenced(diffLines(old, add))...)
			return out
		}
		return []string{head}
	case "read":
		path := firstField(input, "file_path", "path")
		return []string{"**" + name + "** " + path}
	case "skill":
		skill := firstField(input, "skill", "name")
		return []string{"**" + name + "** " + skill}
	case "compress":
		if topic := field(input, "topic"); topic != "" {
			return []string{"**" + name + "** " + topic}
		}
		return []string{"**" + name + "**"}
	default:
		parts := []string{"**" + name + "**"}
		keys := sortedKeys(input)
		for _, key := range keys {
			value := valueString(input[key])
			if value == "" {
				continue
			}
			parts = append(parts, key+"="+value)
		}
		return []string{strings.Join(parts, " ")}
	}
}

func diffLines(old, add string) []string {
	var out []string
	for _, line := range strings.Split(strings.TrimRight(old, "\n"), "\n") {
		if line == "" && len(out) == 0 {
			continue
		}
		out = append(out, "-"+line)
	}
	for _, line := range strings.Split(strings.TrimRight(add, "\n"), "\n") {
		if line == "" && len(out) == 0 {
			continue
		}
		out = append(out, "+"+line)
	}
	return out
}

func patchMarkdown(block atuin.Content) []string {
	var out []string
	if len(block.Files) == 0 {
		return fenced(strings.Split(strings.TrimRight(block.Text, "\n"), "\n"))
	}
	for _, file := range block.Files {
		head := fmt.Sprintf("**patch** %s (+%d -%d)", file.Path, file.Additions, file.Deletions)
		out = append(out, head)
		out = append(out, "")
		out = append(out, fenced(strings.Split(strings.TrimRight(file.Diff, "\n"), "\n"))...)
	}
	return out
}

func quote(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			out = append(out, ">")
			continue
		}
		out = append(out, "> "+line)
	}
	return out
}

func fenced(lines []string) []string {
	marks := "```"
	for _, line := range lines {
		if run := backtickRun(line); run >= len(marks) {
			marks = strings.Repeat("`", run+1)
		}
	}
	out := make([]string, 0, len(lines)+2)
	out = append(out, marks)
	out = append(out, lines...)
	out = append(out, marks)
	return out
}

func backtickRun(line string) int {
	longest, current := 0, 0
	for _, r := range line {
		if r == '`' {
			current++
			if current > longest {
				longest = current
			}
			continue
		}
		current = 0
	}
	return longest
}

func valueString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%g", typed)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	}
	return ""
}

func exportName(shortID string, first, last int) string {
	return fmt.Sprintf("att-%s-%d-%d.md", shortID, first+1, last+1)
}

func writeExport(dir, name, payload string) (string, error) {
	path := filepath.Join(dir, name)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	if _, err := file.WriteString(payload); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return name, nil
}
