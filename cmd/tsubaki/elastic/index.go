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

package elastic

import (
	"arisu.land/tsubaki/prisma/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var index string

func newIndexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index [subcommand] [...args]",
		Short: "Creates the index indices or index your user and projects data.",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.PersistentFlags().StringVarP(&index, "index", "i", "", "The index to use.")
	cmd.AddCommand(newCreateIndexCommand(), newIndexDataCommand())

	return cmd
}

func newCreateIndexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "Creates the newly index and update the mapping of it.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Creating required indices...")

			// Create an Elasticsearch client
			client, err := createElasticClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to connect to elasticsearch: %v", err)
			}

			projectMappings := `{
				"mappings": {
					"_doc": {
						"properties": {
							"description": { "type": "text" },
							"created_at": { "type": "date" },
							"owner_id": { "type": "text" },
							"name": { "type": "text" },
							"id": { "type": "text" }
						}
					}
				}
			}`

			// Create the projects index
			res, err := client.Indices.Create(
				"tsubaki-projects",
				client.Indices.Create.WithBody(strings.NewReader(projectMappings)),
				client.Indices.Create.WithErrorTrace(),
			)

			if err != nil {
				return err
			}

			// unmarshal it from `res`
			defer func() {
				_ = res.Body.Close()
			}()

			var data map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				return err
			}

			if res.IsError() {
				return fmt.Errorf("unable to create index \"tsubaki-projects\":\n%v", data)
			} else {
				fmt.Println("Created the \"tsubaki-projects\" index!")
			}

			// Create the users index
			usersMappings := `{
				"mappings": {
					"_doc": {
						"properties": {
							"description": { "type": "text" },
							"created_at": { "type": "date" },
							"name": { "type": "text" },
							"id": { "type": "text" }
						}
					}
				}
			}`

			res2, err := client.Indices.Create(
				"tsubaki-users",
				client.Indices.Create.WithBody(strings.NewReader(usersMappings)),
				client.Indices.Create.WithErrorTrace(),
			)

			if err != nil {
				return err
			}

			// unmarshal it from `res`
			defer func() {
				_ = res2.Body.Close()
			}()

			var data2 map[string]interface{}
			err = json.NewDecoder(res2.Body).Decode(&data2)
			if err != nil {
				return err
			}

			if res.IsError() {
				return fmt.Errorf("unable to create index \"tsubaki-projects\":\n%v", data2)
			} else {
				fmt.Println("Created the \"tsubaki-users\" index!")
			}

			fmt.Println("Looks like we are done!")
			return nil
		},
	}

	return cmd
}

func newIndexDataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "data",
		Short:        "Indexes all the data from PostgreSQL",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Now indexing projects and users!")

			workers, err := strconv.Atoi(cmd.Flag("workers").Value.String())
			if err != nil {
				return err
			}

			fmt.Printf("Using %d workers...\n", workers)

			// Create the Elasticsearch client
			_, err = createElasticClient(cmd)
			if err != nil {
				return err
			}

			// Create the PostgreSQL client
			prisma := db.NewClient()
			if err := prisma.Connect(); err != nil {
				return err
			}

			fmt.Println("Connected to PostgreSQL!")

			projects, err := prisma.Project.FindMany().Exec(context.TODO())
			if err != nil {
				return err
			}

			_, err = prisma.User.FindMany().Exec(context.TODO())
			if err != nil {
				return err
			}

			var wg sync.WaitGroup
			for i, project := range projects {
				fmt.Printf("[%d/%d] Indexing project %s (%s)...\n", i+1, len(projects), project.Name, project.ID)

				wg.Add(1)
				go func(i int, proj db.ProjectModel) {
					_ = time.Now()

					// Marshal it from JSON
					_, _ = json.Marshal(&proj.InnerProject)
				}(i, project)
			}

			return nil
		},
	}

	cmd.Flags().IntP("workers", "w", runtime.NumCPU(), "how many workers the indexing should use.")
	return cmd
}
