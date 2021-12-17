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
