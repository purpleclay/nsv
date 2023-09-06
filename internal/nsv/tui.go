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

package nsv

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	borderStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false).BorderForeground(lipgloss.Color("#2b0940"))
	hashStyle      = lipgloss.NewStyle().Background(lipgloss.Color("#1d1d1f")).Foreground(lipgloss.Color("#807d8a"))
	highlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("#bf31f7"))
	tagStyle       = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#3a1577"))
	feintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#807d8a"))
	chevron        = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#b769d6")).Render(">>")
	checkMark      = lipgloss.NewStyle().Foreground(lipgloss.Color("#139c20")).SetString("✓ ")
)

func PrintSummary(next Next, opts Options) {
	if opts.NoColor {
		tagStyle = tagStyle.UnsetPadding()
	}

	log := make([]string, 0, len(next.Log))
	for i, entry := range next.Log {
		msg := entry.Message

		var marker string
		if i == next.Match.Index {
			marker = checkMark.Render()

			matched := msg[next.Match.Start:next.Match.End]
			replace := matched
			if opts.NoColor {
				replace = fmt.Sprintf(">>%s<<", replace)
			}
			msg = strings.Replace(msg, matched, highlightStyle.Render(replace), 1)
		}

		log = append(log, lipgloss.JoinHorizontal(
			lipgloss.Left,
			marker,
			lipgloss.JoinVertical(
				lipgloss.Top,
				hashStyle.Render(entry.AbbrevHash),
				wordwrap.String(msg, 80)),
		))
	}

	logDir := ""
	if next.LogDir != "" {
		logDir = feintStyle.Render(fmt.Sprintf(" (%s)", next.LogDir))
	}

	pane := lipgloss.JoinVertical(lipgloss.Top,
		"\n",
		lipgloss.JoinHorizontal(lipgloss.Left, tagStyle.Render("HEAD"), logDir),
		borderStyle.Render(strings.Join(log, "\n\n")),
		tagStyle.Render(next.PrevTag),
	)

	fmt.Fprint(opts.Err, pane)
}

func PrintFormat(tag Tag, opts Options) {
	tagf := tag.Format(opts.VersionFormat)
	header := lipgloss.JoinHorizontal(lipgloss.Left, tag.raw, chevron, opts.VersionFormat, chevron, tagf, "\n")

	pane := lipgloss.JoinVertical(lipgloss.Top,
		header,
		fmt.Sprintf("{{.Prefix}} %s%s", chevron, tag.Prefix),
		fmt.Sprintf("{{.SemVer}} %s%s", chevron, tag.SemVer),
		fmt.Sprintf("{{.Version}}%s%s", chevron, tag.Version))

	fmt.Fprint(opts.Err, pane)
}
