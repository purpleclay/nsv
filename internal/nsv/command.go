package nsv

import (
	"strings"

	"github.com/purpleclay/chomp"
	git "github.com/purpleclay/gitz"
)

const (
	prefix      = "NSV:"
	sep         = "~"
	forceCmd    = "force"
	forceMajor  = "major"
	forceMinor  = "minor"
	forcePatch  = "patch"
	forceIgnore = "ignore"
	preCmd      = "pre"
	preAlpha    = "alpha"
	preBeta     = "beta"
	preRc       = "rc"
)

type Command struct {
	Force      Increment
	Prerelease string
}

func DetectCommand(log []git.LogEntry) (Command, Match) {
	command := Command{}
	match := Match{Index: noMatch}

	for i, entry := range log {
		msg := strings.TrimSpace(entry.Message)
		idx := strings.LastIndex(msg, "\n")
		if idx == -1 {
			continue
		}

		footer := msg[idx+1:]
		if strings.ToUpper(footer[:len(prefix)]) != prefix {
			continue
		}

		cmdLine := strings.TrimSpace(footer[len(prefix):])
		match = Match{Index: i, Start: idx + 1, End: (idx + len(footer)) + 1}

		cmds := commands(cmdLine)
		for _, cmd := range cmds {
			if strings.HasPrefix(cmd, forceCmd) {
				command.Force = chompForce(cmd)
			} else if strings.HasPrefix(cmd, preCmd) {
				command.Prerelease = chompPre(cmd)
			}
		}

		// Detect only the first command
		break
	}

	return command, match
}

func commands(line string) []string {
	_, ext, _ := chomp.ManyN(
		chomp.Suffixed(
			chomp.Opt(chomp.Any(", ")),
			chomp.Not(", ")), 0)(line)

	return ext
}

func chompForce(cmd string) Increment {
	_, out, err := chomp.SepPair(
		chomp.Tag(forceCmd),
		chomp.Tag(sep),
		chomp.First(
			chomp.Tag(forceMajor),
			chomp.Tag(forceMinor),
			chomp.Tag(forcePatch),
			chomp.Tag(forceIgnore),
		))(cmd)
	if err != nil {
		return NoIncrement
	}

	switch out[1] {
	case forceMajor:
		return MajorIncrement
	case forceMinor:
		return MinorIncrement
	case forcePatch:
		return PatchIncrement
	default:
		return NoIncrement
	}
}

func chompPre(cmd string) string {
	_, out, err := chomp.SepPair(
		chomp.Tag(preCmd),
		chomp.Tag(sep),
		chomp.First(
			chomp.Tag(preAlpha),
			chomp.Tag(preBeta),
			chomp.Tag(preRc),
		))(cmd)
	if err != nil {
		return "beta"
	}

	return out[1]
}
