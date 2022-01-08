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
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// InputFlags is a object that represents all the global flags
// for all commands.
type InputFlags struct {
	// Verbose returns a bool if we should turn on debug mode. Default is `false`.
	// This is override in `GO_ENV` environment variable.
	Verbose bool

	// ConfigFile returns the configuration file, in the Docker container,
	// this is usually `/usr/share/arisu/tsubaki/config.yml`. Default is `./config.yml`
	ConfigFile *string
}

var rootCmd = &cobra.Command{
	Use:               "tsubaki",
	Short:             "Tsubaki is the core heart of Arisu.",
	RunE:              run,
	PersistentPreRunE: preRun,
	Long: `
Tsubaki is the core backend for Arisu, this is what powers all services like:
	- Fubuki (frontend)
	- GitHub Bot (automative microservice for GitHub <- -> Arisu)

To run the server, you can run "tsubaki -c ./config.yml" to start up the server.
Want to validate your configuration? Run "tsubaki validate" to check if it is OK.

You can get the current version of Tsubaki by running localhost:<port>/version and retrieving
the "version" field or you can run "tsubaki version" to retrieve the current version.
`,
}

var GlobalFlags = new(InputFlags)

func Execute() int {
	// setup flags
	rootCmd.Flags().BoolVarP(&GlobalFlags.Verbose, "verbose", "d", false, "")
	rootCmd.Flags().StringVarP(GlobalFlags.ConfigFile, "config-file", "c", "./config.yml", "")
	rootCmd.AddCommand(
		newVersionCommand(),
		newValidateCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}

func run(_ *cobra.Command, args []string) error {
	// run the server if there is no args
	if len(args) == 0 {
		return nil
	}

	return fmt.Errorf("unknown tsubaki command %q", "tsubaki "+args[0])
}

func preRun(cmd *cobra.Command, args []string) error {
	parent := cmd.Root()
	if parent != nil {
		prerun := parent.PersistentPreRunE
		if prerun != nil {
			err := prerun(cmd, args)
			if err != nil {
				return err
			}
		}
	}

	if GlobalFlags.Verbose {
		logrus.SetLevel(logrus.TraceLevel)
	}

	return nil
}
