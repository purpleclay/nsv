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

	"github.com/spf13/cobra"
)

type BuildDetails struct {
	Version   string `json:"version,omitempty"`
	GitBranch string `json:"git_branch,omitempty"`
	GitCommit string `json:"git_commit,omitempty"`
	Date      string `json:"build_date,omitempty"`
}

func Execute(out io.Writer, buildInfo BuildDetails) error {
	cmd := &cobra.Command{
		Use:           "nsv",
		Short:         "Manage your semantic versioning without any config",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(versionCmd(out, buildInfo),
		manCmd(out),
		playgroundCmd(out),
		nextCmd(out),
		tagCmd())

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
