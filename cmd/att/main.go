package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const version = "0.1.0"

const usage = `att - browse atuin threads

Usage:
  att [flags]

Flags:
  --version   print version and exit
  -h, --help  show this help
`

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

	fmt.Fprintln(stderr, "att: terminal UI is not wired up yet")
	return 1
}
