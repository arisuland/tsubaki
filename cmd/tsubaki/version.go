package tsubaki

import "github.com/spf13/cobra"

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Short: "Returns the current version of Tsubaki.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
