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

package nsv

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"strings"
	"text/template"
)

const versionFormat = "{{.Prefix}}{{.Version}}"

var versionTmpl = template.Must(template.New("default-format").Parse(versionFormat))

func TemplateFields() map[string]struct{} {
	fields := map[string]struct{}{}

	rType := reflect.TypeOf(Tag{})
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		if field.IsExported() {
			fields["."+field.Name] = struct{}{}
		}
	}

	return fields
}

func CheckTemplate(tmpl string) error {
	if _, err := template.New("version-format").Parse(tmpl); err != nil {
		return err
	}

	lookup := TemplateFields()
	fields := extractFields(tmpl)
	for _, field := range fields {
		if _, supported := lookup[field]; !supported {
			// TODO: custom error that should be raised (reference in documentation)
			return errors.New("template field is not supported and will resolve to an invalid semantic version")
		}
	}

	// TODO: execute the template against a default tag that exercises all of the runs
	return nil
}

func extractFields(tmpl string) []string {
	scanner := bufio.NewScanner(strings.NewReader(tmpl))
	scanner.Split(TokenizeFields)

	fields := make([]string, 0)
	for scanner.Scan() {
		fields = append(fields, scanner.Text())
	}

	return fields
}

// TODO: write basic test to ensure this works
func TokenizeFields(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte{'{', '{'}); i >= 0 {
		if j := bytes.Index(data, []byte{'}', '}'}); j >= i {
			return j + 2, eatField(data[i+2 : j]), nil
		}
	}

	if atEOF {
		return len(data), nil, nil
	}

	return 0, nil, nil
}

func eatField(data []byte) []byte {
	dotPos := bytes.IndexByte(data, '.')
	if dotPos == -1 {
		return nil
	}

	i := dotPos
	for i < len(data) && data[i] != ' ' {
		i++
	}

	return data[dotPos:i]
}
