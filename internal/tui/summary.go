/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/purpleclay/lipgloss-theme"
	"github.com/purpleclay/nsv/internal/nsv"
)

type Pretty string

const (
	Full    Pretty = "full"
	Compact Pretty = "compact"
)

var PrettyFormats = []string{string(Full), string(Compact)}

const (
	minTagCellWidth = 15
	minLogCellWidth = 50
	logWrapAt       = 80
	markerFmt       = ">>%s<<"
)

var (
	highlight = lipgloss.NewStyle().Background(
		lipgloss.AdaptiveColor{
			Light: string(theme.S400),
			Dark:  string(theme.S200),
		}).
		Foreground(lipgloss.AdaptiveColor{Light: "#ffffff"})

	diffMark = lipgloss.NewStyle().Foreground(
		lipgloss.AdaptiveColor{
			Light: string(theme.Green900),
			Dark:  string(theme.Green700),
		}).
		SetString(" ↑↑")

	faint  = lipgloss.NewStyle().Faint(true)
	bullet = faint.Copy().SetString("> ")
)

type SummaryOptions struct {
	NoColor bool
	Out     io.Writer
	Pretty  Pretty
}

func PrintSummary(vers []*nsv.Next, opts SummaryOptions) {
	var rows [][]string
	for _, ver := range vers {
		tagDiff := lipgloss.JoinVertical(lipgloss.Top,
			theme.H1.Render(ver.Tag),
			diffMark.Render(),
			theme.H4.Render(ver.PrevTag),
		)

		var log string
		switch Pretty(opts.Pretty) {
		case Compact:
			log = printCompactSummary(ver, opts)
		default:
			log = printFullSummary(ver, opts)
		}

		var sum []string
		if ver.LogDir != "" {
			sum = append(sum, faint.Render(fmt.Sprintf("(dir: %s)\n", ver.LogDir)))
		}
		sum = append(sum, log)

		logHistory := lipgloss.JoinVertical(lipgloss.Top, strings.Join(sum, "\n"))
		rows = append(rows, []string{tagDiff, logHistory})
	}

	out := theme.NewTable(rows).
		Border(theme.ThinBorder).
		Widths(minTagCellWidth, minLogCellWidth).
		String()

	fmt.Fprint(opts.Out, lipgloss.JoinVertical(
		lipgloss.Top,
		"",
		out,
	))
}

func printFullSummary(next *nsv.Next, opts SummaryOptions) string {
	log := make([]string, 0, len(next.Log))
	for i, entry := range next.Log {
		msg := entry.Message

		marker := bullet.Render()
		if i == next.Match.Index {
			marker = theme.Tick

			matched := msg[next.Match.Start:next.Match.End]
			replace := matched
			if opts.NoColor {
				replace = fmt.Sprintf(markerFmt, replace)
			}
			msg = strings.Replace(msg, matched, highlight.Render(replace), 1)
		}

		log = append(log, lipgloss.JoinHorizontal(
			lipgloss.Left,
			marker,
			lipgloss.JoinVertical(
				lipgloss.Top,
				theme.Mark.Render(entry.AbbrevHash),
				wordwrap.String(msg, logWrapAt)),
		))
	}

	return strings.Join(log, "\n\n")
}

func printCompactSummary(next *nsv.Next, opts SummaryOptions) string {
	entry := next.Log[next.Match.Index]
	msg := entry.Message

	matched := msg[next.Match.Start:next.Match.End]
	replace := matched
	if opts.NoColor {
		replace = fmt.Sprintf(markerFmt, replace)
	}
	msg = strings.Replace(msg, matched, highlight.Render(replace), 1)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		theme.Tick,
		lipgloss.JoinVertical(
			lipgloss.Top,
			theme.Mark.Render(entry.AbbrevHash),
			wordwrap.String(msg, logWrapAt)),
	)
}
