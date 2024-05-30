package cmd

import (
	"github.com/caarlos0/env/v9"
	"github.com/purpleclay/nsv/internal/nsv"
	"github.com/purpleclay/nsv/internal/tui"
	"github.com/spf13/cobra"
)

var playgroundLongDesc = `A playground for discovering go template support.

Discover ways of formatting your repository tag using the in-built
go template annotations.

Environment Variables:

| Name       | Description                                       |
|------------|---------------------------------------------------|
| NSV_FORMAT | set a go template for formatting the provided tag |`

func playgroundCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "playground <tag>",
		Short: "A playground for discovering go template support",
		Long:  playgroundLongDesc,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return env.Parse(opts)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if err := nsv.CheckTemplate(opts.VersionFormat); err != nil {
				return err
			}

			tag, err := nsv.ParseTag(args[0])
			if err != nil {
				return err
			}

			tui.PrintFormat(tag, tui.PlaygroundOptions{
				Out:           opts.Err,
				VersionFormat: opts.VersionFormat,
			})
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.VersionFormat, "format", "f", "", "provide a go template for changing the default version format")

	return cmd
}
