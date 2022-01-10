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

package types

import "arisu.land/tsubaki/prisma/db"

// SearchQueryType is a enum representing "Project" and "User".
// This is defined as a GraphQL enum object.
type SearchQueryType string

var (
	// UNKNOWN is a SearchQueryType to query a user OR project,
	// doesn't really care about what is returned.
	UNKNOWN SearchQueryType = "Unknown"

	// PROJECT is a SearchQueryType to query a project by its
	// ID or name.
	PROJECT SearchQueryType = "Project"

	// USER is a SearchQueryType to query a user by its
	// ID or name.
	USER SearchQueryType = "User"
)

// String stringifies a SearchQueryType value.
func (s SearchQueryType) String() string {
	switch s {
	case PROJECT:
		return "Project"

	case USER:
		return "User"

	default:
		return "<unknown>"
	}
}

// SearchedItem is a disected User to not expose too much data
// since the search query requires no authentication.
type SearchedItem struct {
	Description *string `json:"description"`
	Name        *string `json:"name"`
	ID          string  `json:"id"`
}

// SearchQueryResult is a GraphQL object to determine the results
// from ElasticSearch.
type SearchQueryResult struct {
	// Projects returns a list of projects if the SearchQueryType was PROJECT.
	Projects []Project `json:"projects"`

	// Users returns a list of users if the SearchQueryType was USER.
	Users []SearchedItem `json:"users"`
}

// ConvertUserToItem converts a db.UserModel into a SearchedItem object.
func ConvertUserToItem(user *db.UserModel) SearchedItem {
	var desc *string
	var name *string

	n, ok := user.Name()
	if ok {
		name = &n
	}

	d, ok := user.Description()
	if ok {
		desc = &d
	}

	return SearchedItem{
		Description: desc,
		Name:        name,
		ID:          user.ID,
	}
}
