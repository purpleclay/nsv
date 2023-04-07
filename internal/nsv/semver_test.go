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
	"bytes"
	"os"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextVersion(t *testing.T) {
	log := `(main, origin/main) docs: document new search improvements
fix(search): search is not being aggregated correctly
(tag: 0.1.0) feat(search): support aggregations for search analytics
ci: add parallel testing support to workflow`
	gittest.InitRepository(t, gittest.WithLog(log))
	gitc, _ := git.NewClient()

	var buf bytes.Buffer
	err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

	require.NoError(t, err)
	assert.Equal(t, "0.1.1", buf.String())
}

func TestNextVersionFirstVersion(t *testing.T) {
	tests := []struct {
		name     string
		log      string
		expected string
	}{
		{
			name:     "Patch",
			log:      "fix: incorrectly displaying results in tui",
			expected: "0.0.1",
		},
		{
			name:     "Minor",
			log:      "feat: add new tui for displaying job summary",
			expected: "0.1.0",
		},
		{
			name:     "Major",
			log:      "feat!: somehow we have a breaking change for the first commit",
			expected: "0.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t, gittest.WithLog(tt.log))
			gitc, _ := git.NewClient()

			var buf bytes.Buffer
			err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestNextVersionMajorZeroSemV(t *testing.T) {
	log := `(main, origin/main) refactor!: switch to v0.1 of the schema which has no backwards compatibility
(tag: 0.1.2) fix: issues with v0.0.8 of schema`
	gittest.InitRepository(t, gittest.WithLog(log))
	gitc, _ := git.NewClient()

	var buf bytes.Buffer
	err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

	require.NoError(t, err)
	require.Equal(t, "0.2.0", buf.String())
}

func TestNextVersionMajorZeroSemVForceMajor(t *testing.T) {
	log := `> (main, origin/main) feat!: everything is now stable ready for v1
nsv: force~major
> (tag: 0.9.9) fix: stability issues around long running database connectivity`
	gittest.InitRepository(t, gittest.WithLog(log))
	gitc, _ := git.NewClient()

	var buf bytes.Buffer
	err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

	require.NoError(t, err)
	require.Equal(t, "1.0.0", buf.String())
}

func TestNextVersionIncludesSubDirectoryAsPrefix(t *testing.T) {
	log := `feat(trends)!: breaking change to capturing user trends`
	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithFiles("search/main.go", "store/main.go"))
	gittest.StageFile(t, "search/main.go")
	gittest.Commit(t, "feat(search): add support to search across user trends for recommendations")
	gittest.StageFile(t, "store/main.go")
	gittest.Commit(t, "fix(store): fixed timestamp formatting issues")

	os.Chdir("search")
	gitc, _ := git.NewClient()

	var buf bytes.Buffer
	err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

	require.NoError(t, err)
	assert.Equal(t, "search/0.1.0", buf.String())
}

func TestNextVersionPreservesTagPrefix(t *testing.T) {
	tests := []struct {
		name     string
		log      string
		expected string
	}{
		{
			name: "VPrefix",
			log: `(main, origin/main) feat(ui): rebrand existing components with new theme
docs: update documentation with breaking change
(tag: v1.0.0) feat(ui)!: breaking redesign of search ui`,
			expected: "v1.1.0",
		},
		{
			name: "MonoRepoPrefix",
			log: `(HEAD -> main) fix(cache): incorrect sorting of metadata cache
docs: update documentation with latest cache improvement
ci: configure workflow to run benchmarks
(tag: cache/v0.2.0) feat(cache): `,
			expected: "cache/v0.2.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t, gittest.WithLog(tt.log))
			gitc, _ := git.NewClient()

			var buf bytes.Buffer
			err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestNextVersionWithFormat(t *testing.T) {
	format := "v{{ .Version }}"

	tests := []struct {
		name     string
		log      string
		expected string
	}{
		{
			name:     "NewVersion",
			log:      "(main) feat(broker): support asynchronous publishing to broker",
			expected: "v0.1.0",
		},
		{
			name: "ExistingVersion",
			log: `(main) feat(broker): enable tls authentication
(tag: v0.1.0) feat(broker): support asynchronous publishing to broker`,
			expected: "v0.2.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t, gittest.WithLog(tt.log))
			gitc, _ := git.NewClient()

			var buf bytes.Buffer
			err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf, VersionFormat: format})

			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}
}
