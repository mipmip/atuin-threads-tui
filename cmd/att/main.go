package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mipmip/atuin-threads-tui/internal/atuin"
	"github.com/mipmip/atuin-threads-tui/internal/ui"
)

//go:embed VERSION
var versionFile string

var version = strings.TrimSpace(versionFile)

const usage = `att - browse atuin threads

Usage:
  att [flags]

Flags:
  --version   print version and exit
  -h, --help  show this help
`

var runProgram = func(m tea.Model, out io.Writer) error {
	program := tea.NewProgram(m, tea.WithOutput(out))
	_, err := program.Run()
	return err
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("att", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stdout, usage) }

	showVersion := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *showVersion {
		fmt.Fprintf(stdout, "att %s\n", version)
		return 0
	}

	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "att: unknown argument %q\n", fs.Arg(0))
		return 2
	}

	if err := runProgram(ui.New(atuin.New()), stdout); err != nil {
		fmt.Fprintf(stderr, "att: %v\n", err)
		return 1
	}

	return 0
}
