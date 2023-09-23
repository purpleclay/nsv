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

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/purpleclay/nsv/internal/tui"
	"github.com/spf13/cobra"
)

type TemplateSyntaxError struct {
	Template string
	Err      string
}

func (e TemplateSyntaxError) Error() string {
	parts := strings.Split(e.Err, ":")

	var out strings.Builder
	out.WriteString(fmt.Sprintf(`template "%s" `, e.Template))
	out.WriteString("contains a syntax error at line " + parts[2])
	if len(parts) == 5 {
		out.WriteString(" column " + parts[3])
	}
	out.WriteString(": ")

	if strings.Contains(e.Err, "can't evaluate field") {
		si := strings.Index(e.Err, "<")
		ei := strings.Index(e.Err, ">")
		out.WriteString(fmt.Sprintf(`unrecognised field "%s" in template`, e.Err[si+1:ei]))
	} else {
		out.WriteString(strings.TrimSpace(parts[len(parts)-1]))
	}
	return out.String()
}

type release struct {
	Tag     string
	PrevTag string
}

var tagLongDesc = `Tag the repository with the next semantic version based on the conventional commit history of
your repository.

Environment Variables:

| Name            | Description                                                    |
|-----------------|----------------------------------------------------------------|
| NO_COLOR        | switch to using an ASCII color profile within the terminal     |
| NSV_FORMAT      | provide a go template for changing the default version format  |
| NSV_PRETTY      | pretty-print the output of the next semantic version in a      |
|                 | given format. The format can be one of either full or compact. |
|                 | full is the default. Must be used in conjunction with NSV_SHOW |
| NSV_SHOW        | show how the next semantic version was generated               |
| NSV_TAG_MESSAGE | a custom message for the tag, supports go text templates. The  |
|                 | default is: "chore: tagged release {{.Tag}}"                   |`

func tagCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag [<path>...]",
		Short: "Tag the repository with the next semantic version",
		Long:  tagLongDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := supportedPretty(opts.Pretty); err != nil {
				return err
			}

			if err := verifyTagTemplate(opts.TagMessage); err != nil {
				return err
			}

			opts.Paths = args
			return pathsExist(opts.Paths)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			gitc, err := git.NewClient()
			if err != nil {
				return err
			}

			vers, err := nextVersions(gitc, opts)
			if err != nil {
				return err
			}

			if len(vers) == 0 {
				return nil
			}

			printNext(vers, opts)
			return tagAndPush(gitc, vers, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")
	flags.StringVarP(&opts.TagMessage, "message", "m", "chore: tagged release {{.Tag}}", "a custom message for the tag, supports go text templates")
	flags.StringVarP(&opts.Pretty, "pretty", "p", string(tui.Full), "pretty-print the output of the next semantic version in a given format. "+
		"The format can be one of either full or compact. Must be used in conjunction with --show")
	flags.BoolVarP(&opts.Show, "show", "s", false, "show how the next semantic version was generated")

	cmd.RegisterFlagCompletionFunc("pretty", prettyFlagShellComp)
	return cmd
}

func verifyTagTemplate(tmpl string) error {
	t, err := template.New("tag-template").Parse(tmpl)
	if err != nil {
		return TemplateSyntaxError{Template: tmpl, Err: err.Error()}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, release{Tag: "0.2.0", PrevTag: "0.1.0"}); err != nil {
		return TemplateSyntaxError{Template: tmpl, Err: err.Error()}
	}

	return nil
}

func tagAndPush(gitc *git.Client, vers []*nsv.Next, opts *Options) error {
	impersonate, err := requiresImpersonation(gitc)
	if err != nil {
		return err
	}

	tmpl, _ := template.New("tag-template").Parse(opts.TagMessage)

	var tags []string
	for _, ver := range vers {
		var cfg []string
		var err error
		if impersonate {
			if cfg, err = impersonateConfig(gitc, ver); err != nil {
				return err
			}
		}

		var buf bytes.Buffer
		tmpl.Execute(&buf, release{Tag: ver.Tag, PrevTag: ver.PrevTag})

		if _, err := gitc.Tag(ver.Tag,
			git.WithTagConfig(cfg...),
			git.WithCommitRef(ver.Log[0].Hash),
			git.WithLocalOnly(),
			git.WithAnnotation(buf.String())); err != nil {
			return err
		}
		tags = append(tags, ver.Tag)
	}

	_, err = gitc.Push(git.WithRefSpecs(tags...))
	return err
}

func requiresImpersonation(gitc *git.Client) (bool, error) {
	// If the user.name and user.email config settings are set, then no impersonation is required
	gcfg, err := gitc.Config()
	if err != nil {
		return false, err
	}

	_, nameSet := gcfg["user.name"]
	_, emailSet := gcfg["user.email"]
	if nameSet && emailSet {
		return false, nil
	}

	// If git specific environment variables are set, no impersonation is required,
	// see: https://git-scm.com/book/en/v2/Git-Internals-Environment-Variables
	if os.Getenv("GIT_COMMITTER_NAME") != "" && os.Getenv("GIT_COMMITTER_EMAIL") != "" {
		return false, nil
	}

	return true, nil
}

func impersonateConfig(gitc *git.Client, next *nsv.Next) ([]string, error) {
	hash := next.Log[next.Match.Index].Hash
	commits, err := gitc.ShowCommits(hash)
	if err != nil {
		return nil, err
	}

	commit := commits[hash]
	return []string{"user.name", commit.Committer.Name, "user.email", commit.Committer.Email}, nil
}
