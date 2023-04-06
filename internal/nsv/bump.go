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

import "github.com/Masterminds/semver/v3"

// TODO: need to extract hints from commit message
// TODO: need the largest possible hint force~major, force~minor, force~patch
// TODO: Hint{Force: Increment}

// inputs: version, format, hint, increment (parallel scanning of the log) parallel execution

func Bump(ver, format string, inc Increment) (string, error) {
	pTag, _ := ParseTag(ver)
	semv, err := semver.StrictNewVersion(pTag.SemVer)
	if err != nil {
		return "", err
	}

	var bumpedVer semver.Version
	switch inc {
	case MajorIncrement:
		// TODO: if major version is 0, then fallthrough into minor
		// TODO: if hint exists to force major, then ignore previous logic
		bumpedVer = semv.IncMajor()
	case MinorIncrement:
		bumpedVer = semv.IncMinor()
	case PatchIncrement:
		bumpedVer = semv.IncPatch()
	}
	nextTag := pTag.Bump(bumpedVer.String())
	return nextTag.Format(format), nil
}
