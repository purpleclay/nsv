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

	"github.com/purpleclay/chomp"
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

func commands(line string) []string {
	var cmds []string

	rem := line
	for {
		var out string
		var err error

		// keep chomping until an error is returned
		rem, out, err = cmd()(rem)
		if err != nil {
			break
		}

		cmds = append(cmds, out)
	}

	return cmds
}

func cmd() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		rem, out, err := chomp.Not(", ")(s)
		if err != nil {
			return rem, out, err
		}

		rem, _, _ = chomp.Opt(chomp.Any(", "))(rem)
		return rem, out, nil
	}
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

		cmdLine := strings.TrimSpace(footer[len(command):])
		match = Match{Index: i, Start: idx + 1, End: (idx + len(footer)) + 1}

		cmds := commands(cmdLine)
		for _, cmd := range cmds {
			switch cmd {
			case forceMajor:
				force = MajorIncrement
			case forceMinor:
				force = MinorIncrement
			case forcePatch:
				force = PatchIncrement
			case forceIgnore:
				force = NoIncrement
			case prerelease:
				pre = true
			}
		}

		// Only want the first detected command
		break
	}

	return Command{Force: force, Prerelease: pre}, match
}
