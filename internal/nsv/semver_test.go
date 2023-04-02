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

// TODO: fix commits

func TestNextVersion(t *testing.T) {
	log := `(main, origin/main) docs: bbb
fix: aaa
(tag: 0.1.0) feat: ccc
ci: ddd`
	gittest.InitRepository(t, gittest.WithLog(log))
	gitc, _ := git.NewClient()

	var buf bytes.Buffer
	err := nsv.NextVersion(gitc, nsv.Options{StdOut: &buf})

	require.NoError(t, err)
	assert.Equal(t, "0.1.1", buf.String())
}

// TODO: breaking not increment major {0.1.0 feat!, 0.1.0 feat:, 0.0.1 patch}

func TestNextVersionFirstVersion(t *testing.T) {
	tests := []struct {
		name     string
		log      string
		expected string
	}{
		{
			name:     "Patch",
			log:      "fix: this is a fix",
			expected: "0.0.1",
		},
		{
			name:     "Minor",
			log:      "feat: this is a feature",
			expected: "0.1.0",
		},
		// TODO: handle major (should adhere to semantic spec)
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

// TODO: rename ~ Appends prefix automatically

func TestNextVersionMonoRepoFirstVersion(t *testing.T) {
	log := `ci: sdsdsd
chore: sdsdsdsd
docs: dsfjkhsdljfhsdfljkh`
	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithFiles("search/main.go", "store/main.go"))
	gittest.StageFile(t, "search/main.go")
	gittest.Commit(t, "feat(search): awesome search functionality")
	gittest.StageFile(t, "store/main.go")
	gittest.Commit(t, "fix(store): fix bug in the store")

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
			log: `(main, origin/main) feat: aaaa
docs: sdsd
chore: sdsd
(tag: v1.0.0) feat: bbbb`,
			expected: "v1.1.0",
		},
		{
			// TODO: what does a mono-repo log look like for a directory
			name: "MonoRepoPrefix",
			log: `(main, origin/main) fix: dhdhdhdhdh
docs: hhfhfhfhf
ci: hhdhdhdh
(tag: component/v0.2.0) feat: sdlkjfhsdljhsdlkfjhldfk`,
			expected: "component/v0.2.1",
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
