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
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v7"
	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/spf13/cobra"
)

var tagLongDesc = `Tag the repository with the next semantic version based on the conventional commit history of
your repository.

Environment Variables:

| Name            | Description                                                   |
|-----------------|---------------------------------------------------------------|
| NO_COLOR        | switch to using an ASCII color profile within the terminal    |
| NSV_FORMAT      | provide a go template for changing the default version format |
| NSV_SHOW        | show how the next semantic version was generated              |
| NSV_TAG_MESSAGE | a custom message for the tag, overrides the default message   |
|                 | of: chore: tagged release <version>                           |`

func tagCmd(out io.Writer) *cobra.Command {
	opts := nsv.Options{
		Err: os.Stderr,
		Out: out,
	}

	cmd := &cobra.Command{
		Use:   "tag [path]",
		Short: "Tag the repository with the next semantic version",
		Long:  tagLongDesc,
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := env.Parse(&opts); err != nil {
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

			next, err := nextVersion(gitc, opts)
			if err != nil {
				return err
			}

			if next == nil {
				return nil
			}

			printNext(next, opts)
			return tagAndPush(gitc, next, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")
	flags.StringVarP(&opts.TagMessage, "message", "m", "", "a custom message for the tag, overrides the default message of: chore: tagged release <version>")
	flags.BoolVarP(&opts.Show, "show", "s", false, "show how the next semantic version was generated")

	return cmd
}

func tagAndPush(gitc *git.Client, next *nsv.Next, opts nsv.Options) error {
	if opts.TagMessage == "" {
		opts.TagMessage = fmt.Sprintf("chore: tagged release %s", next.Tag)
	}

	cfg, err := impersonateConfig(gitc, next)
	if err != nil {
		return err
	}
	if _, err := gitc.Tag(next.Tag, git.WithTagConfig(cfg...), git.WithAnnotation(opts.TagMessage)); err != nil {
		return err
	}

	_, err = gitc.PushRef(next.Tag)
	return err
}

func impersonateConfig(gitc *git.Client, next *nsv.Next) ([]string, error) {
	// If the user.name and user.email config settings are set, then no impersonation is required
	gcfg, err := gitc.Config()
	if err != nil {
		return nil, err
	}

	_, nameSet := gcfg["user.name"]
	_, emailSet := gcfg["user.email"]
	if nameSet && emailSet {
		return nil, nil
	}

	// If git specific environment variables are set, no impersonation is required,
	// see: https://git-scm.com/book/en/v2/Git-Internals-Environment-Variables
	if os.Getenv("GIT_COMMITTER_NAME") != "" && os.Getenv("GIT_COMMITTER_EMAIL") != "" {
		return nil, nil
	}

	hash := next.Log[next.Match.Index].Hash
	commits, err := gitc.ShowCommits(hash)
	if err != nil {
		return nil, err
	}

	commit := commits[hash]
	return []string{"user.name", commit.Committer.Name, "user.email", commit.Committer.Email}, nil
}
