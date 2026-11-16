package ui

import (
	"strings"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

type typeMask uint8

const (
	typePrompts typeMask = 1 << iota
	typeReplies
	typeThinking
	typeCalls
	typeResults
	typePatches
	typeSummaries
	typeNotes
)

const (
	allTypes          = typePrompts | typeReplies | typeThinking | typeCalls | typeResults | typePatches | typeSummaries | typeNotes
	conversationTypes = typePrompts | typeReplies
	toolTypes         = typeCalls | typeResults | typePatches
	promptTypes       = typePrompts
)

// shows reports whether a block of this category is visible. A category of
// zero is a block type the program does not know, which is always visible.
func (h typeMask) shows(category typeMask) bool {
	return category == 0 || h&category == 0
}

func (h typeMask) toggle(category typeMask) typeMask {
	return h ^ category
}

type typeRow struct {
	category typeMask
	label    string
	note     string
}

var typeRows = []typeRow{
	{typePrompts, "prompts", "your messages"},
	{typeReplies, "replies", "assistant prose"},
	{typeThinking, "thinking", "hidden reasoning"},
	{typeCalls, "tool calls", "Bash, edit, read"},
	{typeResults, "tool results", "command output"},
	{typePatches, "patches", "diffs"},
	{typeSummaries, "summaries", "session wrap-ups"},
	{typeNotes, "notes", "attachments, titles, empty system"},
}

func categoryOf(block atuin.Content, role string) typeMask {
	switch block.Type {
	case "text":
		if strings.EqualFold(strings.TrimSpace(role), "user") {
			return typePrompts
		}
		return typeReplies
	case "thinking":
		return typeThinking
	case "tool_call":
		return typeCalls
	case "tool_result", "error":
		return typeResults
	case "patch":
		return typePatches
	case "summary":
		return typeSummaries
	}
	return 0
}

// countTypes counts the blocks per category and the messages without content.
func countTypes(conversation atuin.Conversation) map[typeMask]int {
	counts := make(map[typeMask]int, len(typeRows))
	for _, msg := range conversation.Messages {
		if len(msg.Content) == 0 {
			counts[typeNotes]++
			continue
		}
		for _, block := range msg.Content {
			if category := categoryOf(block, msg.Role); category != 0 {
				counts[category]++
			}
		}
	}
	return counts
}

func visibleBlocks(msg atuin.Message, hidden typeMask) []atuin.Content {
	blocks := make([]atuin.Content, 0, len(msg.Content))
	for _, block := range msg.Content {
		if hidden.shows(categoryOf(block, msg.Role)) {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func hasVisibleBlock(msg atuin.Message, hidden typeMask) bool {
	for _, block := range msg.Content {
		if hidden.shows(categoryOf(block, msg.Role)) {
			return true
		}
	}
	return false
}
