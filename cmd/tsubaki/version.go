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
	"arisu.land/tsubaki/internal"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
)

type versionInfo struct {
	Version    string `json:"version"`
	CommitSHA  string `json:"commit_sha"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	GoCompiler string `json:"go_compiler"`
	Platform   string `json:"platform"`
}

func newVersionCommand() *cobra.Command {
	var (
		printJson bool
		pretty    bool
	)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Returns the current version of Tsubaki.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if printJson {
				data := &versionInfo{
					Version:    internal.Version,
					CommitSHA:  internal.CommitSHA,
					BuildDate:  internal.BuildDate,
					GoVersion:  strings.TrimPrefix(runtime.Version(), "go"),
					GoCompiler: runtime.Compiler,
					Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
				}

				if pretty {
					b, err := json.MarshalIndent(data, "", "  ")
					if err != nil {
						return err
					}

					fmt.Println(string(b))
					return nil
				} else {
					b, err := json.Marshal(data)
					if err != nil {
						return err
					}

					fmt.Println(string(b))
					return nil
				}
			}

			fmt.Printf("tsubaki v%s (commit: %s, built at: %s) on %s/%s\n", internal.Version, internal.CommitSHA, internal.BuildDate, runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&printJson, "json", "j", false, "outputs the version in a JSON format.")
	cmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "indents the version if `--json` is true.")

	return cmd
}
