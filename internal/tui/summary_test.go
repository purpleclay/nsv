package tui_test

import (
	"bytes"
	"testing"

	git "github.com/purpleclay/gitz"
	theme "github.com/purpleclay/lipgloss-theme"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/purpleclay/nsv/internal/tui"
	"gotest.tools/v3/golden"
)

// Thanks ChatGPT!
// The Madness Rating is an interpretation of each character's level of psychological instability
// and how unpredictable or mentally unhinged they are portrayed
var data = [][]string{
	{"Name", "Sex", "Distinguishing Features", "Madness Rating"},
	{"The Joker", "Male", "Clown-like appearance, green hair, pale skin, psychopathic smile", "10"},
	{"Harley Quinn", "Female", "Clown-like appearance, mallet weapon, acrobatic and unpredictable", "9"},
	{"Two-Face", "Male", "Half-burned face, split personality (Harvey Dent and Two-Face)", "8"},
	{"Scarecrow", "Male", "Wears a scarecrow mask, uses fear toxins to manipulate victims", "8"},
	{"Mad Hatter", "Male", "Obsession with Alice in Wonderland, mind-control technology", "8"},
	{"Riddler", "Male", "Obsession with riddles, green suit with question marks", "7"},
}

var versions = []*nsv.Next{
	{
		Tag:     "0.2.0",
		PrevTag: "0.1.0",
		LogDir:  "src/ui",
		Log: []git.LogEntry{
			{
				Hash:       "ba1ec8339c0fc97a6e53147aa916308364c8d420",
				AbbrevHash: "ba1ec83",
				Message:    "fix: search options were not being correctly converted into elastic search filters (#63)",
			},
			{
				Hash:       "4e7a2774f16bccdb2da56488f3c2b7ff9eea09b6",
				AbbrevHash: "4e7a277",
				Message: `chore(deps): bump docker/setup-qemu-action from 2 to 3 (#62)

Signed-off-by: dependabot[bot] <support@github.com>
 Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`,
			},
			{
				Hash:       "2c9b178b073e8314603ca942c56c35feda7a14cf",
				AbbrevHash: "2c9b178",
				Message:    "feat: add option toggles to the dashboard that allows dynamic queryies to elastic (#58)",
			},
		},
		Match: nsv.Match{
			Index: 2,
			Start: 0,
			End:   4,
		},
	},
	{
		Tag:     "0.2.1",
		PrevTag: "0.2.0",
		LogDir:  "src/search",
		Log: []git.LogEntry{
			{
				Hash:       "6e6fcac17ac90f351b66006cdc9909ae8cb84479",
				AbbrevHash: "6e6fcac",
				Message:    "feat: add redis caching support (#55)",
			},
			{
				Hash:       "869fd31e9d658ea41f6ca61b8c69026453c455c2",
				AbbrevHash: "869fd31",
				Message: `feat(deps): bump github.com/charmbracelet/lipgloss from 0.7.1 to 0.8.0 (#56)

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`,
			},
		},
		Match: nsv.Match{
			Index: 0,
			Start: 0,
			End:   4,
		},
	},
}

func TestPrintSummary(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	tui.PrintSummary(versions, tui.SummaryOptions{Out: &buf})

	golden.Assert(t, buf.String(), "TestPrintSummary.golden")
}

func TestPrintSummaryWithDiffs(t *testing.T) {
	t.Parallel()

	versionsWithDiffs := copyVersions(t)
	versionsWithDiffs[0].Diffs = []git.FileDiff{
		{Path: "src/ui/go.mod"},
		{Path: "src/ui/dashboard.go"},
	}
	versionsWithDiffs[1].Diffs = []git.FileDiff{
		{Path: "src/search/cache.go"},
		{Path: "src/search/go.mod"},
	}

	var buf bytes.Buffer
	tui.PrintSummary(versionsWithDiffs, tui.SummaryOptions{Out: &buf})

	golden.Assert(t, buf.String(), "TestPrintSummaryWithDiffs.golden")
}

func copyVersions(t *testing.T) []*nsv.Next {
	t.Helper()

	versionsWithDiffs := make([]*nsv.Next, len(versions))
	for i := range versions {
		ver := *versions[i]
		versionsWithDiffs[i] = &ver
	}

	return versionsWithDiffs
}

func TestPrintSummaryNoColor(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	opts := tui.SummaryOptions{
		Out:     &buf,
		NoColor: true,
	}

	tui.PrintSummary(versions, opts)

	golden.Assert(t, buf.String(), "TestPrintSummaryNoColor.golden")
}

func TestPrintSummaryCompact(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	opts := tui.SummaryOptions{
		Out:    &buf,
		Pretty: tui.Compact,
	}

	tui.PrintSummary(versions, opts)

	golden.Assert(t, buf.String(), "TestPrintSummaryCompact.golden")
}

func TestTableNoBorder(t *testing.T) {
	t.Parallel()
	tbl := theme.NewTable(data)

	golden.Assert(t, tbl.String(), "TestTableNoBorder.golden")
}

func TestTableThinBorder(t *testing.T) {
	t.Parallel()
	tbl := theme.NewTable(data).
		Border(theme.ThinBorder)

	golden.Assert(t, tbl.String(), "TestTableThinBorder.golden")
}

func TestTableNoDividers(t *testing.T) {
	t.Parallel()
	tbl := theme.NewTable(data).
		Border(theme.ThinBorder).
		Dividers(false)

	golden.Assert(t, tbl.String(), "TestTableNoDividers.golden")
}

func TestTableCollapsed(t *testing.T) {
	t.Parallel()
	tbl := theme.NewTable(data).
		Collapsed(true)

	golden.Assert(t, tbl.String(), "TestTableCollapsed.golden")
}

func TestTable(t *testing.T) {
	golden.Assert(t, theme.U.Render("sssss"), "TestTable.golden")
}
