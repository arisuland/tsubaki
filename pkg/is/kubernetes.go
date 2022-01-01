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

func hasKubeEnv() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

func hasServiceAccountFile() bool {
	var hasServiceAccountToken = false
	var hasServiceAccountNS = false

	_, err := os.Stat("/run/secrets/kubernetes.io/serviceaccount/token")
	if err == nil {
		hasServiceAccountToken = true
	}

	_, err = os.Stat("/run/secrets/kubernetes.io/serviceaccount/token")
	if err == nil {
		hasServiceAccountNS = true
	}

	return hasServiceAccountToken && hasServiceAccountNS
}

func hasClusterDns() bool {
	contents, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		return false
	}

	return strings.Contains(string(contents), "cluster.local")
}

// Kubernetes returns a boolean if Tsubaki is running under Kubernetes.
func Kubernetes() bool {
	return hasKubeEnv() || hasClusterDns() || hasServiceAccountFile()
}
