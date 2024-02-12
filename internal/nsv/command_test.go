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

package nsv_test

import (
	"fmt"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectCommandForce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		commit string
		inc    nsv.Increment
		match  nsv.Match
	}{
		{
			name: "PatchIncrement",
			commit: `group together filtering and sorting of components to ensure
results are displayed to the user in the expected format
nsv:force~patch`,
			inc:   nsv.PatchIncrement,
			match: nsv.Match{Start: 118, End: 133},
		},
		{
			name: "MinorIncrement",
			commit: `include code examples around test library usage
nsv: force~minor
`,
			inc:   nsv.MinorIncrement,
			match: nsv.Match{Start: 48, End: 64},
		},
		{
			name: "MajorIncrement",
			commit: `fix broken badges on README
NSV: force~major`,
			inc:   nsv.MajorIncrement,
			match: nsv.Match{Start: 28, End: 44},
		},
		{
			name: "NoIncrement",
			commit: `include smart labels for faster searching within results
nsv:force~ignore`,
			inc:   nsv.NoIncrement,
			match: nsv.Match{Start: 57, End: 73},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd, match := nsv.DetectCommand([]git.LogEntry{
				{
					Message: tt.commit,
				},
			})
			require.Equal(t, tt.inc, cmd.Force, "failed to match increment")
			require.Equal(t, tt.match.Start, match.Start, "failed to match starting index")
			require.Equal(t, tt.match.End, match.End, "failed to match end index")
			require.Empty(t, cmd.Prerelease, "prerelease should be empty")
		})
	}
}

func TestDetectCommandFirstChosen(t *testing.T) {
	t.Parallel()

	log := []git.LogEntry{
		{Message: `switch to using bootstrap for styling dashboard
nsv: force~minor`},
		{Message: `fix cache key generation as items are not being correctly evicted
nsv: force~patch`},
		{Message: `switch to using an in-memory cache that supports item eviction
nsv: force~major`},
		{Message: `bump aquasecurity/trivy-action from 0.5.0 to 0.6.2 (#12)

Bumps [aquasecurity/trivy-action](https://github.com/aquasecurity/trivy-action) from 0.5.0 to 0.6.2.
- [Release notes](https://github.com/aquasecurity/trivy-action/releases)
- [Commits](https://github.com/aquasecurity/trivy-action/compare/0.5.0...0.6.2)

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`},
	}

	cmd, _ := nsv.DetectCommand(log)
	assert.Equal(t, nsv.MinorIncrement, cmd.Force)
}

func TestDetectCommandPrerelease(t *testing.T) {
	tests := []struct {
		name    string
		command string
		label   string
	}{
		{
			name:    "DefaultToBeta",
			command: "pre",
			label:   "beta",
		},
		{
			name:    "Alpha",
			command: "pre~alpha",
			label:   "alpha",
		},
		{
			name:    "Beta",
			command: "pre~beta",
			label:   "beta",
		},
		{
			name:    "ReleaseCandidate",
			command: "pre~rc",
			label:   "rc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd, _ := nsv.DetectCommand([]git.LogEntry{
				{
					Message: fmt.Sprintf(`experimental feature has been added to search
nsv:%s`, tt.command),
				},
			})
			require.Equal(t, tt.label, cmd.Prerelease, "failed to match prerelease label")
		})
	}
}

func TestDetectMultipleCommands(t *testing.T) {
	t.Parallel()

	cmd, match := nsv.DetectCommand([]git.LogEntry{
		{
			Message: `experimental use of a file cache
nsv:pre,force~major`,
		},
	})
	assert.Equal(t, "beta", cmd.Prerelease)
	assert.Equal(t, nsv.MajorIncrement, cmd.Force)
	assert.Equal(t, 33, match.Start)
	assert.Equal(t, 52, match.End)
}

// globals are used to prevent any compiler optimizations
var gCmd nsv.Command

func BenchmarkDetectCommand(b *testing.B) {
	log := []git.LogEntry{
		{Message: "generate documentation using material for mkdocs"},
		{Message: "turn on automatic analysis of code using deepsource"},
		{Message: `add support to drag and drop columns within table
nsv: force~patch`},
		{Message: "update project description within README"},
		{Message: `bump github.com/aws/aws-sdk-go-v2/service/route53 from 1.27.0 to 1.27.1

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`},
		{Message: "incorrect tokenization of conventional commit messages"},
		{Message: `bump actions/dependency-review-action from 2 to 3 (#156)

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`},
		{Message: "add ability to take and skip data retrieved from the store"},
		{Message: "support parallel running of go tests"},
		{Message: `combine existing methods to reduce complexity of api

BREAKING-CHANGE: backwards compatibility with v1 is no longer supported
nsv:force~minor`},
	}
	b.ResetTimer()

	var cmd nsv.Command
	var match nsv.Match
	for n := 0; n < b.N; n++ {
		cmd, match = nsv.DetectCommand(log)
	}
	gCmd = cmd
	gMatch = match
}
