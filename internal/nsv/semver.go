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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	git "github.com/purpleclay/gitz"
)

const (
	vPrefix       = 'v'
	firstVer      = "0.0.0"
	versionFormat = "{{ .Prefix }}{{ .Version }}"
)

var versionTmpl = template.Must(template.New("version").Parse(versionFormat))

type Options struct {
	StdOut io.Writer
	StdErr io.Writer
	Show   bool
}

type context struct {
	TagPrefix string
	LogPath   string
}

func execContext(gitc *git.Client) (*context, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	relPath, err := gitc.ToRelativePath(cwd)
	if err != nil {
		return nil, err
	}

	if relPath == git.RelativeAtRoot {
		return &context{TagPrefix: "", LogPath: ""}, nil
	}

	logPath := relPath
	if logPath == filepath.Base(cwd) {
		logPath = ""
	}

	return &context{TagPrefix: relPath, LogPath: logPath}, nil
}

type tag struct {
	Prefix  string
	SemVer  string
	Version string
}

func parseTag(raw string) tag {
	lastSlash := 0
	if idx := strings.LastIndex(raw, "/"); idx > -1 {
		lastSlash = idx + 1
	}

	semv := raw[lastSlash:]
	if semv[0] == vPrefix {
		semv = semv[1:]
	}

	return tag{
		Prefix:  raw[:lastSlash],
		SemVer:  semv,
		Version: raw[lastSlash:],
	}
}

func (t tag) bump(semv string) tag {
	ver := semv
	if t.Version[0] == vPrefix {
		ver = fmt.Sprintf("%c%s", vPrefix, ver)
	}

	return tag{
		Prefix:  t.Prefix,
		SemVer:  semv,
		Version: ver,
	}
}

func NextVersion(gitc *git.Client, opts Options) error {
	ctx, err := execContext(gitc)
	if err != nil {
		return err
	}

	ltag, err := latestTag(gitc, ctx.TagPrefix)
	if err != nil {
		return err
	}

	log, err := gitc.Log(git.WithPaths(ctx.LogPath), git.WithRefRange(git.HeadRef, ltag))
	if err != nil {
		return err
	}

	inc, pos := DetectIncrement(log.Commits)
	if inc == NoIncrement {
		return nil
	}

	if ltag == "" {
		ltag = firstVersion(ctx.TagPrefix)
	}

	pTag := parseTag(ltag)
	ver, err := semver.StrictNewVersion(pTag.SemVer)
	if err != nil {
		return err
	}

	var bumpedVer semver.Version
	switch inc {
	case MajorIncrement:
		bumpedVer = ver.IncMajor()
	case MinorIncrement:
		bumpedVer = ver.IncMinor()
	case PatchIncrement:
		bumpedVer = ver.IncPatch()
	}
	nextTag := pTag.bump(bumpedVer.String())

	var nextV bytes.Buffer
	versionTmpl.Execute(&nextV, nextTag)
	fmt.Fprint(opts.StdOut, nextV.String())

	if opts.Show {
		PrintSummary(opts.StdErr, Summary{
			Tags:  []string{git.HeadRef, nextV.String()},
			Log:   log.Commits,
			Match: pos,
		})
	}
	return nil
}

func latestTag(gitc *git.Client, prefix string) (string, error) {
	prefixFilter := func(tag string) bool {
		if prefix == "" {
			return true
		}

		return strings.HasPrefix(tag, prefix+"/")
	}

	tags, err := gitc.Tags(git.WithShellGlob("**/*.*.*"),
		git.WithSortBy(git.VersionDesc),
		git.WithFilters(prefixFilter),
		git.WithCount(1))
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", nil
	}

	return tags[0], nil
}

func firstVersion(prefix string) string {
	if prefix == "" {
		return firstVer
	}
	return fmt.Sprintf("%s/%s", prefix, firstVer)
}
