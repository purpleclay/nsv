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
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
)

const versionFormat = "{{.Raw}}"

var versionTmpl = template.Must(template.New("default-format").Parse(versionFormat))

type UnknownTemplateFieldError struct {
	Field string
}

func (e UnknownTemplateFieldError) Error() string {
	return fmt.Sprintf("template field {{%s}} is not recognised and would resolve to an invalid semantic version", e.Field)
}

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
	checkTmpl, err := template.New("version-format").Parse(tmpl)
	if err != nil {
		return err
	}

	lookup := TemplateFields()
	fields := extractFields(tmpl)
	for _, field := range fields {
		if _, supported := lookup[field]; !supported {
			return UnknownTemplateFieldError{Field: field}
		}
	}

	return checkTmpl.Execute(io.Discard, Tag{Prefix: "ui/", SemVer: "0.1.0", Version: "v0.1.0"})
}

func extractFields(tmpl string) []string {
	scanner := bufio.NewScanner(strings.NewReader(tmpl))
	scanner.Split(TokenizeFields)

	fields := make([]string, 0)
	for scanner.Scan() {
		if field := scanner.Text(); field != "" {
			fields = append(fields, field)
		}
	}

	return fields
}

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
