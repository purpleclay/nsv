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
	"fmt"
	"io"
	"strings"

	"github.com/Masterminds/semver/v3"
	git "github.com/purpleclay/gitz"
)

type Increment int

const (
	NoIncrement Increment = iota
	PatchIncrement
	MinorIncrement
	MajorIncrement
)

type Options struct {
	StdOut io.Writer
	StdErr io.Writer
	Show   bool
}

// NextVer ...
func NextVer(gitc *git.Client, opts Options) error {
	tags, _ := gitc.Tags(git.WithSortBy(git.VersionDesc))

	// TODO: inspect the tags
	// TODO: if there is 1, then get log using range ...
	// TODO: if there are 2, then get log using range ...
	// TODO: move this logic into gitz

	// TODO: from HEAD to last tag (is this supported in gitz?)

	log, _ := gitc.Log(git.WithRefRange(tags[0], tags[1]))
	inc, pos := DetectIncrement(log.Commits)

	if inc == NoIncrement {
		return nil
	}

	// Determine the next semantic version

	// TODO: if it starts with a 'v' strip off when calculating next semver
	ver, err := semver.StrictNewVersion(tags[0])
	if err != nil {
		return err
	}

	// TODO: understand the semantic version specification

	var bumpedVer semver.Version
	switch inc {
	case MajorIncrement:
		bumpedVer = ver.IncMajor()
	case MinorIncrement:
		bumpedVer = ver.IncMinor()
	case PatchIncrement:
		bumpedVer = ver.IncPatch()
	}

	if opts.Show {
		PrintSummary(opts.StdErr, Summary{
			Tags:  tags,
			Log:   log.Commits,
			Match: pos,
		})
	} else {
		fmt.Fprintf(opts.StdOut, bumpedVer.String())
	}

	return nil
}

// DetectIncrement ...
func DetectIncrement(log []git.LogEntry) (Increment, int) {
	mode := NoIncrement
	match := -1
	for i, entry := range log {
		// Check for the existence of a conventional commit type
		idx := strings.Index(entry.Message, ": ")
		if idx == -1 {
			continue
		}

		leadingType := strings.ToUpper(entry.Message[:idx])
		if leadingType[idx-1] == '!' || multilineBreaking(entry.Message) {
			return MajorIncrement, i
		}

		// Only feat and fix types now make a difference
		if leadingType[0] != 'F' {
			continue
		}

		if mode == MinorIncrement && match > -1 {
			continue
		}

		if contains(leadingType, "FEAT") {
			mode = MinorIncrement
			match = i
		} else if contains(leadingType, "FIX") {
			mode = PatchIncrement
			match = i
		}
	}

	return mode, match
}

func contains(str string, prefix string) bool {
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
