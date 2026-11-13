package atuin

import (
	"context"
	"time"
)

// Tokens is the usage atuin attributes to a thread.
type Tokens struct {
	Input      int `json:"input"`
	Output     int `json:"output"`
	CacheRead  int `json:"cache_read"`
	CacheWrite int `json:"cache_write"`
	Reasoning  int `json:"reasoning"`
}

// Thread is one captured agent session as atuin lists it.
type Thread struct {
	Harness      string    `json:"harness"`
	SessionID    string    `json:"session_id"`
	AtuinID      string    `json:"atuin_id"`
	Host         string    `json:"host"`
	CWD          string    `json:"cwd"`
	GitBranch    string    `json:"git_branch"`
	Model        string    `json:"model"`
	StartedAt    time.Time `json:"started_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MessageCount int       `json:"message_count"`
	Tokens       Tokens    `json:"tokens"`
	Title        string    `json:"title"`
	Preview      string    `json:"preview"`
	LastReply    string    `json:"last_reply"`
}

// Content is one block inside a message.
type Content struct {
	Type      string `json:"type"`
	Text      string `json:"text"`
	Name      string `json:"name"`
	Input     string `json:"input"`
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error"`
}

// Message is one recorded turn in a thread.
type Message struct {
	ID         string    `json:"id"`
	Role       string    `json:"role"`
	Timestamp  time.Time `json:"timestamp"`
	CWD        string    `json:"cwd"`
	Content    []Content `json:"content"`
	StopReason string    `json:"stop_reason"`
}

// Conversation is a thread with its messages.
type Conversation struct {
	Session  Thread    `json:"session"`
	Messages []Message `json:"messages"`
}

// Source provides threads to the UI.
type Source interface {
	List(ctx context.Context) ([]Thread, error)
	Show(ctx context.Context, id string) (Conversation, error)
	Search(ctx context.Context, query string, limit int) ([]Thread, error)
}
