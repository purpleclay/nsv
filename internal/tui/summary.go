package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	theme "github.com/purpleclay/lipgloss-theme"
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

	faint    = lipgloss.NewStyle().Faint(true)
	bullet   = faint.SetString(">")
	padRight = lipgloss.NewStyle().PaddingRight(1)
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
			padRight.Render(marker),
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
		padRight.Render(theme.Tick),
		lipgloss.JoinVertical(
			lipgloss.Top,
			theme.Mark.Render(entry.AbbrevHash),
			wordwrap.String(msg, logWrapAt)),
	)
}
