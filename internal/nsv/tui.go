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
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	git "github.com/purpleclay/gitz"
)

var (
	borderStyle  = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false).BorderForeground(lipgloss.Color("#2b0940"))
	matchedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#139c20")).Bold(true)
	hashStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#1d1d1f")).Foreground(lipgloss.Color("#807d8a"))
	tagStyle     = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#3a1577"))
)

type Summary struct {
	Tags  []string
	Log   []git.LogEntry
	Match int
}

func PrintSummary(out io.Writer, summary Summary) {
	log := make([]string, 0, len(summary.Log))
	for i, entry := range summary.Log {
		var marker string
		if i == summary.Match {
			marker = " << matched"
		}

		log = append(log, lipgloss.JoinHorizontal(lipgloss.Left,
			hashStyle.Render(entry.AbbrevHash),
			" ",
			wordwrap.String(entry.Message, 100),
			matchedStyle.Render(marker)))
	}

	pane := lipgloss.JoinVertical(lipgloss.Top,
		"\n",
		lipgloss.JoinHorizontal(lipgloss.Left, tagStyle.Render(summary.Tags[0]), " ... ", tagStyle.Render(summary.Tags[1])),
		borderStyle.Render(strings.Join(log, "\n")),
	)

	fmt.Fprint(out, pane)
}

func PrintFormat(out io.Writer, tag Tag, format string) {
	tagf := tag.Format(format)
	header := lipgloss.JoinHorizontal(lipgloss.Left, tag.raw, " >> ", format, " >> ", tagf, "\n")

	pane := lipgloss.JoinVertical(lipgloss.Top,
		header,
		fmt.Sprintf("{{.Prefix}}%s%s", " >> ", tag.Prefix),
		fmt.Sprintf("{{.SemVer}}%s%s", " >> ", tag.SemVer),
		fmt.Sprintf("{{.Version}}%s%s", " >> ", tag.Version))

	fmt.Fprint(out, pane)
}
