package ui

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type clipboardTool struct {
	name string
	args []string
}

var clipboardTools = []clipboardTool{
	{name: "wl-copy"},
	{name: "pbcopy"},
	{name: "xclip", args: []string{"-selection", "clipboard"}},
	{name: "xsel", args: []string{"--clipboard", "--input"}},
}

type detector func(name string) (string, error)

func detectClipboard(lookPath detector) (clipboardTool, bool) {
	for _, tool := range clipboardTools {
		if _, err := lookPath(tool.name); err == nil {
			return tool, true
		}
	}
	return clipboardTool{}, false
}

func execLookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func runClipboard(tool clipboardTool, text string) error {
	path, err := exec.LookPath(tool.name)
	if err != nil {
		return err
	}
	cmd := exec.Command(path, tool.args...)
	cmd.Stdin = strings.NewReader(text)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return errors.New(msg)
		}
		return err
	}
	return nil
}
