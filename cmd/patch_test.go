package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatch(t *testing.T) {
	log := `feat: ensure redis perists to disk to prevent loss of data during restarts
(tag: 0.1.0) feat: distributed caching`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("Chart.yaml"),
		gittest.WithFileContent("Chart.yaml", "version: 0.1.0"),
	)

	execFile(t, "patch.sh", `#!/bin/bash
sed -i "s/^version\: $NSV_PREV_TAG/version\: $NSV_NEXT_TAG/" "${NSV_WORKING_DIRECTORY}/Chart.yaml"`)

	cmd := patchCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--hook", "./patch.sh"})
	err := cmd.Execute()
	require.NoError(t, err)

	logs := gittest.RemoteLog(t)
	assert.Equal(t, "chore: patched files for release 0.2.0 [skip ci]", logs[0].Message)
	assert.Equal(t, "version: 0.2.0", readFile(t, "Chart.yaml"))
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(contents)
}
