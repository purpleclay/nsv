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
