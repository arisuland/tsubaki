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

// BaseStorageProvider represents the bare-bones methods of what a storage provider should be.
type BaseStorageProvider interface {
	// HandleUpload is a function to handle file uploads to this specific
	// BaseStorageProvider instance.
	HandleUpload()

	// GetMetadata is a function to retrieve metadata about this project. This is usually
	// embedded under `id/project/metadata.json`.
	GetMetadata(id string, project string) *ProjectMetadata

	// Name returns the name of this BaseStorageProvider.
	Name() string
}

// ProjectMetadata represents the metadata stored under `id/project/metadata.json`
type ProjectMetadata struct {
	// FormatVersion returns the specific format version this metadata table
	// is in.
	FormatVersion int `json:"format_version"`

	// Description is the description of the project.
	Description string `json:"description"`

	// Owner returns the owner's ID to whoever owns this project.
	Owner string `json:"owner"`

	// Files returns the files of this project.
	Files []FileMetadata `json:"files"`

	// Path returns the storage path if using GCS or S3, or the absolute
	// path if using the Filesystem provider.
	Path string `json:"path"`

	// Name returns the project name
	Name string `json:"name"`
}

// FileMetadata returns the metadata about a specific file.
type FileMetadata struct {
	// ContentType returns the content type of this file.
	ContentType string `json:"content_type"`

	// Path returns the storage path if using GCS or S3, or the relative
	// path if using the Filesystem provider.
	Path string `json:"path"`

	// Size returns in bytes, how big the file is.
	Size int `json:"size"`
}

// Config is the configuration for storing projects.
//
// Prefix: TSUBAKI_STORAGE
type Config struct {
	// Configures using the filesystem to host your projects, once the data
	// is removed, Arisu will fix it but cannot restore your projects.
	// If you're using Docker or Kubernetes, it is best assured that you
	// must create a volume so Arisu can interact with it.
	//
	// Aliases: `fs` | Prefix: TSUBAKI_STORAGE_FS_*
	Filesystem FilesystemStorageConfig `json:"filesystem" json:"fs"`
}
