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

func DetectIncrement(log []git.LogEntry) (Increment, Match) {
	mode := NoIncrement
	match := Match{Index: noMatch}

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

		// Only feat and fix types now make a difference. Both have the same first letter
		if leadingType[0] != featUpper[0] {
			continue
		}

		if mode == MinorIncrement && match.Index > noMatch {
			continue
		}

		if contains(leadingType, featUpper) {
			mode = MinorIncrement
			match = Match{Index: i, End: idx}
		} else if contains(leadingType, fixUpper) {
			mode = PatchIncrement
			match = Match{Index: i, End: idx}
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

func multilineBreaking(msg string) (bool, int, int) {
	n := strings.Count(msg, "\n")
	if n == 0 {
		return false, noMatch, noMatch
	}

	idx := strings.LastIndex(msg, "\n")
	if idx == len(msg) {
		// There is a newline at the end of the string, so jump back one
		if idx = strings.LastIndex(msg[:len(msg)-1], "\n"); idx == -1 {
			return false, noMatch, noMatch
		}
	}

	footer := msg[idx+1:]
	if strings.HasPrefix(footer, breaking) ||
		strings.HasPrefix(footer, breakingHyphen) {
		return true, idx + 1, (idx + len(breaking)) - 1
	}

	return false, noMatch, noMatch
}
