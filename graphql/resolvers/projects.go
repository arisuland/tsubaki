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
	"arisu.land/tsubaki/graphql/types"
	"arisu.land/tsubaki/prisma/db"
	"context"
	"fmt"
)

func (r *Resolver) Projects(ctx context.Context, args struct {
	Pagination *types.PaginationOptions
}) ([]*types.Project, error) {
	query := r.Container.Database.Client.Project.FindMany()
	if args.Pagination != nil {
		if args.Pagination.Take != nil {
			query = query.Take(int(*args.Pagination.Take))
		}

		if args.Pagination.Skip != nil {
			query = query.Skip(int(*args.Pagination.Skip))
		}

		if args.Pagination.SortBy != nil {
			var order db.SortOrder
			if *args.Pagination.SortBy == types.ASC {
				order = db.SortOrderAsc
			} else if *args.Pagination.SortBy == types.DESC {
				order = db.SortOrderAsc
			} else {
				return nil, fmt.Errorf("unknown sort order: %s", args.Pagination.SortBy.String())
			}

			query = query.OrderBy(db.Project.CreatedAt.Order(order))
		}
	}

	data, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var projects []*types.Project
	for _, i := range data {
		model := types.FromProjectModel(&i)
		projects = append(projects, &model)
	}

	return projects, nil
}
