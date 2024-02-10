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

package cmd

import (
	"fmt"
	"os"
	"strings"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/purpleclay/nsv/internal/tui"
	"github.com/spf13/cobra"
)

type MissingPathsError struct {
	Paths []string
}

func (e MissingPathsError) Error() string {
	return "paths do not exist within the current repository: " + strings.Join(e.Paths, ", ")
}

type InvalidPrettyFormatError struct {
	Pretty string
}

func (e InvalidPrettyFormatError) Error() string {
	return fmt.Sprintf("pretty format '%s' is not supported, must be one of either: %s",
		e.Pretty, strings.Join(tui.PrettyFormats, ", "))
}

var nextLongDesc = `Generate the next semantic version based on the conventional commit history of your repository.

Environment Variables:

| Name       | Description                                                    |
|------------|----------------------------------------------------------------|
| NO_COLOR   | switch to using an ASCII color profile within the terminal     |
| NSV_FORMAT | provide a go template for changing the default version format  |
| NSV_PRETTY | pretty-print the output of the next semantic version in a      |
|            | given format. The format can be one of either full or compact. |
|            | full is the default. Must be used in conjunction with NSV_SHOW |
| NSV_SHOW   | show how the next semantic version was generated               |`

func nextCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next [<path>...]",
		Short: "Generate the next semantic version",
		Long:  nextLongDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := supportedPretty(opts.Pretty); err != nil {
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
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")
	flags.StringVarP(&opts.Pretty, "pretty", "p", string(tui.Full), "pretty-print the output of the next semantic version in a given format. "+
		"The format can be one of either full or compact. Must be used in conjunction with --show")
	flags.BoolVarP(&opts.Show, "show", "s", false, "show how the next semantic version was generated")

	cmd.RegisterFlagCompletionFunc("pretty", prettyFlagShellComp)
	return cmd
}

func prettyFlagShellComp(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return tui.PrettyFormats, cobra.ShellCompDirectiveDefault
}

func supportedPretty(pretty string) error {
	for _, p := range tui.PrettyFormats {
		if p == pretty {
			return nil
		}
	}
	return InvalidPrettyFormatError{Pretty: pretty}
}

func pathsExist(paths []string) error {
	var notFound []string
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			notFound = append(notFound, path)
		}
	}

	if len(notFound) > 0 {
		return MissingPathsError{Paths: notFound}
	}

	return nil
}

func nextVersions(gitc *git.Client, opts *Options) ([]*nsv.Next, error) {
	if err := nsv.CheckTemplate(opts.VersionFormat); err != nil {
		return nil, err
	}

	if len(opts.Paths) == 0 {
		opts.Paths = append(opts.Paths, git.RelativeAtRoot)
	}

	var vers []*nsv.Next
	for _, path := range opts.Paths {
		next, err := nsv.NextVersion(gitc, nsv.Options{
			Path:          path,
			VersionFormat: opts.VersionFormat,
		})
		if err != nil {
			return nil, err
		}

		if next != nil {
			vers = append(vers, next)
		}
	}

	return vers, nil
}

func printNext(vers []*nsv.Next, opts *Options) {
	var tags []string
	for _, ver := range vers {
		tags = append(tags, ver.Tag)
	}

	fmt.Fprint(opts.Out, strings.Join(tags, ","))

	if opts.Show {
		tui.PrintSummary(vers, tui.SummaryOptions{
			NoColor: opts.NoColor,
			Out:     opts.Err,
			Pretty:  tui.Pretty(opts.Pretty),
		})
	}
}
