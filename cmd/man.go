package cmd

import (
	"fmt"
	"io"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

func manCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generate man pages",
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mp, err := mcobra.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(out, mp.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}
