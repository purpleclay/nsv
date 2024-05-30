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

	cmd := tagCmd(&Options{Out: io.Discard, Logger: noopLogger})
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

	cmd := tagCmd(&Options{Out: io.Discard, Logger: noopLogger})
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

	cmd := tagCmd(&Options{Out: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--message", "chore: tagged {{.Tag}} from {{.PrevTag}}"})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 2)
	out := gittest.Show(t, tags[1])
	assert.Contains(t, out, "chore: tagged 0.2.0 from 0.1.1")
}
