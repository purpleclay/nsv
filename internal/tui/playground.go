/*
Copyright (c) 2023 - 2024 Purple Clay

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

	"github.com/charmbracelet/lipgloss"
	theme "github.com/purpleclay/lipgloss-theme"
	"github.com/purpleclay/nsv/internal/nsv"
)

var chevron = lipgloss.NewStyle().
	Padding(0, 1).
	Foreground(lipgloss.AdaptiveColor{
		Light: string(theme.S500),
		Dark:  string(theme.S200),
	}).
	Render(">>")

type PlaygroundOptions struct {
	Out           io.Writer
	VersionFormat string
}

func PrintFormat(tag nsv.Tag, opts PlaygroundOptions) {
	tagf := tag.Format(opts.VersionFormat)
	header := lipgloss.JoinHorizontal(lipgloss.Left, tag.Raw, chevron, opts.VersionFormat, chevron, tagf, "\n")

	pane := lipgloss.JoinVertical(lipgloss.Top,
		header,
		fmt.Sprintf("{{.Prefix}} %s%s", chevron, tag.Prefix),
		fmt.Sprintf("{{.SemVer}} %s%s", chevron, tag.SemVer),
		fmt.Sprintf("{{.Version}}%s%s", chevron, tag.Version))

	fmt.Fprint(opts.Out, pane)
}
