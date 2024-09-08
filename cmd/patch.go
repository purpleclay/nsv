package cmd

import (
	"bytes"
	"text/template"

	"github.com/purpleclay/chomp"
	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/ci"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/purpleclay/nsv/internal/tui"
	"github.com/spf13/cobra"
)

var (
	commitMessageTmpl = "chore: patched files for release {{.Tag}}"
	commitTmpl        *template.Template

	patchLongDesc = `Patch files in a repository with the next semantic version based on the conventional commit
history of your repository.

Environment Variables:

| Name               | Description                                                    |
|--------------------|----------------------------------------------------------------|
| LOG_LEVEL          | the level of logging when printing to stderr (default: info)   |
| NO_COLOR           | switch to using an ASCII color profile within the terminal     |
| NO_LOG             | disable all log output                                         |
| NSV_COMMIT_MESSAGE | a custom message when committing file changes, supports go     |
|                    | text templates. The default is: "chore: patched files for      |
|                    | release {{.Tag}}"                                              |
| NSV_DRY_RUN        | no changes will be made to the repository                      |
| NSV_FORMAT         | provide a go template for changing the default version format  |
| NSV_HOOK           | a user-defined hook that will be executed before any file      |
|                    | changes are committed with the next semantic version           |
| NSV_MAJOR_PREFIXES | a comma separated list of conventional commit prefixes for     |
|                    | triggering a major semantic version increment                  |
| NSV_MINOR_PREFIXES | a comma separated list of conventional commit prefixes for     |
|                    | triggering a minor semantic version increment                  |
| NSV_PATCH_PREFIXES | a comma separated list of conventional commit prefixes for     |
|                    | triggering a patch semantic version increment                  |
| NSV_PRETTY         | pretty-print the output of the next semantic version in a      |
|                    | given format. The format can be one of either full or compact. |
|                    | Must be used in conjunction with NSV_SHOW (default: full)      |
| NSV_SHOW           | show how the next semantic version was generated               |

Hook Environment Variables:

| Name                  | Description                                                 |
|-----------------------|-------------------------------------------------------------|
| NSV_NEXT_TAG          | the next calculated semantic version                        |
| NSV_PREV_TAG          | the last semantic version as identified within the tag      |
|                       | history of the current repository                           |
| NSV_WORKING_DIRECTORY | the working directory (or path) relative to the root of the |
|                       | current repository. It will be empty if not a monorepo      |`
)

func patchCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch [<path>...]",
		Short: "Patch files within a repository with the next semantic version",
		Long:  patchLongDesc,
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Paths = defaultIfEmpty(args, []string{git.RelativeAtRoot})

			if err := verifyTextTemplate(opts.CommitMessage); err != nil {
				return err
			}
			commitTmpl, _ = template.New("commit-template").Parse(opts.CommitMessage)

			return preRunChecks(opts)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			gitc, err := git.NewClient()
			if err != nil {
				return err
			}

			if opts.DryRun {
				opts.Logger.Warn("no changes will be made in dry run mode")
			}

			return doPatch(gitc, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.CommitMessage, "commit-message", "M", commitMessageTmpl, "a custom message when committing file "+
		"changes, supports go text templates")
	flags.BoolVar(&opts.DryRun, "dry-run", false, "no changes will be made to the repository")
	flags.StringVar(&opts.Hook, "hook", "", "a user-defined hook that will be executed before any file changes are committed "+
		"with the next semantic version")
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

	return cmd
}

func doPatch(gitc *git.Client, opts *Options) error {
	impersonate, err := requiresImpersonation(gitc)
	if err != nil {
		return err
	}

	var vers []*nsv.Next
	for _, path := range opts.Paths {
		next, err := nsv.NextVersion(gitc, nsv.Options{
			Hook:          opts.Hook,
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

		if next == nil {
			continue
		}

		if err := commitChanges(gitc, next, impersonate, opts); err != nil {
			return err
		}

		vers = append(vers, next)
	}

	if len(vers) == 0 {
		return nil
	}

	var noTags []string
	if err := pushAll(gitc, noTags, opts); err != nil {
		return err
	}

	printNext(vers, opts)
	return nil
}

func commitChanges(gitc *git.Client, ver *nsv.Next, impersonate bool, opts *Options) error {
	if opts.DryRun {
		statuses, err := gitc.PorcelainStatus()
		if err != nil {
			return err
		}
		return gitc.RestoreUsing(statuses)
	}

	var cfg []string
	var err error
	if impersonate {
		if cfg, err = impersonateConfig(gitc, ver); err != nil {
			return err
		}
	}

	rel := release{
		Tag:             ver.Tag,
		PrevTag:         ver.PrevTag,
		SkipPipelineTag: ci.Detect().SkipPipelineTag,
	}

	_, err = stageAndCommit(gitc, cfg, ver.Diffs, rel)
	return err
}

func stageAndCommit(gitc *git.Client, cfg []string, changes []git.FileDiff, rel release) (string, error) {
	if len(changes) == 0 {
		return "", nil
	}

	var paths []string
	for _, change := range changes {
		paths = append(paths, change.Path)
	}

	if _, err := gitc.Stage(git.WithPathSpecs(paths...)); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	commitTmpl.Execute(&buf, rel)

	msg, err := gitc.Commit(buf.String(), git.WithCommitConfig(cfg...))
	if err != nil {
		return "", err
	}

	_, marker, err := chomp.BracketSquare()(msg)
	if err != nil {
		return "", err
	}

	_, ext, err := chomp.SepPair(chomp.Until(" "), chomp.Tag(" "), chomp.Eol())(marker)
	if err != nil {
		return "", err
	}
	return ext[1], nil
}
