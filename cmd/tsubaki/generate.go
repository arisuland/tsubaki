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
	"arisu.land/tsubaki/pkg/storage"
	"arisu.land/tsubaki/util"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "generate <\"jwt\" | \"config\" | \"ssl\">",
		Short:        "Generates a multitude of configurations for you.",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	}

	cmd.AddCommand(newGenerateJwtCommand(), newGenerateConfigCommand())
	return cmd
}

func newGenerateJwtCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "jwt",
		Short:        "Generates a JWT secret key to validate JWT sessions.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			hash := util.GenerateHash(32)
			if hash == "" {
				return errors.New("unable to generate, try again later")
			}

			fmt.Println(hash)
			return nil
		},
	}
}

func newGenerateConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "config [path]",
		Short:        "Generates a config.toml file for you",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If there is a path available, let's generate it there
			path := "./config.yml"
			if len(args) == 1 {
				path = args[0]
			}

			ext := filepath.Ext(path)
			if ext == ".yml" || ext == ".yaml" {
				hash := util.GenerateHash(32)
				if hash == "" {
					return errors.New("unable to generate, try again later")
				}

				defaultConfig := pkg.Config{
					Registrations: true,
					Environment:   pkg.Production,
					Telemetry:     false,
					SecretKeyBase: hash,
					InviteOnly:    false,
					Storage: pkg.StorageConfig{
						Filesystem: &storage.FilesystemStorageConfig{
							Directory: "./.arisu",
						},
					},

					Redis: pkg.RedisConfig{
						Host:    "localhost",
						Port:    6379,
						DbIndex: 7,
					},
				}

				out, err := yaml.Marshal(&defaultConfig)
				if err != nil {
					return err
				}

				fd, err := util.CreateFile(path)
				if err != nil {
					return err
				}

				defer func() {
					_ = fd.Close()
				}()

				_, err = fd.Write(out)
				if err != nil {
					return err
				} else {
					return nil
				}
			} else {
				return fmt.Errorf("path %s didn't end with .yml or .yaml", path)
			}
		},
	}
}
