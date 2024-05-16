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
