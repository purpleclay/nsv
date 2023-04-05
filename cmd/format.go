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
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/spf13/cobra"
)

var formatLongDesc = `A playground for discovering go template support.

Discover ways of formatting your repository tag using the in-built
go template annotations.

Environment Variables:

| Name       | Description                                       |
|------------|---------------------------------------------------|
| NSV_FORMAT | set a go template for formatting the provided tag |`

func formatCmd(out io.Writer) *cobra.Command {
	opts := nsv.Options{
		StdOut: out,
		StdErr: os.Stderr,
	}

	cmd := &cobra.Command{
		Use:   "format <tag>",
		Short: "A playground for discovering go template support",
		Long:  formatLongDesc,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return env.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := nsv.CheckTemplate(opts.VersionFormat); err != nil {
				return err
			}

			tag, err := nsv.ParseTag(args[0])
			if err != nil {
				return err
			}

			nsv.PrintFormat(out, tag, opts.VersionFormat)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.VersionFormat, "with", "", "a go template for formatting the provided tag")

	return cmd
}
