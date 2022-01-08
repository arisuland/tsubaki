package tsubaki

import "github.com/spf13/cobra"

func newValidateCommand() *cobra.Command {
	var validateEnv bool

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates a config.yml file or environment variables.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Long: `
The "validate" command will validate any .yml files or the environment variables to check
if your configuration is correct. This is used before launching your instance with the official
Docker image or Kubernetes Helm Chart.

To use it, simply run "tsubaki validate <config>" to validate a YAML file.
To check the environment variables, run "tsubaki validate --env"
`,
	}

	cmd.Flags().BoolVarP(&validateEnv, "env", "e", false, "")
	return cmd
}
