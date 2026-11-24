package atuin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Runner executes one command and returns its stdout.
type Runner func(ctx context.Context, args ...string) ([]byte, error)

// Client reads threads through the atuin CLI.
type Client struct {
	binary string
	run    Runner
}

// Option adjusts a Client.
type Option func(*Client)

// WithBinary sets the atuin executable to run.
func WithBinary(path string) Option {
	return func(c *Client) { c.binary = path }
}

// WithRunner replaces the process runner, for tests.
func WithRunner(r Runner) Option {
	return func(c *Client) { c.run = r }
}

// New returns a client that runs the atuin binary on PATH.
func New(opts ...Option) *Client {
	c := &Client{binary: "atuin", run: execRunner}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func execRunner(ctx context.Context, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command given")
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return nil, fmt.Errorf("%w: %s", err, msg)
		}
		return nil, err
	}
	return out, nil
}

// List returns captured threads, newest first.
func (c *Client) List(ctx context.Context) ([]Thread, error) {
	out, err := c.run(ctx, c.binary, "threads", "list", "--style", "json")
	if err != nil {
		return nil, fmt.Errorf("listing threads: %w", err)
	}

	var threads []Thread
	if err := json.Unmarshal(out, &threads); err != nil {
		return nil, fmt.Errorf("decoding thread list: %w", err)
	}
	if threads == nil {
		threads = []Thread{}
	}
	return threads, nil
}

// Show returns one thread with its messages.
func (c *Client) Show(ctx context.Context, id string) (Conversation, error) {
	out, err := c.run(ctx, c.binary, "threads", "show", "--style", "json", id)
	if err != nil {
		return Conversation{}, fmt.Errorf("reading thread %s: %w", id, err)
	}

	var conversation Conversation
	if err := json.Unmarshal(out, &conversation); err != nil {
		return Conversation{}, fmt.Errorf("decoding thread %s: %w", id, err)
	}
	return conversation, nil
}

// Search returns threads matching every word in query.
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Thread, error) {
	words := strings.Fields(query)
	if len(words) == 0 {
		return []Thread{}, nil
	}

	args := []string{
		c.binary, "threads", "search",
		"--style", "json",
		"--limit", strconv.Itoa(limit),
		"--",
	}
	args = append(args, words...)

	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("searching threads: %w", err)
	}

	var threads []Thread
	if err := json.Unmarshal(out, &threads); err != nil {
		return nil, fmt.Errorf("decoding search result: %w", err)
	}
	if threads == nil {
		threads = []Thread{}
	}
	return threads, nil
}

// Host returns the id atuin uses for this machine.
func (c *Client) Host(ctx context.Context) (string, error) {
	out, err := c.run(ctx, c.binary, "info")
	if err != nil {
		return "", fmt.Errorf("reading host id: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		rest, ok := strings.CutPrefix(strings.TrimSpace(line), "host id:")
		if ok && strings.TrimSpace(rest) != "" {
			return strings.TrimSpace(rest), nil
		}
	}
	return "", fmt.Errorf("no host id in atuin info output")
}
