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
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	git "github.com/purpleclay/gitz"
)

const (
	vPrefix  = 'v'
	firstVer = "0.0.0"
)

type Increment int

const (
	NoIncrement Increment = iota
	PatchIncrement
	MinorIncrement
	MajorIncrement
)

type Options struct {
	Err           io.Writer `env:"-"`
	NoColor       bool      `env:"NO_COLOR"`
	Out           io.Writer `env:"-"`
	Show          bool      `env:"NSV_SHOW"`
	TagMessage    string    `env:"NSV_TAG_MESSAGE"`
	VersionFormat string    `env:"NSV_FORMAT"`
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
	if strings.HasSuffix(cwd, logPath) {
		logPath = git.RelativeAtRoot
	}

	return &context{TagPrefix: relPath, LogPath: logPath}, nil
}

type Tag struct {
	Prefix  string
	SemVer  string
	Version string
	raw     string
}

func ParseTag(raw string) (Tag, error) {
	lastSlash := 0
	if idx := strings.LastIndex(raw, "/"); idx > -1 {
		lastSlash = idx + 1
	}

	semv := raw[lastSlash:]
	if semv[0] == vPrefix {
		semv = semv[1:]
	}

	if _, err := semver.StrictNewVersion(semv); err != nil {
		return Tag{}, err
	}

	return Tag{
		Prefix:  raw[:lastSlash],
		raw:     raw,
		SemVer:  semv,
		Version: raw[lastSlash:],
	}, nil
}

func (t Tag) Bump(semv string) Tag {
	ver := semv
	if t.Version[0] == vPrefix {
		ver = fmt.Sprintf("%c%s", vPrefix, ver)
	}

	return Tag{
		Prefix:  t.Prefix,
		SemVer:  semv,
		Version: ver,
	}
}

func (t Tag) Format(format string) string {
	var tagf bytes.Buffer
	if format == "" {
		_ = versionTmpl.Execute(&tagf, t)
		return tagf.String()
	}

	tmpl, _ := template.New("custom-format").Parse(format)
	_ = tmpl.Execute(&tagf, t)

	fmted := tagf.String()
	lastSlash := 0
	if idx := strings.LastIndex(fmted, "/"); idx > -1 {
		lastSlash = idx + 1
	}

	if fmted[lastSlash:lastSlash+2] == "vv" {
		return fmted[:lastSlash] + fmted[lastSlash+1:]
	}
	return fmted
}

type Next struct {
	Log     []git.LogEntry
	LogDir  string
	Match   Match
	PrevTag string
	Tag     string
}

type Match struct {
	End   int
	Index int
	Start int
}

func NextVersion(gitc *git.Client, opts Options) (*Next, error) {
	ctx, err := execContext(gitc)
	if err != nil {
		return nil, err
	}

	ltag, err := latestTag(gitc, ctx.TagPrefix)
	if err != nil {
		return nil, err
	}

	log, err := gitc.Log(git.WithPaths(ctx.LogPath), git.WithRefRange(git.HeadRef, ltag))
	if err != nil {
		return nil, err
	}

	// Detect commands first as they have a higher precedence over conventional commits
	var inc Increment
	cmd, match := DetectCommand(log.Commits)

	inc = cmd.Force
	if inc == NoIncrement {
		inc, match = DetectIncrement(log.Commits)
	}
	if inc == NoIncrement {
		return nil, nil
	}

	if ltag == "" {
		ltag = firstVersion(ctx.TagPrefix)
	}

	nextVer, err := bump(ltag, opts.VersionFormat, inc, cmd)
	if err != nil {
		return nil, err
	}

	return &Next{
		PrevTag: ltag,
		Tag:     nextVer,
		Log:     log.Commits,
		LogDir:  ctx.TagPrefix,
		Match:   match,
	}, nil
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

func bump(ver, format string, inc Increment, cmd Command) (string, error) {
	pTag, _ := ParseTag(ver)
	semv, err := semver.StrictNewVersion(pTag.SemVer)
	if err != nil {
		return "", err
	}

	var bumpedVer semver.Version
	if inc == MajorIncrement && semv.Major() == 0 {
		// Support SemVer Major 0 (0.y.z) workflow, https://semver.org/#spec-item-4
		inc = MinorIncrement
	}

	if cmd.Force != NoIncrement {
		inc = cmd.Force
	}

	switch inc {
	case MajorIncrement:
		bumpedVer = semv.IncMajor()
	case MinorIncrement:
		bumpedVer = semv.IncMinor()
	case PatchIncrement:
		bumpedVer = semv.IncPatch()
	}

	nextTag := pTag.Bump(bumpedVer.String())
	return nextTag.Format(format), nil
}
