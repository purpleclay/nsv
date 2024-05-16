package nsv_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/stretchr/testify/assert"
)

func TestCheckTemplateUnsupportedField(t *testing.T) {
	t.Parallel()

	err := nsv.CheckTemplate("{{ .Unknown }}")
	assert.EqualError(t, err, "template field {{.Unknown}} is not recognised and would resolve to an invalid semantic version")
}

func TestTokenizeFields(t *testing.T) {
	t.Parallel()

	scanner := bufio.NewScanner(strings.NewReader("{{ .A }}{{.B}}{{  func .C.D  }}"))
	scanner.Split(nsv.TokenizeFields)

	fields := readUntilEOF(t, scanner)
	assert.ElementsMatch(t, []string{".A", ".B", ".C.D"}, fields)
}

func readUntilEOF(t *testing.T, scanner *bufio.Scanner) []string {
	t.Helper()

	fields := make([]string, 0)
	for scanner.Scan() {
		fields = append(fields, scanner.Text())
	}

	return fields
}
