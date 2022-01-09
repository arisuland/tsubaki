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
	"arisu.land/tsubaki/prisma/db"
	"time"
)

type Subproject struct {
	Description *string `json:"description"`
	UpdatedAt   string  `json:"updatedAt"`
	CreatedAt   string  `json:"createdAt"`
	Parent      Project `json:"parent"`
	Name        string  `json:"name"`
	ID          string  `json:"id"`
}

func FromSubprojectDbModel(model *db.SubprojectModel) Subproject {
	var desc *string

	description, ok := model.Description()
	if ok {
		desc = &description
	}

	parent := model.Parent()
	return Subproject{
		Description: desc,
		UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		Parent:      FromProjectModel(parent),
		Name:        model.Name,
		ID:          model.ID,
	}
}
