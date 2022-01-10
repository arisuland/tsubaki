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
	"arisu.land/tsubaki/server"
	"arisu.land/tsubaki/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// InputFlags is a object that represents all the global flags
// for all commands.
type InputFlags struct {
	// Verbose returns a bool if we should turn on debug mode. Default is `false`.
	// This is override in `GO_ENV` environment variable.
	Verbose bool

	// ConfigFile returns the configuration file, in the Docker container,
	// this is usually `/usr/share/arisu/tsubaki/config.yml`. Default is `./config.yml`
	ConfigFile string
}

var rootCmd = &cobra.Command{
	Use:               "tsubaki [command] [...args]",
	Short:             "Tsubaki is the core heart of Arisu.",
	RunE:              run,
	PersistentPreRunE: preRun,
	SilenceUsage:      true,
	Long: `Tsubaki is the core backend for Arisu, this is what powers all services like:
	- Fubuki (frontend)
	- GitHub Bot (automative microservice for GitHub <- -> Arisu)

To run the server, you can run "tsubaki -c ./config.yml" to start up the server.
Want to validate your configuration? Run "tsubaki validate" to check if it is OK.

You can get the current version of Tsubaki by running localhost:<port>/version and retrieving
the "version" field or you can run "tsubaki version" to retrieve the current version.
`,
}

var globalFlags = new(InputFlags)

func Execute() int {
	// setup flags
	rootCmd.Flags().BoolVarP(&globalFlags.Verbose, "verbose", "v", false, "if debug logging should be enabled or not.")
	rootCmd.Flags().StringVarP(&globalFlags.ConfigFile, "config-file", "c", "./config.yml", "the configuration file to use when starting up the server")
	rootCmd.AddCommand(
		newVersionCommand(),
		newValidateCommand(),
		newGenerateCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}

func run(_ *cobra.Command, args []string) error {
	// run the server if there is no args
	if len(args) == 0 {
		util.PrintBanner()
		return server.Start(globalFlags.ConfigFile)
	}

	return fmt.Errorf("unknown tsubaki command %q", "tsubaki "+args[0])
}

func preRun(_ *cobra.Command, _ []string) error {
	if globalFlags.Verbose {
		logrus.SetLevel(logrus.TraceLevel)
		_ = os.Setenv("PRISMA_CLIENT_GO_LOG", "true")
	}

	return nil
}
