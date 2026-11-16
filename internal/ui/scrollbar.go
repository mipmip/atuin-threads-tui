package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	scrollTrackCell = "│"
	scrollThumbCell = "█"
)

var (
	scrollTrackStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	scrollThumbStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type scrollbar struct {
	top   int
	thumb int
	shown bool
}

// newScrollbar places a thumb of thumb cells inside a track of track cells.
func newScrollbar(scroll, total, track int) scrollbar {
	if track <= 0 || total <= track {
		return scrollbar{}
	}

	thumb := (track*track + total/2) / total
	if thumb < 1 {
		thumb = 1
	}
	if thumb > track {
		thumb = track
	}

	span := total - track
	top := 0
	if span > 0 && track > thumb {
		top = (scroll*(track-thumb) + span/2) / span
	}
	if top < 0 {
		top = 0
	}
	if top > track-thumb {
		top = track - thumb
	}

	return scrollbar{top: top, thumb: thumb, shown: true}
}

func (m Model) scrollbar() scrollbar {
	if m.mode != modeDetail || m.rendering || !m.detail.built {
		return scrollbar{}
	}
	return newScrollbar(m.detailScroll, m.detailTotal(), m.bodyHeight())
}

// withScrollbar fills the last column of the body with the bar.
func (m Model) withScrollbar(lines []string) []string {
	if m.width < 2 {
		return lines
	}

	bar := m.scrollbar()
	for i, line := range lines {
		line = ansi.Truncate(line, m.width-1, "")
		if gap := m.width - 1 - ansi.StringWidth(line); gap > 0 {
			line += strings.Repeat(" ", gap)
		}

		switch {
		case !bar.shown:
			line += " "
		case i >= bar.top && i < bar.top+bar.thumb:
			line += scrollThumbStyle.Render(scrollThumbCell)
		default:
			line += scrollTrackStyle.Render(scrollTrackCell)
		}
		lines[i] = line
	}
	return lines
}
