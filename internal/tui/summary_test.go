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

package tui_test

import (
	"bytes"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/purpleclay/nsv/internal/tui"
	"gotest.tools/v3/golden"
)

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
	var buf bytes.Buffer
	tui.PrintSummary(versions, tui.SummaryOptions{Out: &buf})

	golden.Assert(t, buf.String(), "TestPrintSummary.golden")
}

func TestPrintSummaryNoColor(t *testing.T) {
	var buf bytes.Buffer
	opts := tui.SummaryOptions{
		Out:     &buf,
		NoColor: true,
	}

	tui.PrintSummary(versions, opts)

	golden.Assert(t, buf.String(), "TestPrintSummaryNoColor.golden")
}

func TestPrintSummaryCompact(t *testing.T) {
	var buf bytes.Buffer
	opts := tui.SummaryOptions{
		Out:    &buf,
		Pretty: tui.Compact,
	}

	tui.PrintSummary(versions, opts)

	golden.Assert(t, buf.String(), "TestPrintSummaryCompact.golden")
}
