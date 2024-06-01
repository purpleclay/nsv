package nsv

import (
	"strings"

	git "github.com/purpleclay/gitz"
)

const (
	colonSpace     = ": "
	breaking       = "BREAKING CHANGE: "
	breakingHyphen = "BREAKING-CHANGE: "
	breakingBang   = '!'
	noMatchIdx     = -1
)

var (
	minorDefault = []string{"FEAT"}
	patchDefault = []string{"FIX"}
	NoMatch      = Match{Index: noMatchIdx}
)

type ConventionalStrategy struct {
	MajorPrefixes []string
	MinorPrefixes []string
	PatchPrefixes []string
}

func Angular() ConventionalStrategy {
	return ConventionalStrategy{
		MajorPrefixes: []string{},
		MinorPrefixes: minorDefault,
		PatchPrefixes: patchDefault,
	}
}

func AngularMerge(major, minor, patch []string) ConventionalStrategy {
	minorPrefs := minorDefault
	if len(minor) > 0 {
		minorPrefs = toUpper(minor)
	}
	patchPrefs := patchDefault
	if len(patch) > 0 {
		patchPrefs = toUpper(patch)
	}

	return ConventionalStrategy{
		MajorPrefixes: toUpper(major),
		MinorPrefixes: minorPrefs,
		PatchPrefixes: patchPrefs,
	}
}

func toUpper(s []string) []string {
	upper := make([]string, 0, len(s))
	for _, se := range s {
		upper = append(upper, strings.ToUpper(se))
	}

	return upper
}

func (s ConventionalStrategy) DetectIncrement(log []git.LogEntry) (Increment, Match) {
	mode := NoIncrement
	match := NoMatch

	for i, entry := range log {
		// Check for the existence of a conventional commit type
		idx := strings.Index(entry.Message, colonSpace)
		if idx == -1 {
			continue
		}

		leadingType := strings.ToUpper(entry.Message[:idx])

		if leadingType[idx-1] == breakingBang {
			return MajorIncrement, Match{Index: i, End: idx}
		}

		if found, start, end := multilineBreaking(entry.Message); found {
			return MajorIncrement, Match{Index: i, Start: start, End: end}
		}

		if contains(s.MajorPrefixes, leadingType) {
			return MajorIncrement, Match{Index: i, End: idx}
		}

		if mode == MinorIncrement && match.Index > noMatchIdx {
			continue
		}

		if contains(s.MinorPrefixes, leadingType) {
			mode = MinorIncrement
			match = Match{Index: i, End: idx}
		} else if contains(s.PatchPrefixes, leadingType) {
			mode = PatchIncrement
			match = Match{Index: i, End: idx}
		}
	}

	return mode, match
}

func contains(prefixes []string, str string) bool {
	if len(prefixes) == 0 {
		return false
	}

	for _, prefix := range prefixes {
		if str == prefix {
			return true
		}

		if strings.HasPrefix(str, prefix) {
			if len(str) > len(prefix) &&
				(str[len(prefix)] == '(' && str[len(str)-1] == ')') {
				return true
			}
		}
	}

	return false
}

func multilineBreaking(msg string) (bool, int, int) {
	n := strings.Count(msg, "\n")
	if n == 0 {
		return false, noMatchIdx, noMatchIdx
	}

	idx := strings.LastIndex(msg, "\n")
	if idx == len(msg) {
		// There is a newline at the end of the string, so jump back one
		if idx = strings.LastIndex(msg[:len(msg)-1], "\n"); idx == -1 {
			return false, noMatchIdx, noMatchIdx
		}
	}

	footer := msg[idx+1:]
	if strings.HasPrefix(footer, breaking) ||
		strings.HasPrefix(footer, breakingHyphen) {
		return true, idx + 1, (idx + len(breaking)) - 1
	}

	return false, noMatchIdx, noMatchIdx
}
