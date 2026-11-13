package main

import (
	"bytes"
	"strings"
	"testing"
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
		{
			name:       "no arguments reports missing ui",
			wantCode:   1,
			wantStderr: "terminal UI is not wired up yet",
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
