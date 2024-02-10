/*
Copyright (c) 2023 - 2024 Purple Clay

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
	command     = "NSV:"
	forceMajor  = "force~major"
	forceMinor  = "force~minor"
	forcePatch  = "force~patch"
	forceIgnore = "force~ignore"
	prerelease  = "pre"
)

type Command struct {
	Force      Increment
	Prerelease bool
}

func DetectCommand(log []git.LogEntry) (Command, Match) {
	force := NoIncrement
	pre := false
	match := Match{}

	for i, entry := range log {
		msg := strings.TrimSpace(entry.Message)
		idx := strings.LastIndex(msg, "\n")
		if idx == -1 {
			continue
		}

		footer := msg[idx+1:]
		if strings.ToUpper(footer[:len(command)]) != command {
			continue
		}

		cmd := strings.TrimSpace(footer[len(command):])
		match = Match{Index: i, Start: idx + 1, End: (idx + len(footer)) + 1}

		if cmd == forceIgnore {
			force = NoIncrement
			goto command
		}

		if force == MajorIncrement {
			break
		}

		switch cmd {
		case forceMajor:
			force = MajorIncrement
		case forceMinor:
			force = MinorIncrement
		case forcePatch:
			force = PatchIncrement
		case prerelease:
			pre = true
		}
	}

command:
	return Command{Force: force, Prerelease: pre}, match
}
