package atuin

import (
	"context"
	"errors"
	"strings"
	"testing"
)

const listFixture = `[
  {
    "harness": "claude-code",
    "session_id": "s-2",
    "atuin_id": "01a1acb62fa877ab83c0c981da8da47e",
    "host": "host-1",
    "cwd": "/home/pim/work/two",
    "git_branch": "main",
    "model": "claude-opus-5-5",
    "started_at": "2026-11-05T17:16:23.334+00:00",
    "updated_at": "2026-11-13T11:00:00.000+00:00",
    "message_count": 12,
    "tokens": {"input": 10, "output": 20, "cache_read": 30, "cache_write": 0, "reasoning": 5},
    "title": "Second thread",
    "preview": "Preview two",
    "last_reply": "Last two",
    "future_field": true
  },
  {
    "harness": "opencode",
    "session_id": "s-1",
    "atuin_id": "01a1d43c1f237541a0e1c6c42cf53ec4",
    "cwd": "/home/pim/work/one",
    "started_at": "2026-11-12T09:00:00Z",
    "updated_at": "2026-11-12T09:30:00Z",
    "message_count": 3,
    "title": "First thread"
  }
]`

const showFixture = `{
  "session": {"atuin_id": "01a1d43c1f237541a0e1c6c42cf53ec4", "title": "First thread"},
  "messages": [
    {
      "id": "m1",
      "role": "user",
      "timestamp": "2026-11-12T09:00:00Z",
      "cwd": "/home/pim/work/one",
      "content": [{"type": "text", "text": "hello"}],
      "stop_reason": "end_turn"
    },
    {
      "id": "m2",
      "role": "assistant",
      "timestamp": "2026-11-12T09:00:01Z",
      "content": [
        {"type": "thinking", "text": "thinking"},
        {"type": "tool_call", "id": "c1", "name": "shell", "input": "{\"command\":\"ls\"}"},
        {"type": "tool_result", "tool_use_id": "c1", "content": "file", "is_error": false}
      ]
    }
  ]
}`

type call struct {
	args []string
}

func fakeRunner(out string, err error, calls *[]call) Runner {
	return func(_ context.Context, args ...string) ([]byte, error) {
		if calls != nil {
			*calls = append(*calls, call{args: args})
		}
		return []byte(out), err
	}
}

func TestList(t *testing.T) {
	var calls []call
	c := New(WithRunner(fakeRunner(listFixture, nil, &calls)))

	threads, err := c.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(threads) != 2 {
		t.Fatalf("got %d threads, want 2", len(threads))
	}
	if threads[0].Title != "Second thread" || threads[1].Title != "First thread" {
		t.Errorf("order not preserved: %q, %q", threads[0].Title, threads[1].Title)
	}
	if threads[0].Tokens.CacheRead != 30 || threads[0].Tokens.Reasoning != 5 {
		t.Errorf("tokens not decoded: %+v", threads[0].Tokens)
	}
	if threads[0].MessageCount != 12 {
		t.Errorf("message count = %d, want 12", threads[0].MessageCount)
	}
	if got := threads[1].UpdatedAt.Year(); got != 2026 {
		t.Errorf("timestamp not decoded: %v", threads[1].UpdatedAt)
	}

	want := []string{"atuin", "threads", "list", "--style", "json"}
	if strings.Join(calls[0].args, " ") != strings.Join(want, " ") {
		t.Errorf("args = %v, want %v", calls[0].args, want)
	}
}

func TestListEmpty(t *testing.T) {
	c := New(WithRunner(fakeRunner("[]", nil, nil)))

	threads, err := c.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if threads == nil || len(threads) != 0 {
		t.Fatalf("threads = %#v, want empty non-nil slice", threads)
	}
}

func TestListCommandFailure(t *testing.T) {
	c := New(WithRunner(fakeRunner("", errors.New("exit status 2: no such command"), nil)))

	_, err := c.List(context.Background())
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "listing threads") || !strings.Contains(err.Error(), "no such command") {
		t.Errorf("error = %q, want operation and stderr text", err)
	}
}

func TestListMalformed(t *testing.T) {
	c := New(WithRunner(fakeRunner("[{", nil, nil)))

	_, err := c.List(context.Background())
	if err == nil || !strings.Contains(err.Error(), "decoding") {
		t.Fatalf("error = %v, want a decoding error", err)
	}
}

func TestShow(t *testing.T) {
	var calls []call
	c := New(WithRunner(fakeRunner(showFixture, nil, &calls)))

	conversation, err := c.Show(context.Background(), "01a1d43c1f237541a0e1c6c42cf53ec4")
	if err != nil {
		t.Fatalf("Show: %v", err)
	}

	if conversation.Session.Title != "First thread" {
		t.Errorf("session title = %q", conversation.Session.Title)
	}
	if len(conversation.Messages) != 2 {
		t.Fatalf("got %d messages, want 2", len(conversation.Messages))
	}
	if conversation.Messages[1].Content[1].Name != "shell" {
		t.Errorf("tool call not decoded: %+v", conversation.Messages[1].Content[1])
	}
	if conversation.Messages[1].Content[2].Content != "file" {
		t.Errorf("tool result not decoded: %+v", conversation.Messages[1].Content[2])
	}

	want := []string{"atuin", "threads", "show", "--style", "json", "01a1d43c1f237541a0e1c6c42cf53ec4"}
	if strings.Join(calls[0].args, " ") != strings.Join(want, " ") {
		t.Errorf("args = %v, want %v", calls[0].args, want)
	}
}

