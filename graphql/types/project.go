package types

import (
	"time"

	"arisu.land/tsubaki/prisma/db"
)

// Project is a project that belongs to a User or Organization.
type Project struct {
	Description *string      `json:"description"`
	UpdatedAt   string       `json:"updated_at"`
	CreatedAt   string       `json:"created_at"`
	Subprojects []Subproject `json:"subprojects"`
	Flags       int32        `json:"flags"`
	Owner       User         `json:"owner"`
	Name        string       `json:"name"`
	ID          string       `json:"id"`
}

// FromProjectModel returns a new Project entity based off the db result.
func FromProjectModel(model *db.ProjectModel) Project {
	var desc *string
	description, ok := model.Description()
	if !ok {
		desc = nil
	} else {
		desc = &description
	}

	subprojects := make([]Subproject, 0)
	if model.RelationsProject.Subprojects != nil {
		for _, sub := range model.RelationsProject.Subprojects {
			subprojects = append(subprojects, FromSubprojectDbModel(&sub))
		}
	}

	owner := model.Owner()
	return Project{
		Description: desc,
		UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		Subprojects: subprojects,
		Owner:       FromDbModel(owner),
		Flags:       int32(model.Flags),
		Name:        model.Name,
		ID:          model.ID,
	}
}
