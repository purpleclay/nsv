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
	Format string
}

func (e InvalidPrettyFormatError) Error() string {
	return fmt.Sprintf("pretty format '%s' is not supported, must be one of either: %s",
		e.Format, strings.Join(tui.PrettyFormats, ", "))
}

var nextLongDesc = `Generate the next semantic version based on the conventional commit history of your repository.

Environment Variables:

| Name               | Description                                                    |
|--------------------|----------------------------------------------------------------|
| LOG_LEVEL          | the level of logging when printing to stderr (default: info)   |
| NO_COLOR           | switch to using an ASCII color profile within the terminal     |
| NO_LOG             | disable all log output                                         |
| NSV_FORMAT         | provide a go template for changing the default version format  |
| NSV_MAJOR_PREFIXES | a comma separated list of conventional commit prefixes for     |
|                    | triggering a major semantic version increment                  |
| NSV_MINOR_PREFIXES | a comma separated list of conventional commit prefixes for     |
|                    | triggering a minor semantic version increment                  |
| NSV_PATCH_PREFIXES | a comma separated list of conventional commit prefixes for     |
|                    | triggering a patch semantic version increment                  |
| NSV_PRETTY         | pretty-print the output of the next semantic version in a      |
|                    | given format. The format can be one of either full or compact. |
|                    | Must be used in conjunction with NSV_SHOW (default: full)      |
| NSV_SHOW           | show how the next semantic version was generated               |`

func nextCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next [<path>...]",
		Short: "Generate the next semantic version",
		Long:  nextLongDesc,
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Paths = defaultIfEmpty(args, []string{git.RelativeAtRoot})

			return preRunChecks(opts)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			gitc, err := git.NewClient()
			if err != nil {
				return err
			}

			return doNext(gitc, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")
	flags.StringSliceVar(&opts.MajorPrefixes, "major-prefixes", []string{}, "a comma separated list of conventional commit prefixes for "+
		"triggering a major semantic version increment")
	flags.StringSliceVar(&opts.MinorPrefixes, "minor-prefixes", []string{}, "a comma separated list of conventional commit prefixes for "+
		"triggering a minor semantic version increment")
	flags.StringSliceVar(&opts.PatchPrefixes, "patch-prefixes", []string{}, "a comma separated list of conventional commit prefixes for "+
		"triggering a patch semantic version increment")
	flags.StringVarP(&opts.Pretty, "pretty", "p", string(tui.Full), "pretty-print the output of the next semantic version in a given format. "+
		"The format can be one of either full or compact. Must be used in conjunction with --show")
	flags.BoolVarP(&opts.Show, "show", "s", false, "show how the next semantic version was generated")

	cmd.RegisterFlagCompletionFunc("pretty", prettyFlagShellComp)
	return cmd
}

func prettyFlagShellComp(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return tui.PrettyFormats, cobra.ShellCompDirectiveDefault
}

func preRunChecks(opts *Options) error {
	if err := supportedPrettyFormat(opts.Pretty); err != nil {
		return err
	}

	if err := nsv.CheckTemplate(opts.VersionFormat); err != nil {
		return err
	}

	return pathsExist(opts.Paths)
}

func supportedPrettyFormat(format string) error {
	for _, p := range tui.PrettyFormats {
		if p == format {
			return nil
		}
	}
	return InvalidPrettyFormatError{Format: format}
}

func defaultIfEmpty(paths, def []string) []string {
	if len(paths) == 0 {
		return def
	}

	return paths
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

func doNext(gitc *git.Client, opts *Options) error {
	var vers []*nsv.Next
	for _, path := range opts.Paths {
		next, err := nsv.NextVersion(gitc, nsv.Options{
			MajorPrefixes: opts.MajorPrefixes,
			MinorPrefixes: opts.MinorPrefixes,
			Logger:        opts.Logger,
			PatchPrefixes: opts.PatchPrefixes,
			Path:          path,
			VersionFormat: opts.VersionFormat,
		})
		if err != nil {
			return err
		}

		if next != nil {
			vers = append(vers, next)
		}
	}

	if len(vers) == 0 {
		return nil
	}

	printNext(vers, opts)
	return nil
}

func printNext(vers []*nsv.Next, opts *Options) {
	var tags []string
	for _, ver := range vers {
		tags = append(tags, ver.Tag)
	}

	if !opts.NoLog {
		fmt.Fprintln(opts.Err)
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