func TestShowFailure(t *testing.T) {
	c := New(WithRunner(fakeRunner("", errors.New("thread not found"), nil)))

	_, err := c.Show(context.Background(), "abc")
	if err == nil || !strings.Contains(err.Error(), "reading thread abc") {
		t.Fatalf("error = %v, want a read error naming the thread", err)
	}
}

func TestSearch(t *testing.T) {
	var calls []call
	c := New(WithRunner(fakeRunner(listFixture, nil, &calls)))

	threads, err := c.Search(context.Background(), "openspec beans", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(threads) != 2 {
		t.Fatalf("got %d threads, want 2", len(threads))
	}

	want := []string{
		"atuin", "threads", "search", "--style", "json",
		"--limit", "5", "--", "openspec", "beans",
	}
	if strings.Join(calls[0].args, " ") != strings.Join(want, " ") {
		t.Errorf("args = %v, want %v", calls[0].args, want)
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	c := New(WithRunner(fakeRunner("", errors.New("must not be called"), nil)))

	threads, err := c.Search(context.Background(), "   ", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(threads) != 0 {
		t.Errorf("threads = %v, want none", threads)
	}
}

func TestSearchMalformed(t *testing.T) {
	c := New(WithRunner(fakeRunner("nope", nil, nil)))

	_, err := c.Search(context.Background(), "x", 5)
	if err == nil || !strings.Contains(err.Error(), "decoding search result") {
		t.Fatalf("error = %v, want a decode error", err)
	}
}

func TestMissingBinary(t *testing.T) {
	c := New(WithBinary("att-atuin-does-not-exist"))

	_, err := c.List(context.Background())
	if err == nil || !strings.Contains(err.Error(), "listing threads") {
		t.Fatalf("error = %v, want a listing error", err)
	}
}

const usageFixture = `{
  "session": {"atuin_id": "01a1d43c1f237541a0e1c6c42cf53ec4", "title": "Usage"},
  "messages": [
    {
      "id": "m1",
      "role": "assistant",
      "timestamp": "2026-11-12T09:00:01Z",
      "cwd": "/home/pim/work/one",
      "model": "deepseek-v4.1-flash",
      "tokens": {"input": 2040, "output": 517, "cache_read": 49792, "cache_write": 0, "reasoning": 62},
      "content": [
        {"type": "patch", "tool_use_id": "p1", "files": [
          {"path": "flake.nix", "change": "update", "additions": 2, "deletions": 2, "diff": "@@ -1 +1 @@"}
        ]}
      ]
    },
    {"id": "m2", "role": "tool", "timestamp": "2026-11-12T09:00:02Z", "content": []}
  ]
}`

func TestShowDecodesUsageAndPatchFiles(t *testing.T) {
	c := New(WithRunner(fakeRunner(usageFixture, nil, nil)))

	conversation, err := c.Show(context.Background(), "01a1d43c1f237541a0e1c6c42cf53ec4")
	if err != nil {
		t.Fatalf("Show: %v", err)
	}

	first := conversation.Messages[0]
	if first.Model != "deepseek-v4.1-flash" {
		t.Errorf("model = %q", first.Model)
	}
	if first.Tokens == nil {
		t.Fatal("tokens not decoded")
	}
	if first.Tokens.Output != 517 || first.Tokens.Reasoning != 62 || first.Tokens.CacheRead != 49792 {
		t.Errorf("tokens = %+v", *first.Tokens)
	}

	files := first.Content[0].Files
	if len(files) != 1 {
		t.Fatalf("files = %#v, want one", files)
	}
	want := PatchFile{Path: "flake.nix", Change: "update", Additions: 2, Deletions: 2, Diff: "@@ -1 +1 @@"}
	if files[0] != want {
		t.Errorf("file = %+v, want %+v", files[0], want)
	}

	if conversation.Messages[1].Tokens != nil {
		t.Errorf("message without usage decoded tokens: %+v", conversation.Messages[1].Tokens)
	}
}

const infoFixture = `name: atuin
version: 18.24.6
commit:  NO_GIT

Host info:
host id: 01a18745351c74618f8f9d94eac717e8
`

func TestHostReadsInfo(t *testing.T) {
	var calls []call
	c := New(WithRunner(fakeRunner(infoFixture, nil, &calls)))

	host, err := c.Host(context.Background())
	if err != nil {
		t.Fatalf("Host: %v", err)
	}
	if host != "01a18745351c74618f8f9d94eac717e8" {
		t.Fatalf("host = %q", host)
	}
	want := []string{"atuin", "info"}
	if strings.Join(calls[0].args, " ") != strings.Join(want, " ") {
		t.Errorf("args = %v, want %v", calls[0].args, want)
	}
}

func TestHostWithoutId(t *testing.T) {
	c := New(WithRunner(fakeRunner("Host info:\nno id here\n", nil, nil)))

	if _, err := c.Host(context.Background()); err == nil {
		t.Fatal("want an error when atuin info carries no host id")
	}
}

func TestHostCommandFailure(t *testing.T) {
	c := New(WithRunner(fakeRunner("", errors.New("boom"), nil)))

	if _, err := c.Host(context.Background()); err == nil || !strings.Contains(err.Error(), "reading host id") {
		t.Fatalf("err = %v", err)
	}
}
