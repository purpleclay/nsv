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

package cmd

import (
	"io"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagSkipsImpersonationIfGitConfigExists(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "poison-ivy", "user.email", "poison-ivy@dc.com")
	gittest.CommitEmptyWithAuthor(t, "scarecrow", "scarecrow@dc.com", "feat: capture metrics for populating dashboards")

	cmd := tagCmd(&Options{Out: io.Discard})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 1)
	out := gittest.Show(t, tags[0])
	assert.Contains(t, out, "Tagger: poison-ivy <poison-ivy@dc.com>")
}

func TestTagSkipsImpersonationIfGitEnvVarsExist(t *testing.T) {
	gittest.InitRepository(t)
	gittest.CommitEmptyWithAuthor(t, "penguin", "penguin@dc.com", "feat(data): support data export from mongodb")

	t.Setenv("GIT_COMMITTER_NAME", "joker")
	t.Setenv("GIT_COMMITTER_EMAIL", "joker@dc.com")

	cmd := tagCmd(&Options{Out: io.Discard})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 1)
	out := gittest.Show(t, tags[0])
	assert.Contains(t, out, "Tagger: joker <joker@dc.com>")
}

func TestTagWithTemplatedMessage(t *testing.T) {
	log := `feat(ui): include timeline support for display historical events
(tag: 0.1.1) fix(ui): events are not being sorted as per user filters`
	gittest.InitRepository(t, gittest.WithLog(log))

	cmd := tagCmd(&Options{Out: io.Discard})
	cmd.SetArgs([]string{"--message", "chore: tagged {{.Tag}} from {{.PrevTag}}"})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 2)
	out := gittest.Show(t, tags[1])
	assert.Contains(t, out, "chore: tagged 0.2.0 from 0.1.1")
}
