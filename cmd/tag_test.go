package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	log := `feat: ensure logs are attached to spans during tracing
(tag: 0.1.0) feat: support distributed tracing`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("VERSION"),
		gittest.WithFileContent("VERSION", "0.1.0"),
	)

	execFile(t, "patch-version.sh", `#!/bin/bash
echo -n $NSV_NEXT_TAG > VERSION`)

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--hook", "./patch-version.sh"})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.ElementsMatch(t, []string{"0.1.0", "0.2.0"}, tags)

	logs := gittest.Log(t)
	assert.Equal(t, "chore: tagged release 0.2.0 [skip ci]", logs[0].Message)

	out := gittest.Show(t, "0.2.0")
	assert.Contains(t, out, logs[0].Hash)
}

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

func TestTagWithTemplatedTagMessage(t *testing.T) {
	log := `feat(ui): include timeline support for display historical events
(tag: 0.1.1) fix(ui): events are not being sorted as per user filters`
	gittest.InitRepository(t, gittest.WithLog(log))

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--tag-message", "chore: tagged {{.Tag}} from {{.PrevTag}}"})
	err := cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 2)
	out := gittest.Show(t, tags[1])
	assert.Contains(t, out, "chore: tagged 0.2.0 from 0.1.1")
}

func TestTagWithTemplatedCommitMessage(t *testing.T) {
	log := `feat(search): include prebuilt queries for common search scenarios
(tag: 0.2.1) fix(search): prevent truncation of search criteria entered by users`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("VERSION"),
		gittest.WithFileContent("VERSION", "0.1.0"),
	)

	execFile(t, "patch-version.sh", `#!/bin/bash
echo -n $NSV_NEXT_TAG > VERSION`)

	cmd := tagCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--commit-message", "chore: tagged {{.Tag}} from {{.PrevTag}}", "--hook", "./patch-version.sh"})
	err := cmd.Execute()
	require.NoError(t, err)

	logs := gittest.Log(t)
	assert.Equal(t, "chore: tagged 0.3.0 from 0.2.1", logs[0].Message)
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
