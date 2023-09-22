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
	minTagCellWidth = 20
	minLogCellWidth = 50
	logWrapAt       = 80
	markerFmt       = ">>%s<<"
)

type SummaryOptions struct {
	NoColor bool
	Out     io.Writer
	Pretty  Pretty
}

type row struct {
	tagCell string
	logCell string
}

func PrintSummary(vers []*nsv.Next, opts SummaryOptions) {
	// Build the internals of the table first. This allows the borders of the table
	// to accurately be drawn around the contents
	rows := make([]row, 0, len(vers))
	for _, ver := range vers {
		tagDiff := cell.Render(
			lipgloss.JoinVertical(lipgloss.Top,
				tagStyle.Render(ver.Tag),
				diffMark.Render(),
				tagStyle.Render(ver.PrevTag)))

		var log string
		switch Pretty(opts.Pretty) {
		case Compact:
			log = printCompactSummary(ver, opts)
		default:
			log = printFullSummary(ver, opts)
		}

		logDir := ""
		if ver.LogDir != "" {
			logDir = feintStyle.Render(fmt.Sprintf("(dir: %s)\n", ver.LogDir))
		}

		logHistory := lipgloss.JoinVertical(lipgloss.Top,
			logDir,
			log,
		)

		rows = append(rows, row{tagCell: tagDiff, logCell: logHistory})
	}

	// Determine the maximum cell widths based on the current row data
	longestTag := minTagCellWidth
	longestLog := minLogCellWidth
	for _, row := range rows {
		longestTag = max(longestTag, lipgloss.Width(row.tagCell))
		longestLog = max(longestLog, lipgloss.Width(row.logCell))
	}

	top, middle, bottom := horizontalDividers(longestTag, longestLog)

	tableRows := make([]string, 0, len(vers))
	for i, row := range rows {
		if i > 0 {
			tableRows = append(tableRows, middle)
		}

		rowHeight := lipgloss.Height(row.logCell)

		row := lipgloss.JoinHorizontal(lipgloss.Left,
			verticalDivider(rowHeight),
			cell.Copy().Width(longestTag).Render(row.tagCell),
			verticalDivider(rowHeight),
			cell.Copy().Width(longestLog).Render(row.logCell),
			verticalDivider(rowHeight),
		)
		tableRows = append(tableRows, row)
	}

	fmt.Fprint(opts.Out, lipgloss.JoinVertical(
		lipgloss.Top,
		"\n",
		top,
		strings.Join(tableRows, "\n"),
		bottom))
}

func verticalDivider(h int) string {
	return borderStyle.Render(strings.TrimRight(strings.Repeat("│\n", h), "\n"))
}

func horizontalDividers(w1, w2 int) (string, string, string) {
	top := lipgloss.JoinHorizontal(lipgloss.Left,
		borderStyle.Render(topLeft),
		borderStyle.Render(strings.Repeat(top, w1)),
		borderStyle.Render(topDivider),
		borderStyle.Render(strings.Repeat(top, w2)),
		borderStyle.Render(topRight))

	middle := lipgloss.JoinHorizontal(lipgloss.Left,
		borderStyle.Render(middleLeft),
		borderStyle.Render(strings.Repeat(middle, w1)),
		borderStyle.Render(middleDivider),
		borderStyle.Render(strings.Repeat(middle, w2)),
		borderStyle.Render(middleRight),
	)

	bottom := lipgloss.JoinHorizontal(
		lipgloss.Left,
		borderStyle.Render(bottomLeft),
		borderStyle.Render(strings.Repeat(bottom, w1)),
		borderStyle.Render(bottomDivider),
		borderStyle.Render(strings.Repeat(bottom, w2)),
		borderStyle.Render(bottomRight))

	return top, middle, bottom
}

func printFullSummary(next *nsv.Next, opts SummaryOptions) string {
	log := make([]string, 0, len(next.Log))
	for i, entry := range next.Log {
		msg := entry.Message

		marker := bullet.Render()
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

	return cell.Render(strings.Join(log, "\n\n"))
}

func printCompactSummary(next *nsv.Next, opts SummaryOptions) string {
	entry := next.Log[next.Match.Index]
	msg := entry.Message

	matched := msg[next.Match.Start:next.Match.End]
	replace := matched
	if opts.NoColor {
		replace = fmt.Sprintf(markerFmt, replace)
	}
	msg = strings.Replace(msg, matched, highlightStyle.Render(replace), 1)

	return cell.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		checkMark.Render(),
		lipgloss.JoinVertical(
			lipgloss.Top,
			hashStyle.Render(entry.AbbrevHash),
			wordwrap.String(msg, logWrapAt)),
	))
}
