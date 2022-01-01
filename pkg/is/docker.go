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

package is

import (
	"io/ioutil"
	"os"
	"strings"
)

// Docker returns if Tsubaki is running under a Docker container.
// This package is a port of `is-docker` by Sindre Sorhus
//
// Package: https://npm.im/is-docker
func Docker() bool {
	return hasDockerEnv() || hasDockerCGroup()
}

func hasDockerEnv() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func hasDockerCGroup() bool {
	contents, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		return false
	}

	return strings.Contains(string(contents), "docker")
}
