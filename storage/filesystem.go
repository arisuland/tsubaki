// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2021 Noelware
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

package storage

// FilesystemProvider is a BaseStorageProvider for using the local filesystem
// to handle projects in.
type FilesystemProvider struct {
	// Directory is the directory to use to store projects. An `arisu.lock` file will be
	// generated if it doesn't exist.
	Directory string
}

// FilesystemStorageConfig is the configuration for the FilesystemProvider.
type FilesystemStorageConfig struct {
	// Directory is the directory to use to store projects. An `arisu.lock` file will be
	// generated if it doesn't exist.
	Directory string `json:"directory"`
}
