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
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/caarlos0/env/v7"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/spf13/cobra"
)

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

| Name     | Description                                                   |
|----------|---------------------------------------------------------------|
| NO_COLOR | switch to using an ASCII color profile within the terminal    |`

type BuildDetails struct {
	Version   string `json:"version,omitempty"`
	GitBranch string `json:"git_branch,omitempty"`
	GitCommit string `json:"git_commit,omitempty"`
	Date      string `json:"build_date,omitempty"`
}

func Execute(out io.Writer, buildInfo BuildDetails) error {
	opts := nsv.Options{}

	cmd := &cobra.Command{
		Use:           "nsv",
		Short:         "Manage your semantic versioning without any config",
		Long:          rootLongDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := env.Parse(&opts); err != nil {
				return err
			}

			if opts.NoColor {
				lipgloss.SetColorProfile(termenv.Ascii)
			}
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&opts.NoColor, "no-color", false, "switch to using an ASCII color profile within the terminal")

	cmd.AddCommand(versionCmd(out, buildInfo),
		manCmd(out),
		playgroundCmd(out),
		nextCmd(out),
		tagCmd(out))

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
