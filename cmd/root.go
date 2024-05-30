package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/caarlos0/env/v9"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	theme "github.com/purpleclay/lipgloss-theme"
	"github.com/spf13/cobra"
)

// TODO: should be set environment defaults that align to CLI flags?

type Options struct {
	Err           io.Writer   `env:"-"`
	Logger        *log.Logger `env:"_"`
	LogLevel      string      `env:"LOG_LEVEL"`
	NoColor       bool        `env:"NO_COLOR"`
	Out           io.Writer   `env:"-"`
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

| Name      | Description                                                    |
|-----------|----------------------------------------------------------------|
| LOG_LEVEL | the level of logging when outputting to stderr (default: info) |
| NO_COLOR  | switch to using an ASCII color profile within the terminal     |`

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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := env.Parse(opts); err != nil {
				return err
			}

			if opts.NoColor {
				lipgloss.SetColorProfile(termenv.Ascii)
			}

			// TODO: support io.Discard as a writer option, if logging is to be turned off
			logLevel, _ := log.ParseLevel(opts.LogLevel)

			opts.Logger = log.NewWithOptions(os.Stderr, log.Options{
				Level:           logLevel,
				ReportCaller:    false,
				ReportTimestamp: false,
			})
			opts.Logger.SetStyles(theme.Logging)
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&opts.LogLevel, "log-level", "info", "the level of logging when outputting to stderr")
	flags.BoolVar(&opts.NoColor, "no-color", false, "switch to using an ASCII color profile within the terminal")

	cmd.AddCommand(versionCmd(out, buildInfo),
		manCmd(out),
		playgroundCmd(opts),
		nextCmd(opts),
		tagCmd(opts),
	)

	cmd.SetUsageTemplate(customUsageTemplate)
	return cmd.Execute()
}

func versionCmd(out io.Writer, buildInfo BuildDetails) *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print build time version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if short {
				fmt.Fprintf(out, buildInfo.Version)
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
