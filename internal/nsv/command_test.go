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

package nsv_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectCommandForce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		commit   string
		expected nsv.Increment
	}{
		{
			name: "PatchIncrement",
			commit: `group together filtering and sorting components
nsv:force~patch`,
			expected: nsv.PatchIncrement,
		},
		{
			name: "MinorIncrement",
			commit: `include code examples around test library usage
nsv: force~minor`,
			expected: nsv.MinorIncrement,
		},
		{
			name: "MajorIncrement",
			commit: `fix broken badges on README
NSV: force~major`,
			expected: nsv.MajorIncrement,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := nsv.DetectCommand([]git.LogEntry{
				{
					Message: tt.commit,
				},
			})
			require.Equal(t, tt.expected, cmd.Force, "failed to match increment")
		})
	}
}

func TestDetectCommandForceHighestChosen(t *testing.T) {
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

	cmd := nsv.DetectCommand(log)
	assert.Equal(t, nsv.MajorIncrement, cmd.Force)
}

func TestDetectCommandForceIgnore(t *testing.T) {
	t.Parallel()

	log := []git.LogEntry{
		{Message: `support the use of websockets to provide dynamic notifications to dashboard
Closes #8
nsv: force~major`},
		{Message: `add additional checks to ensure no secrets escape the repository
Closes #7
nsv: force~ignore`},
	}

	cmd := nsv.DetectCommand(log)
	assert.Equal(t, nsv.NoIncrement, cmd.Force)
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
nsv: force~minor`},
	}
	b.ResetTimer()

	var cmd nsv.Command
	for n := 0; n < b.N; n++ {
		cmd = nsv.DetectCommand(log)
	}
	gCmd = cmd
}
