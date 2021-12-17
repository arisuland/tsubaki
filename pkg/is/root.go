package is

import "os"

// Root returns if Tsubaki was ran using root or with `sudo`.
// This package is a port of `is-root` by Sindre Sorhus
//
// Package: https://npm.im/is-root
func Root() bool {
	id := os.Getuid()
	return id == 0
}
