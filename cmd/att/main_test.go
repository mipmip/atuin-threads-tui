package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "version",
			args:       []string{"--version"},
			wantCode:   0,
			wantStdout: "att " + version + "\n",
		},
		{
			name:       "help short",
			args:       []string{"-h"},
			wantCode:   0,
			wantStdout: "att - browse atuin threads",
		},
		{
			name:       "help long",
			args:       []string{"--help"},
			wantCode:   0,
			wantStdout: "Usage:",
		},
		{
			name:       "unknown flag",
			args:       []string{"--nope"},
			wantCode:   2,
			wantStderr: "flag provided but not defined",
		},
		{
			name:       "unknown argument",
			args:       []string{"extra"},
			wantCode:   2,
			wantStderr: "unknown argument",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tc.args, &stdout, &stderr)

			if code != tc.wantCode {
				t.Fatalf("exit code = %d, want %d", code, tc.wantCode)
			}
			if tc.wantStdout != "" && !strings.Contains(stdout.String(), tc.wantStdout) {
				t.Errorf("stdout = %q, want it to contain %q", stdout.String(), tc.wantStdout)
			}
			if tc.wantStderr != "" && !strings.Contains(stderr.String(), tc.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tc.wantStderr)
			}
		})
	}
}

func TestVersionComesFromFile(t *testing.T) {
	content, err := os.ReadFile("VERSION")
	if err != nil {
		t.Fatalf("reading VERSION: %v", err)
	}
	want := strings.TrimSpace(string(content))
	if want == "" {
		t.Fatal("VERSION is empty")
	}
	if version != want {
		t.Fatalf("version = %q, want %q", version, want)
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{"--version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if got := stdout.String(); got != "att "+want+"\n" {
		t.Fatalf("stdout = %q, want %q", got, "att "+want+"\n")
	}
}

func TestRunStartsUI(t *testing.T) {
	original := runProgram
	t.Cleanup(func() { runProgram = original })

	var started bool
	runProgram = func(m tea.Model, _ io.Writer) error {
		started = true
		if _, ok := m.(interface{ View() string }); !ok {
			t.Errorf("model %T has no View method", m)
		}
		return nil
	}

	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr %q)", code, stderr.String())
	}
	if !started {
		t.Fatal("program was not started")
	}
}

func TestRunReportsProgramFailure(t *testing.T) {
	original := runProgram
	t.Cleanup(func() { runProgram = original })

	runProgram = func(tea.Model, io.Writer) error {
		return errors.New("could not open a new TTY")
	}

	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "could not open a new TTY") {
		t.Errorf("stderr = %q, want the program error", stderr.String())
	}
}
