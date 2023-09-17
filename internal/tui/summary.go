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
	"github.com/purpleclay/nsv/internal/nsv"
)

type Pretty string

const (
	Full    Pretty = "full"
	Compact Pretty = "compact"
)

var PrettyFormats = []string{string(Full), string(Compact)}

const (
	logWrapAt = 80
	markerFmt = ">>%s<<"
)

type SummaryOptions struct {
	NoColor bool
	Out     io.Writer
	Pretty  Pretty
}

func PrintSummary(next nsv.Next, opts SummaryOptions) {
	if opts.NoColor {
		tagStyle = tagStyle.UnsetPadding()
	}

	var log string
	switch Pretty(opts.Pretty) {
	case Compact:
		log = printCompactSummary(next, opts)
	default:
		log = printFullSummary(next, opts)
	}

	logDir := ""
	if next.LogDir != "" {
		logDir = feintStyle.Render(fmt.Sprintf(" (dir: %s)", next.LogDir))
	}

	pane := lipgloss.JoinVertical(lipgloss.Top,
		"\n",
		lipgloss.JoinHorizontal(lipgloss.Left, tagStyle.Render("HEAD"), logDir),
		borderStyle.Render(log),
		tagStyle.Render(next.PrevTag),
	)

	fmt.Fprint(opts.Out, pane)
}

func printFullSummary(next nsv.Next, opts SummaryOptions) string {
	log := make([]string, 0, len(next.Log))
	for i, entry := range next.Log {
		msg := entry.Message

		var marker string
		if i == next.Match.Index {
			marker = checkMark.Render()

			matched := msg[next.Match.Start:next.Match.End]
			replace := matched
			if opts.NoColor {
				replace = fmt.Sprintf(markerFmt, replace)
			}
			msg = strings.Replace(msg, matched, highlightStyle.Render(replace), 1)
		}

		log = append(log, lipgloss.JoinHorizontal(
			lipgloss.Left,
			marker,
			lipgloss.JoinVertical(
				lipgloss.Top,
				hashStyle.Render(entry.AbbrevHash),
				wordwrap.String(msg, logWrapAt)),
		))
	}

	return strings.Join(log, "\n\n")
}

func printCompactSummary(next nsv.Next, opts SummaryOptions) string {
	entry := next.Log[next.Match.Index]
	msg := entry.Message

	matched := msg[next.Match.Start:next.Match.End]
	replace := matched
	if opts.NoColor {
		replace = fmt.Sprintf(markerFmt, replace)
	}
	msg = strings.Replace(msg, matched, highlightStyle.Render(replace), 1)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		hashStyle.Render(entry.AbbrevHash),
		wordwrap.String(msg, logWrapAt))
}
