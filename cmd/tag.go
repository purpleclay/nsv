package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/nsv/internal/ci"
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
	out.WriteString(fmt.Sprintf(`template %q `, e.Template))
	out.WriteString("contains a syntax error at line " + parts[2])
	if len(parts) == 5 {
		out.WriteString(" column " + parts[3])
	}
	out.WriteString(": ")

	if strings.Contains(e.Err, "can't evaluate field") {
		si := strings.Index(e.Err, "<")
		ei := strings.Index(e.Err, ">")
		out.WriteString(fmt.Sprintf(`unrecognised field %q in template`, e.Err[si+1:ei]))
	} else {
		out.WriteString(strings.TrimSpace(parts[len(parts)-1]))
	}
	return out.String()
}

type release struct {
	Tag             string
	PrevTag         string
	SkipPipelineTag string
}

var (
	tagMessageTmpl       = "chore: tagged release {{.Tag}}"
	tagCommitMessageTmpl = "chore: patched files for release {{.Tag}} {{.SkipPipelineTag}}"
	tagTmpl              *template.Template

	tagLongDesc = `Tag the repository with the next semantic version based on the conventional commit history of
your repository.

Environment Variables:

| Name               | Description                                                    |
|--------------------|----------------------------------------------------------------|
| LOG_LEVEL          | the level of logging when printing to stderr (default: info)   |
| NO_COLOR           | switch to using an ASCII color profile within the terminal     |
| NO_LOG             | disable all log output                                         |
| NSV_COMMIT_MESSAGE | a custom message when committing file changes, supports go     |
|                    | text templates. The default is: "chore: patched files for      |
|                    | release {{.Tag}} {{.SkipPipelineTag}}"                         |
| NSV_DRY_RUN        | no changes will be made to the repository                      |
| NSV_FIX_SHALLOW    | fix a shallow clone of a repository if detected                |
| NSV_FORMAT         | provide a go template for changing the default version format  |
| NSV_HOOK           | a user-defined hook that will be executed before the           |
|                    | repository is tagged with the next semantic version            |
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
| NSV_TAG_MESSAGE    | a custom message for the annotated tag, supports go text       |
|                    | templates. The default is: "chore: tagged release {{.Tag}}"    |

Hook Environment Variables:

| Name                  | Description                                                 |
|-----------------------|-------------------------------------------------------------|
| NSV_NEXT_TAG          | the next calculated semantic version                        |
| NSV_PREV_TAG          | the last semantic version as identified within the tag      |
|                       | history of the current repository                           |
| NSV_WORKING_DIRECTORY | the working directory (or path) relative to the root of the |
|                       | current repository. It will be empty if not a monorepo      |`
)

func tagCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag [<path>...]",
		Short: "Tag the repository with the next semantic version",
		Long:  tagLongDesc,
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Paths = defaultIfEmpty(args, []string{git.RelativeAtRoot})

			for _, templatedText := range []string{opts.TagMessage, opts.CommitMessage} {
				if err := verifyTextTemplate(templatedText); err != nil {
					return err
				}
			}

			tagTmpl, _ = template.New("tag-template").Parse(opts.TagMessage)
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

			return doTag(gitc, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.CommitMessage, "commit-message", "M", tagCommitMessageTmpl, "a custom message when committing file "+
		"changes, supports go text templates")
	flags.BoolVar(&opts.DryRun, "dry-run", false, "no changes will be made to the repository")
	flags.BoolVar(&opts.FixShallow, "fix-shallow", false, "fix a shallow clone of a repository if detected")
	flags.StringVar(&opts.Hook, "hook", "", "a user-defined hook that will be executed before the repository is tagged "+
		"with the next semantic version")
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")
	flags.StringVarP(&opts.TagMessage, "tag-message", "A", tagMessageTmpl, "a custom message for the annotated tag, supports go text templates")
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

func verifyTextTemplate(tmpl string) error {
	t, err := template.New("verify-template").Parse(tmpl)
	if err != nil {
		return TemplateSyntaxError{Template: tmpl, Err: err.Error()}
	}

	rel := release{Tag: "0.2.0", PrevTag: "0.1.0", SkipPipelineTag: "[skip ci]"}

	var buf bytes.Buffer
	if err := t.Execute(&buf, rel); err != nil {
		return TemplateSyntaxError{Template: tmpl, Err: err.Error()}
	}

	return nil
}

func doTag(gitc *git.Client, opts *Options) error {
	impersonate, err := requiresImpersonation(gitc)
	if err != nil {
		return err
	}

	var tags []string
	var vers []*nsv.Next
	for _, path := range opts.Paths {
		next, err := nsv.NextVersion(gitc, nsv.Options{
			FixShallow:    opts.FixShallow,
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

		if err := commitAndTag(gitc, next, impersonate, opts); err != nil {
			return err
		}

		vers = append(vers, next)
		tags = append(tags, next.Tag)
	}

	if len(vers) == 0 {
		opts.Logger.Info("nothing to release for given paths", "paths", opts.Paths)
		return nil
	}

	if err := pushAll(gitc, tags, opts); err != nil {
		return err
	}

	printNext(vers, opts)
	return nil
}

func commitAndTag(gitc *git.Client, ver *nsv.Next, impersonate bool, opts *Options) error {
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
		opts.Logger.Info("any tag will be annotated as", "user", cfg[0], "email", cfg[1])
	}

	rel := release{
		Tag:             ver.Tag,
		PrevTag:         ver.PrevTag,
		SkipPipelineTag: ci.Detect().SkipPipelineTag,
	}

	hash, err := stageAndCommit(gitc, cfg, ver.Diffs, rel, opts)
	if err != nil {
		return err
	}

	if hash == "" {
		hash = ver.Log[0].Hash
	}

	opts.Logger.Debug("inputs to annotated tag template", "tag", rel.Tag, "prev_tag", rel.PrevTag, "skip_ci", rel.SkipPipelineTag)
	var buf bytes.Buffer
	tagTmpl.Execute(&buf, rel)

	if _, err := gitc.Tag(ver.Tag,
		git.WithTagConfig(cfg...),
		git.WithCommitRef(hash),
		git.WithLocalOnly(),
		git.WithAnnotation(buf.String())); err != nil {
		return err
	}

	opts.Logger.Info("tagged release with", "annotation", buf.String(), "hash", hash)
	return nil
}

func pushAll(gitc *git.Client, tags []string, opts *Options) error {
	if opts.DryRun {
		return nil
	}

	// NOTE: this won't work when if the repository has a detached HEAD
	branch, err := gitc.Exec("git branch --show-current")
	if err != nil {
		return err
	}

	notPushed, err := gitc.Exec(fmt.Sprintf("git log %s --not --remotes", branch))
	if err != nil {
		return err
	}

	if notPushed == "" && len(tags) == 0 {
		opts.Logger.Debug("nothing to push to remote")
		return nil
	}

	var refs []string
	refs = append(refs, branch)
	refs = append(refs, tags...)

	_, err = gitc.Push(git.WithRefSpecs(refs...))
	opts.Logger.Debug("pushed all changes to remote", "ref_specs", refs)
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
