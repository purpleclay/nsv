package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagSkipsImpersonationIfGitConfigExists(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "poison-ivy", "user.email", "poison-ivy@dc.com")
	gittest.CommitEmptyWithAuthor(t, "scarecrow", "scarecrow@dc.com", "feat: capture metrics for populating dashboards")

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
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

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
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

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--message", "chore: tagged {{.Tag}} from {{.PrevTag}}"})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 2)
	out := gittest.Show(t, tags[1])
	assert.Contains(t, out, "chore: tagged 0.2.0 from 0.1.1")
}

func TestTagRevertsChangesInDryRunModeForHook(t *testing.T) {
	log := `> (main, origin/main) feat: support patching files using a hook
> (tag: 0.1.0) feat: ensure diffs can be displayed in the summary`
	gittest.InitRepository(t, gittest.WithLog(log))

	execFile(t, "patch-version.sh", `#!/bin/bash
echo -n $NSV_NEXT_TAG > VERSION`)

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--dry-run", "--hook", "./patch-version.sh"})
	err := cmd.Execute()
	require.NoError(t, err)

	statuses := gittest.PorcelainStatus(t)
	assert.Empty(t, statuses)
}

func execFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o755))
}
