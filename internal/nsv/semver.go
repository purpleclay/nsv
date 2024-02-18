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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	git "github.com/purpleclay/gitz"
)

const (
	vPrefix          = 'v'
	firstVer         = "0.0.0"
	prefixedFirstVer = "v0.0.0"
)

type Increment int

const (
	NoIncrement Increment = iota
	PatchIncrement
	MinorIncrement
	MajorIncrement
)

type Options struct {
	Path          string
	VersionFormat string
}

type context struct {
	TagPrefix string
	LogPath   string
}

func resolveContext(gitc *git.Client, opts Options) (*context, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var relPath string
	if opts.Path != "" {
		relPath = opts.Path
	} else {
		if relPath, err = gitc.ToRelativePath(cwd); err != nil {
			return nil, err
		}
	}

	if relPath == git.RelativeAtRoot {
		return &context{TagPrefix: "", LogPath: ""}, nil
	}

	logPath := relPath
	if strings.HasSuffix(cwd, logPath) {
		logPath = git.RelativeAtRoot
	}

	return &context{
		TagPrefix: filepath.Base(relPath),
		LogPath:   logPath,
	}, nil
}

type Tag struct {
	Prefix   string
	SemVer   string
	Version  string
	Raw      string
	Pre      string
	Metadata string
}

func ParseTag(raw string) (Tag, error) {
	lastSlash := 0
	if idx := strings.LastIndex(raw, "/"); idx > -1 {
		lastSlash = idx + 1
	}

	prefix := ""
	if lastSlash > 0 {
		prefix = raw[:lastSlash-1]
	}

	semv := raw[lastSlash:]
	if semv[0] == vPrefix {
		semv = semv[1:]
	}

	sv, err := semver.StrictNewVersion(semv)
	if err != nil {
		return Tag{}, err
	}

	return Tag{
		Prefix:   prefix,
		Raw:      raw,
		SemVer:   semv,
		Version:  raw[lastSlash:],
		Pre:      sv.Prerelease(),
		Metadata: sv.Metadata(),
	}, nil
}

func (t Tag) Bump(semv string) Tag {
	ver := semv
	if t.Version[0] == vPrefix {
		ver = fmt.Sprintf("%c%s", vPrefix, ver)
	}

	raw := ver
	if t.Prefix != "" {
		raw = fmt.Sprintf("%s/%s", t.Prefix, raw)
	}

	tag, _ := ParseTag(raw)
	return tag
}

func (t Tag) Format(format string) string {
	var tagf bytes.Buffer
	if format == "" {
		_ = versionTmpl.Execute(&tagf, t)
		return tagf.String()
	}

	tmpl, _ := template.New("custom-format").Parse(format)
	_ = tmpl.Execute(&tagf, t)
	return tagf.String()
}

func (t Tag) Prerelease() bool {
	return len(t.Pre) > 0
}

func (t Tag) PrereleaseWithLabel(label string) bool {
	return strings.HasPrefix(t.Pre, label)
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
	ctx, err := resolveContext(gitc, opts)
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
		ltag = firstVersion(ctx)
	}
	ver, _ := ParseTag(ltag)

	if cmd.Prerelease != "" && !ver.PrereleaseWithLabel(cmd.Prerelease) {
		// To prevent any conflict with prerelease tags, query git for the latest tag based
		// on the prerelease label. Patch existing tag as needed
		if preTag, _ := latestPrereleaseTag(gitc, ctx.TagPrefix, cmd.Prerelease); preTag != "" {
			ver, _ = ParseTag(preTag)
		}
	}

	nextVer, err := bump(ver, opts.VersionFormat, inc, cmd)
	if err != nil {
		return nil, err
	}

	return &Next{
		PrevTag: ltag,
		Tag:     nextVer,
		Log:     log.Commits,
		LogDir:  ctx.LogPath,
		Match:   match,
	}, nil
}

func latestTag(gitc *git.Client, prefix string) (string, error) {
	return latestTagByGlob(gitc, prefix, "**/*.*.*")
}

func latestTagByGlob(gitc *git.Client, prefix, glob string) (string, error) {
	prefixFilter := func(tag string) bool {
		if prefix == "" {
			return true
		}

		return strings.HasPrefix(tag, prefix+"/")
	}

	tags, err := gitc.Tags(git.WithShellGlob(glob),
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

func latestPrereleaseTag(gitc *git.Client, prefix, label string) (string, error) {
	return latestTagByGlob(gitc, prefix, fmt.Sprintf("**/*.*.*-%s*", label))
}

func firstVersion(ctx *context) string {
	fv := firstVer
	if detectLanguageIsPrefixed(ctx.LogPath) {
		fv = prefixedFirstVer
	}

	if ctx.TagPrefix == "" {
		return fv
	}

	return fmt.Sprintf("%s/%s", ctx.TagPrefix, fv)
}

func bump(ver Tag, format string, inc Increment, cmd Command) (string, error) {
	semv, err := semver.StrictNewVersion(ver.SemVer)
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

	if ver.Prerelease() {
		// Prerelease versions have a lower precedence than a normal version, so invoke
		// a patch to ensure (0.1.0-beta.1) is bumped to (0.1.0), http://semver.org/#spec-item-9
		inc = PatchIncrement
	}

	switch inc {
	case MajorIncrement:
		bumpedVer = semv.IncMajor()
	case MinorIncrement:
		bumpedVer = semv.IncMinor()
	case PatchIncrement:
		bumpedVer = semv.IncPatch()
	}

	if cmd.Prerelease != "" {
		if ver.Prerelease() {
			// As this is an existing prerelease, increment it as expected
			label, inc, _ := strings.Cut(ver.Pre, ".")
			i, _ := strconv.Atoi(inc)
			bumpedVer, _ = bumpedVer.SetPrerelease(fmt.Sprintf("%s.%d", label, i+1))
		} else {
			bumpedVer, _ = bumpedVer.SetPrerelease(cmd.Prerelease + ".1")
		}
	}

	nextTag := ver.Bump(bumpedVer.String())
	return nextTag.Format(format), nil
}
