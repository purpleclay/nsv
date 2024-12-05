package nsv

import (
	"io"
	"testing"

	"github.com/charmbracelet/log"
	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var noopLogger = log.New(io.Discard)

func TestCheckAndHealRepositoryFixShallow(t *testing.T) {
	log := `(main, origin/main) docs: document new search improvements
fix(search): search is not being aggregated correctly
(tag: 0.1.0) feat(search): support aggregations for search analytics
ci: add parallel testing support to workflow`
	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithCloneDepth(1))

	client, _ := git.NewClient()
	err := checkAndHealRepository(client, Options{FixShallow: true, Logger: noopLogger})
	require.NoError(t, err)

	logEntries := gittest.Log(t)
	assert.Len(t, logEntries, 5)
}
