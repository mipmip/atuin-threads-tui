package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
)

func scopeThreads() []atuin.Thread {
	threads := sampleThreads()
	threads[0].Host = "local-id"
	threads[0].CWD = "/home/pim/gh.mipmip/atuin-threads-tui"
	threads[1].Host = "remote-id"
	threads[1].CWD = "/home/pim/gh.mipmip/atuin-threads-tui"
	threads[2].Host = "local-id"
	threads[2].CWD = "/tmp"
	return threads
}

func scopedModel(t *testing.T, threads []atuin.Thread) Model {
	t.Helper()
	m := New(&fakeSource{threads: threads, host: "local-id"})
	m.width = 120
	m.height = 24
	m.dir = "/home/pim/gh.mipmip/atuin-threads-tui"
	return runInit(t, m)
}

func TestDirScopeKeepsOnlyThisDirectory(t *testing.T) {
	m := scopedModel(t, scopeThreads())

	m, _ = pressRune(t, m, '.')

	vis := m.visible()
	if len(vis) != 2 {
		t.Fatalf("visible = %d, want 2", len(vis))
	}
	for _, thread := range vis {
		if thread.CWD != m.dir {
			t.Fatalf("%s is outside the directory scope", thread.AtuinID)
		}
	}
}

func TestHostScopeKeepsOnlyThisHost(t *testing.T) {
	m := scopedModel(t, scopeThreads())

	m, _ = pressRune(t, m, '@')

	vis := m.visible()
	if len(vis) != 2 {
		t.Fatalf("visible = %d, want 2", len(vis))
	}
	for _, thread := range vis {
		if thread.Host != "local-id" {
			t.Fatalf("%s is not from this host", thread.AtuinID)
		}
	}
}

func TestScopesCombine(t *testing.T) {
	m := scopedModel(t, scopeThreads())

	m, _ = pressRune(t, m, '.')
	m, _ = pressRune(t, m, '@')

	vis := m.visible()
	if len(vis) != 1 || vis[0].AtuinID != "01a1d43c1f237541a0e1c6c42cf53ec4" {
		t.Fatalf("visible = %v, want only the first thread", vis)
	}
}

func TestScopeCoversEveryThread(t *testing.T) {
	m := scopedModel(t, scopeThreads())
	m.cursor = 2
	m.offset = 2

	m, _ = pressRune(t, m, '@')
	if m.cursor != 0 || m.offset != 0 {
		t.Fatalf("cursor = %d offset = %d, want 0 and 0", m.cursor, m.offset)
	}

	m, _ = pressRune(t, m, '@')
	if len(m.visible()) != 3 {
		t.Fatalf("visible = %d, want 3", len(m.visible()))
	}
}

func TestUnknownHostIsNotAScope(t *testing.T) {
	m := New(&fakeSource{threads: scopeThreads(), hostErr: errors.New("no host id")})
	m.width = 120
	m.height = 24
	m = runInit(t, m)

	m, _ = pressRune(t, m, '@')

	if m.scopeHost {
		t.Fatal("host scope is on although the host is unknown")
	}
	if len(m.visible()) != 3 {
		t.Fatalf("visible = %d, want 3", len(m.visible()))
	}
	if !strings.Contains(plain(m.View()), "host unknown") {
		t.Fatalf("status does not explain the missing host: %q", plain(m.View()))
	}
}

func TestScopeStatusNamesTheScope(t *testing.T) {
	m := scopedModel(t, scopeThreads())

	m, _ = pressRune(t, m, '.')
	view := plain(m.View())
	if !strings.Contains(view, "dir: atuin-threads-tui · 2 of 3 threads") {
		t.Fatalf("dir scope missing from status: %q", view)
	}

	m, _ = pressRune(t, m, '@')
	view = plain(m.View())
	if !strings.Contains(view, "dir: atuin-threads-tui · host: local · 1 of 3 threads") {
		t.Fatalf("both scopes missing from status: %q", view)
	}
}

func TestListHintsMentionTheScopes(t *testing.T) {
	m := scopedModel(t, scopeThreads())

	view := plain(m.View())
	for _, want := range []string{". dir", "@ host", "/ filter"} {
		if !strings.Contains(view, want) {
			t.Fatalf("hints do not mention %q: %q", want, view)
		}
	}
}

func TestScopeCombinesWithTheTextFilter(t *testing.T) {
	m := scopedModel(t, scopeThreads())

	m, _ = pressRune(t, m, '.')
	m, _ = pressRune(t, m, '/')
	m = typeText(t, m, "second")

	vis := m.visible()
	if len(vis) != 1 || vis[0].Title != "Second thread" {
		t.Fatalf("visible = %v, want the second thread only", vis)
	}
}
