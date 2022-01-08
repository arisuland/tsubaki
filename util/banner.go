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

package util

import (
	"arisu.land/tsubaki/internal"
	"html/template"
	"io/ioutil"
	"os"
)

type BannerTemplateData struct {
	Telemetry  bool
	CommitHash string
	Version    string
	BuildDate  string
}

func PrintBanner() {
	contents, err := ioutil.ReadFile("./assets/banner.txt")
	if err != nil {
		panic(err)
	}

	value := string(contents)
	t := template.New("banner template")

	data := BannerTemplateData{
		Version:    internal.Version,
		CommitHash: internal.CommitSHA,
		BuildDate:  internal.BuildDate,
	}

	ta, err := t.Parse(value)
	if err != nil {
		panic(err)
	}

	if err := ta.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
