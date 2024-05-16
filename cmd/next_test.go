package cmd

import (
	"bytes"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNext(t *testing.T) {
	log := `(main, origin/main) docs: document new pagination improvements
feat: support pagination of search results`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	cmd := nextCmd(&Options{Out: &buf})
	err := cmd.Execute()

	require.NoError(t, err)
	assert.Equal(t, "0.1.0", buf.String())
}
