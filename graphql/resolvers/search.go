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

// example:
// {
//     search(query: "arisu", args: { type: PROJECT, page: 1 }) {
//         projects { name id description }
//         users { name id description }
//     }
// }
// =>
// {"data":{"projects":[{"name":"arisu","id":"...","description":"Translation made with simplicity, yet robust. - Translation project for Fubuki."}],"users":[]}

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
