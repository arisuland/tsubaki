// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2021 Noelware
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

var Codes = map[int]string{
	1001: "Invalid username. All usernames must be in a range of 8-16 characters.",
	1002: "Invalid password provided.",
	1003: "Email is already taken.",
	1004: "Username is already taken.",
	1005: "User is not logged in.",
	1006: "The subproject parent ID is not equal to the parent project.",
	1007: "Subproject is already created.",
	1008: "Length for string was over the maximum amount of characters.",
}

func Get(code int) (int, string) {
	for k, v := range Codes {
		if k == code {
			return k, v
		}
	}

	return 0, "<unknown>"
}
