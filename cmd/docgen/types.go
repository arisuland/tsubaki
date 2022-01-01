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

package main

type SDLDefinition struct {
	//Types []NamedType `json:"types"`
}

type GenericType struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type ObjectType struct {
	Location []int              `json:"location"`
	Fields   []FieldsDefinition `json:"fields"`
	Desc     string             `json:"desc"`
	Name     string             `json:"name"`
}

type FieldsDefinition struct {
	Location []int            `json:"location"`
	Args     []ArgsDefinition `json:"args"`
	Desc     string           `json:"desc"`
	Type     GenericType      `json:"type"`
	Name     string           `json:"name"`
}

type ArgsDefinition struct {
	Location []int       `json:"location"`
	Default  interface{} `json:"default_value"`
	Desc     string      `json:"desc"`
	Type     GenericType `json:"type"`
	Name     string      `json:"name"`
}
