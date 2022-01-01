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

func fallbackEnvString(envString string, fallback string) string {
	if envString == "" {
		return fallback
	} else {
		return envString
	}
}

func PrintBanner(version string, commitHash string, buildDate string) {
	contents, err := ioutil.ReadFile("./assets/banner.txt")
	if err != nil {
		panic(err)
	}

	value := string(contents)
	t := template.New("banner template")
	telemetry := fallbackEnvString(os.Getenv("TSUBAKI_TELEMETRY"), "no")

	var enabled bool
	if telemetry == "no" || telemetry == "false" {
		enabled = false
	} else {
		enabled = true
	}

	data := BannerTemplateData{
		Version:    version,
		CommitHash: commitHash,
		Telemetry:  enabled,
		BuildDate:  buildDate,
	}

	ta, err := t.Parse(value)
	if err != nil {
		panic(err)
	}

	if err := ta.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
