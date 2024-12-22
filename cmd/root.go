package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	theme "github.com/purpleclay/lipgloss-theme"
	"github.com/spf13/cobra"
)

var logLevels = []string{"debug", "info", "warn", "error", "fatal"}

type Options struct {
	CommitMessage string      `env:"NSV_COMMIT_MESSAGE"`
	DryRun        bool        `env:"NSV_DRY_RUN"`
	Err           io.Writer   `env:"-"`
	FixShallow    bool        `env:"NSV_FIX_SHALLOW"`
	Hook          string      `env:"NSV_HOOK"`
	Logger        *log.Logger `env:"-"`
	LogLevel      string      `env:"LOG_LEVEL"`
	MajorPrefixes []string    `env:"NSV_MAJOR_PREFIXES"`
	MinorPrefixes []string    `env:"NSV_MINOR_PREFIXES"`
	NoColor       bool        `env:"NO_COLOR"`
	NoLog         bool        `env:"NO_LOG"`
	Out           io.Writer   `env:"-"`
	PatchPrefixes []string    `env:"NSV_PATCH_PREFIXES"`
	Paths         []string    `env:"-"`
	Pretty        string      `env:"NSV_PRETTY"`
	Show          bool        `env:"NSV_SHOW"`
	TagMessage    string      `env:"NSV_TAG_MESSAGE"`
	VersionFormat string      `env:"NSV_FORMAT"`
}

var rootLongDesc = `NSV (Next Semantic Version) is a convention-based semantic versioning tool that
leans on the power of conventional commits to make versioning your software a breeze!.

## Why another versioning tool

There are many semantic versioning tools already out there! But they typically require some
configuration or custom scripting in your CI system to make them work. No one likes managing
config; it is error-prone, and the slightest tweak ultimately triggers a cascade of change
across your projects.

Step in NSV. Designed to make intelligent semantic versioning decisions about your project
without needing a config file. Entirely convention-based, you can adapt your workflow from
within your commit message.

The power is at your fingertips.

Global Environment Variables:

| Name      | Description                                                  |
|-----------|--------------------------------------------------------------|
| LOG_LEVEL | the level of logging when printing to stderr (default: info) |
| NO_COLOR  | switch to using an ASCII color profile within the terminal   |
| NO_LOG    | disable all log output                                       |`

type BuildDetails struct {
	Version   string `json:"version,omitempty"`
	GitBranch string `json:"git_branch,omitempty"`
	GitCommit string `json:"git_commit,omitempty"`
	Date      string `json:"build_date,omitempty"`
}

func Execute(out io.Writer, buildInfo BuildDetails) error {
	opts := &Options{
		Err: os.Stderr,
		Out: out,
	}

	cmd := &cobra.Command{
		Use:           "nsv",
		Short:         "Manage your semantic versioning without any config",
		Long:          rootLongDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := env.Parse(opts); err != nil {
				return err
			}

			if opts.NoColor {
				lipgloss.SetColorProfile(termenv.Ascii)
			}

			var logw io.Writer

			logw = opts.Err
			if opts.NoLog {
				logw = io.Discard
			}

			logLevel, _ := log.ParseLevel(opts.LogLevel)
			opts.Logger = log.NewWithOptions(logw, log.Options{
				Level:           logLevel,
				ReportCaller:    false,
				ReportTimestamp: false,
			})
			opts.Logger.SetStyles(theme.Logging)
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&opts.LogLevel, "log-level", "info", "the level of logging when printing to stderr")
	flags.BoolVar(&opts.NoColor, "no-color", false, "switch to using an ASCII color profile within the terminal")
	flags.BoolVar(&opts.NoLog, "no-log", false, "disable all log output")

	cmd.RegisterFlagCompletionFunc("log-level", logLevelFlagShellComp)

	cmd.AddCommand(versionCmd(out, buildInfo),
		manCmd(out),
		playgroundCmd(opts),
		nextCmd(opts),
		tagCmd(opts),
		patchCmd(opts),
	)

	cmd.SetUsageTemplate(customUsageTemplate)
	return cmd.Execute()
}

func logLevelFlagShellComp(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return logLevels, cobra.ShellCompDirectiveDefault
}

func versionCmd(out io.Writer, buildInfo BuildDetails) *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print build time version information",
		RunE: func(_ *cobra.Command, _ []string) error {
			if short {
				fmt.Fprint(out, buildInfo.Version)
				return nil
			}

			ver := struct {
				Go     string `json:"go"`
				GoArch string `json:"go_arch"`
				GoOS   string `json:"go_os"`
				BuildDetails
			}{
				Go:           runtime.Version(),
				GoArch:       runtime.GOARCH,
				GoOS:         runtime.GOOS,
				BuildDetails: buildInfo,
			}
			return json.NewEncoder(out).Encode(&ver)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&short, "short", false, "only print the version number")

	return cmd
}
