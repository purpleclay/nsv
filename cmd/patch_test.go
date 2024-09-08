package cmd

import (
	"io"
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
		gittest.WithCommittedFiles("VERSION"),
		gittest.WithFileContent("VERSION", "0.1.0"),
	)

	execFile(t, "patch-version.sh", `#!/bin/bash
echo -n $NSV_NEXT_TAG > VERSION`)

	cmd := patchCmd(&Options{Out: io.Discard, Err: io.Discard, Logger: noopLogger})
	cmd.SetArgs([]string{"--hook", "./patch-version.sh"})
	err := cmd.Execute()
	require.NoError(t, err)

	logs := gittest.Log(t)
	assert.Equal(t, "chore: patched files for release 0.2.0", logs[0].Message)
}
