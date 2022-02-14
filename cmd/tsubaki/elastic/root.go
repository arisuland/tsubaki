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
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"strings"
)

func NewElasticCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "elastic <subcommand> [...args]",
		Short:        "Command to use Elasticsearch to index all users / projects and to update the index settings.",
		SilenceUsage: true,
		Long: `The "elastic" command is a one-way gateway to interact with Elasticsearch
when using Tsubaki. By default, Tsubaki doesn't create an ES client so you can start using the /api/v1/search
endpoint, so this is completely up to you.

The "index" subcommand, by default, will create the indexes for you, in which you can point to "es-master:9200/<index>"
to get information about the newly indexes. Or, it will connect to PostgreSQL and index all the users and projects for
you.

The "settings" subcommand updates the index settings without you writing the settings by hand from Kibana's Development
Tools page or using cURL.
`,
	}

	cmd.PersistentFlags().StringP("endpoints", "e", "localhost:9200", "Connects to your Elasticsearch cluster.")
	cmd.PersistentFlags().StringP("username", "u", "", "If Basic authentication is enabled, this is the username you use to connect to Elasticsearch.")
	cmd.PersistentFlags().StringP("password", "p", "", "If Basic authentication is enabled, this is the password you use to connect to Elasticsearch.")
	cmd.AddCommand(
		newIndexCommand(),
		newIndexSettingsCommand(),
	)

	return cmd
}

func createElasticClient(cmd *cobra.Command) (*elasticsearch.Client, error) {
	addresses := strings.Split(cmd.Flag("endpoints").Value.String(), ",")
	username := ""
	password := ""

	uflag := cmd.Flag("username")
	pflag := cmd.Flag("password")

	if uflag != nil && uflag.Value.String() != "" {
		username = uflag.Value.String()
	}

	if pflag != nil && pflag.Value.String() != "" {
		password = pflag.Value.String()
	}

	config := elasticsearch.Config{
		Addresses: addresses,
	}

	if username != "" {
		config.Username = username
	}

	if password != "" {
		config.Password = password
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to elasticsearch: %v", err)
	}

	// Try to create a request
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve server information: %v", err)
	}

	// unmarshal it from `res`
	defer func() {
		_ = res.Body.Close()
	}()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	serverVersion := data["version"].(map[string]interface{})["number"].(string)
	fmt.Printf("Server: %s | Client: %s\n", serverVersion, elasticsearch.Version)
	return client, nil
}
