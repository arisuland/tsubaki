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

import (
	"time"

	"arisu.land/tsubaki/prisma/db"
)

// Scope is a Access Token scope.
type Scope string

var (
	// PUBLIC_WRITE is a Scope that allows a OAuth2 application to
	// write any data to any project.
	PUBLIC_WRITE Scope = "PUBLIC_WRITE"

	// REPO_CREATE is a Scope that creates projects beyond the user's behalf.
	REPO_CREATE Scope = "REPO_CREATE"

	// REPO_DELETE is a Scope that allows project deletion beyond the user's behalf.
	REPO_DELETE Scope = "REPO_DELETE"

	// REPO_UPDATE is a Scope that allows updating a project beyond the user's behalf.
	REPO_UPDATE Scope = "REPO_UPDATE"
)

func (s Scope) String() string {
	switch s {
	case PUBLIC_WRITE:
		return "PUBLIC_WRITE"

	case REPO_CREATE:
		return "REPO_CREATE"

	case REPO_DELETE:
		return "REPO_DELETE"

	case REPO_UPDATE:
		return "REPO_UPDATE"

	default:
		return "<unknown>"
	}
}

// AccessToken is a object of a User's access token. This usually
// long-lived for API usage and nothing more. Yes, you can retrieve
// a user's session, but it will only retrieve the session type and
// expiration date, nothing more.
//
// You can see a list of access tokens by a user from the `accessTokens` query
// using the session token or access token to determine the user to grab for.
type AccessToken struct {
	// ExpiresIn returns a ISO-8601 formatted string on when this access
	// token is supposed to live for. This token is deed as "expired"
	// and removed from the database + invalidated if it has expired.
	ExpiresIn *string `json:"expires_in"`

	// Owner returns the User that this AccessToken belongs to.
	Owner User `json:"owner"`

	// Scopes returns a list of Access Token scopes.
	Scopes []Scope `json:"scopes"`

	// Token is the actual JWT token.
	Token string `json:"token"`

	// ID is the current access token ID.
	ID string `json:"id"`
}

func FromAccessTokenDbModel(model *db.AccessTokenModel) AccessToken {
	var expiresIn *string

	owner := model.Owner()
	if data, ok := model.ExpiresIn(); ok {
		owo := data.Format(time.RFC3339)
		expiresIn = &owo
	}

	scopes := make([]Scope, 0)
	for _, scope := range model.Scopes {
		switch scope {
		case db.AccessTokenScopeREPOCREATE:
			scopes = append(scopes, REPO_CREATE)

		case db.AccessTokenScopeREPOUPDATE:
			scopes = append(scopes, REPO_UPDATE)

		case db.AccessTokenScopeREPODELETE:
			scopes = append(scopes, REPO_DELETE)

		case db.AccessTokenScopePUBLICWRITE:
			scopes = append(scopes, PUBLIC_WRITE)
		}
	}

	return AccessToken{
		ExpiresIn: expiresIn,
		Scopes:    scopes,
		Owner:     *FromDbModel(owner),
		Token:     model.Token,
		ID:        model.ID,
	}
}
