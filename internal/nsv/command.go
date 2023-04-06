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
	command     = "NSV:"
	forceMajor  = "force~major"
	forceMinor  = "force~minor"
	forcePatch  = "force~patch"
	forceIgnore = "force~ignore"
)

type Command struct {
	Force Increment
}

func DetectCommand(log []git.LogEntry) Command {
	force := NoIncrement
	for _, entry := range log {
		eol := strings.Index(entry.Message, "\n")
		if eol == -1 || eol == len(entry.Message) {
			continue
		}

		footers := strings.SplitAfter(entry.Message[eol:], "\n")
		for _, footer := range footers {
			if len(footer) < len(command) || strings.ToUpper(footer[:len(command)]) != command {
				continue
			}

			command := strings.TrimSpace(footer[len(command):])
			if command == forceIgnore {
				force = NoIncrement
				goto command
			}

			if force == MajorIncrement {
				break
			}

			switch command {
			case forceMajor:
				force = MajorIncrement
			case forceMinor:
				force = MinorIncrement
			case forcePatch:
				force = PatchIncrement
			}
		}
	}

command:
	return Command{Force: force}
}
