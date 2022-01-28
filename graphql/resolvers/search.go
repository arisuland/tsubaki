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

package resolvers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"arisu.land/tsubaki/graphql/types/result"
	"github.com/sirupsen/logrus"
)

// example:
// {
//     search(query: "arisu", args: { type: PROJECT, page: 1 }) {
//         projects { name id description }
//         users { name id description }
//     }
// }
// =>
// {"data":{"projects":[{"name":"arisu","id":"...","description":"Translation made with simplicity, yet robust. - Translation project for Fubuki."}],"users":[]}

func (r *Resolver) Search(ctx context.Context, args struct {
	Query string
}) (*result.QueryResult, error) {
	if r.Container.ElasticSearch == nil {
		return nil, errors.New("elasticsearch is not available on this instance")
	}

	// build the query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": args.Query,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := r.Container.ElasticSearch.Search(
		r.Container.ElasticSearch.Search.WithContext(context.TODO()),
		r.Container.ElasticSearch.Search.WithIndex("arisu:tsubaki:users", "arisu:tsubaki:projects"),
		r.Container.ElasticSearch.Search.WithBody(&buf),
		r.Container.ElasticSearch.Search.WithTrackTotalHits(true),
		r.Container.ElasticSearch.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			logrus.Errorf("Unable to request to ElasticSearch [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)

			return nil, fmt.Errorf("unable to request to es: [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil, nil
}

/*
  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
    log.Fatalf("Error parsing the response body: %s", err)
  }
  // Print the response status, number of results, and request duration.
  log.Printf(
    "[%s] %d hits; took: %dms",
    res.Status(),
    int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
    int(r["took"].(float64)),
  )
  // Print the ID and document source for each hit.
  for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
    log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
  }

  log.Println(strings.Repeat("=", 37))
*/

// Search is a GraphQL query to return data about a user or project
// that is searched on the dashboard.
//func (r *Resolver) Search(ctx context.Context, args struct {
//	Query string
//	Args  struct {
//		Type types.SearchQueryType
//		Page *int32
//	}
//}) (*types.SearchQueryResult, error) {
//	if pkg.GlobalContainer.ElasticSearch == nil {
//		return nil, errors.New("elasticsearch is not available on this instance")
//	}
//
//	// Now we check if the query type is unknown!
//	if args.Args.Type == types.UNKNOWN {
//		//pkg.GlobalContainer.ElasticSearch.Search
//	}
//
//	return nil, nil
//}
