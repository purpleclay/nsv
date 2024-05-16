package nsv_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectIncrement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		commit string
		inc    nsv.Increment
		match  nsv.Match
	}{
		{
			name:   "NoIncrement",
			commit: "chore(deps): bump action from 1.0.1 to 1.1.0",
			inc:    nsv.NoIncrement,
			match:  nsv.Match{},
		},
		{
			name:   "FixPatchIncrement",
			commit: "fix: incorrect cache retrieval based on key",
			inc:    nsv.PatchIncrement,
			match:  nsv.Match{Start: 0, End: 3},
		},
		{
			name:   "FixScopedPatchIncrement",
			commit: "fix(ui): paging issue within the search table",
			inc:    nsv.PatchIncrement,
			match:  nsv.Match{Start: 0, End: 7},
		},
		{
			name:   "FeatMinorIncrement",
			commit: "feat: add new grep feature to cli",
			inc:    nsv.MinorIncrement,
			match:  nsv.Match{Start: 0, End: 4},
		},
		{
			name:   "FeatScopedMinorIncrement",
			commit: "feat(search): add fuzzy finding to predictive search",
			inc:    nsv.MinorIncrement,
			match:  nsv.Match{Start: 0, End: 12},
		},
		{
			name:   "BangMajorIncrement",
			commit: "refactor!: rename fields within cli context",
			inc:    nsv.MajorIncrement,
			match:  nsv.Match{Start: 0, End: 9},
		},
		{
			name:   "BangScopedMajorIncrement",
			commit: "fix(api)!: rewrite of existing API",
			inc:    nsv.MajorIncrement,
			match:  nsv.Match{Start: 0, End: 9},
		},
		{
			name: "BreakingFooterMajorIncrement",
			commit: `perf: rewrite sorting algorithm and public interface
BREAKING CHANGE: this is a breaking change`,
			inc:   nsv.MajorIncrement,
			match: nsv.Match{Start: 53, End: 68},
		},
		{
			name: "HyphenatedBreakingFooterMajorIncrement",
			commit: `fix(ui): incorrect sizing of components during window resize
BREAKING-CHANGE: this is a breaking change`,
			inc:   nsv.MajorIncrement,
			match: nsv.Match{Start: 61, End: 76},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inc, match := nsv.DetectIncrement([]git.LogEntry{
				{
					Message: tt.commit,
				},
			})
			require.Equal(t, tt.inc, inc, "failed to match increment")
			require.Equal(t, tt.match.Start, match.Start, "failed to match starting index")
			require.Equal(t, tt.match.End, match.End, "failed to match end index")
		})
	}
}

func TestDetectIncrementFindsHighestIncrement(t *testing.T) {
	t.Parallel()

	log := []git.LogEntry{
		{Message: "ci: support dependency scanning in workflows"},
		{Message: "feat: support dependency injection of database tier"},
		{Message: "fix: tui not redrawing correctly during terminal window resize"},
		{Message: "feat!: depecrate existing api and switch to GraphQL"},
		{
			Message: `refactor(cli): remove use of existing flags and rename
BREAKING CHANGE: this is a breaking change`,
		},
		{Message: "docs: build documentation using material for mkdocs"},
	}
	inc, match := nsv.DetectIncrement(log)

	assert.Equal(t, nsv.MajorIncrement, inc)
	assert.Equal(t, 3, match.Index)
}

func TestDetectIncrementStrictness(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		commit   string
		expected nsv.Increment
	}{
		{
			name:     "NoSpaceAfterTypeColon",
			commit:   "feat:add support for multiple search terms",
			expected: nsv.NoIncrement,
		},
		{
			name: "NoSpaceAfterBreakingColon",
			commit: `refactor: add caching layer in front of database
BREAKING CHANGE:this is a breaking change`,
			expected: nsv.NoIncrement,
		},
		{
			name: "NoSpaceAfterHyphenatedBreakingColon",
			commit: `perf: change to write through caching
BREAKING-CHANGE:this is a breaking change`,
			expected: nsv.NoIncrement,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inc, _ := nsv.DetectIncrement([]git.LogEntry{
				{
					Message: tt.commit,
				},
			})
			require.Equal(t, tt.expected, inc, "failed to match increment")
		})
	}
}

func TestDetectIncrementCaseInsensitiveLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		commit   string
		expected nsv.Increment
	}{
		{
			name:     "ForFixType",
			commit:   "FiX(ui): incorrect styling of elements on page",
			expected: nsv.PatchIncrement,
		},
		{
			name:     "ForFeatType",
			commit:   "FeAt: add tui support using bubbletea",
			expected: nsv.MinorIncrement,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inc, _ := nsv.DetectIncrement([]git.LogEntry{
				{
					Message: tt.commit,
				},
			})
			require.Equal(t, tt.expected, inc, "failed to match increment")
		})
	}
}

func TestDetectIncrementCaseSensitiveBreaking(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		commit   string
		expected nsv.Increment
	}{
		{
			name: "MixedCase",
			commit: `fix: incorrectly closing handles to database
Breaking Change: this is a breaking change`,
			expected: nsv.PatchIncrement,
		},
		{
			name: "MixedCaseHyphenated",
			commit: `feat: switched from mongodb to surrealdb
Breaking-Change: this is a breaking change`,
			expected: nsv.MinorIncrement,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inc, _ := nsv.DetectIncrement([]git.LogEntry{
				{
					Message: tt.commit,
				},
			})
			require.Equal(t, tt.expected, inc, "failed to match increment")
		})
	}
}

// globals are used to prevent any compiler optimizations
var (
	gInc   nsv.Increment
	gMatch nsv.Match
)

func BenchmarkDetectIncrement(b *testing.B) {
	log := []git.LogEntry{
		{Message: "docs: generate documentation using material for mkdocs"},
		{Message: "ci: turn on automatic analysis of code using deepsource"},
		{Message: "feat(ui): add support to drag and drop columns within table"},
		{Message: "chore: update project description within README"},
		{Message: `feat(deps): bump github.com/aws/aws-sdk-go-v2/service/route53 from 1.27.0 to 1.27.1

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`},
		{Message: "fix(parser): incorrect tokenization of conventional commit messages"},
		{Message: `chore(deps): bump actions/dependency-review-action from 2 to 3 (#156)

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`},
		{Message: "feat(store): add ability to take and skip data retrieved from the store"},
		{Message: "ci: support parallel running of go tests"},
		{Message: `refactor(api): combine existing methods to reduce complexity of api

BREAKING-CHANGE: backwards compatibility with v1 is no longer supported`},
	}
	b.ResetTimer()

	var inc nsv.Increment
	var m nsv.Match
	for n := 0; n < b.N; n++ {
		inc, m = nsv.DetectIncrement(log)
	}
	gInc = inc
	gMatch = m
}
