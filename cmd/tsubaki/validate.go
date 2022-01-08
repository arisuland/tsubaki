// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2022 Noelware
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package tsubaki

import (
	"arisu.land/tsubaki/pkg"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newValidateCommand() *cobra.Command {
	var validateEnv bool

	cmd := &cobra.Command{
		Use:          "validate [file]",
		Short:        "Validates a config.yml file or environment variables.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				if err := pkg.TestConfigFromPath(args[0]); err != nil {
					return err
				}

				_, _ = fmt.Fprintf(os.Stdout, "Config file at path %s looks fine.", args[0])
				return nil
			}

			if validateEnv {
				if _, err := pkg.LoadFromEnv(); err != nil {
					return err
				} else {
					_, _ = fmt.Fprintf(os.Stdout, "Environment variables looks alright!")
					return nil
				}
			}

			return errors.New("I need a config.yml file to validate")
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
