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

	"github.com/Masterminds/semver/v3"
	git "github.com/purpleclay/gitz"
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
