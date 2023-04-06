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
	"strings"

	git "github.com/purpleclay/gitz"
)

const (
	colonSpace     = ": "
	featUpper      = "FEAT"
	fixUpper       = "FIX"
	breaking       = "BREAKING CHANGE: "
	breakingHyphen = "BREAKING-CHANGE: "
	breakingBang   = '!'
	noMatch        = -1
)

func DetectIncrement(log []git.LogEntry) (Increment, int) {
	mode := NoIncrement
	match := noMatch
	for i, entry := range log {
		// Check for the existence of a conventional commit type
		idx := strings.Index(entry.Message, colonSpace)
		if idx == -1 {
			continue
		}

		leadingType := strings.ToUpper(entry.Message[:idx])
		if leadingType[idx-1] == breakingBang || multilineBreaking(entry.Message) {
			return MajorIncrement, i
		}

		// Only feat and fix types now make a difference. Both have the same first letter
		if leadingType[0] != featUpper[0] {
			continue
		}

		if mode == MinorIncrement && match > noMatch {
			continue
		}

		if contains(leadingType, featUpper) {
			mode = MinorIncrement
			match = i
		} else if contains(leadingType, fixUpper) {
			mode = PatchIncrement
			match = i
		}
	}

	return mode, match
}

func contains(str, prefix string) bool {
	if str == prefix {
		return true
	}

	if strings.HasPrefix(str, prefix) {
		if len(str) > len(prefix) &&
			(str[len(prefix)] == '(' && str[len(str)-1] == ')') {
			return true
		}
	}

	return false
}

func multilineBreaking(msg string) bool {
	n := strings.Count(msg, "\n")
	if n == 0 {
		return false
	}

	idx := strings.LastIndex(msg, "\n")

	if idx == len(msg) {
		// There is a newline at the end of the string, so jump back one
		if idx = strings.LastIndex(msg[:len(msg)-1], "\n"); idx == -1 {
			return false
		}
	}

	footer := msg[idx+1:]
	return strings.HasPrefix(footer, "BREAKING CHANGE: ") ||
		strings.HasPrefix(footer, "BREAKING-CHANGE: ")
}
