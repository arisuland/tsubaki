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

package result

// LoginResult is a Result object but with a `token` property.
type LoginResult struct {
	// Success represents if the result of this action
	// was successful or not.
	Success bool `json:"success"`

	// Errors represents a list of Error objects
	// if `success=false`.
	Errors []Error `json:"errors"`

	// Token results the session token, you will be only
	// able to see this once.
	Token string `json:"token"`
}
