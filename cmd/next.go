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
	"io"
	"os"

	"github.com/caarlos0/env/v7"
	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/spf13/cobra"
)

var nextLongDesc = `Generate the next semantic version based on the conventional commit history
of your repository.

Environment Variables:

| Name       | Description                                                   |
|------------|---------------------------------------------------------------|
| NSV_FORMAT | provide a go template for changing the default version format |
| NSV_SHOW   | show how the next semantic version was generated              |`

func nextCmd(out io.Writer) *cobra.Command {
	opts := nsv.Options{
		StdOut: out,
		StdErr: os.Stderr,
	}

	cmd := &cobra.Command{
		Use:   "next",
		Short: "Generate the next semantic version",
		Long:  nextLongDesc,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return env.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			gitc, err := git.NewClient()
			if err != nil {
				return err
			}

			return nextVersion(gitc, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Show, "show", "s", false, "show how the next semantic version was generated")
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")

	return cmd
}

func nextVersion(gitc *git.Client, opts nsv.Options) error {
	if err := nsv.CheckTemplate(opts.VersionFormat); err != nil {
		return err
	}

	return nsv.NextVersion(gitc, opts)
}
