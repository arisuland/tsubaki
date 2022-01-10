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

package pkg

import "arisu.land/tsubaki/prisma/db"

type IndexedUser struct {
	Description *string `json:"description"`
	Name        *string `json:"name"`
	ID          string  `json:"id"`
}

func FromDbUserModel(model db.UserModel) IndexedUser {
	var name *string
	var desc *string

	if d, ok := model.Description(); ok {
		desc = &d
	}

	if n, ok := model.Name(); ok {
		name = &n
	}

	return IndexedUser{
		Description: desc,
		Name:        name,
		ID:          model.ID,
	}
}
