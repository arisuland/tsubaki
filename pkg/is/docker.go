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
