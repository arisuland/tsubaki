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

package storage

// FormatVersion refers to the format version of the `metadata.lock` file.
type FormatVersion int

// FormatV1 is the first format version of the metadata lock file.
var FormatV1 FormatVersion = 1

// Int returns the integer value of a specific FormatVersion.
func (t FormatVersion) Int() int {
	switch {
	case t == FormatV1:
		return 1

	default:
		return -1
	}
}

// BaseStorageProvider represents the bare-bones methods of what a storage provider should be.
type BaseStorageProvider interface {
	// HandleUpload is a function to handle file uploads to this specific
	// BaseStorageProvider instance.
	HandleUpload(files []UploadRequest) error

	// GetMetadata is a function to retrieve metadata about this project. This is usually
	// embedded under `id/project/metadata.lock`.
	GetMetadata(id string, project string) (*ProjectMetadata, error)

	// Init is a function to call to initialize the provider.
	Init() error

	// Name returns the name of this BaseStorageProvider.
	Name() string
}

// ProjectMetadata represents the metadata stored under `id/project/metadata.json`
type ProjectMetadata struct {
	// FormatVersion returns the specific format version this metadata table
	// is in.
	FormatVersion FormatVersion `json:"format_version"`

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

// UploadRequest is a object that represents a request
// to upload one or more files into this BaseStorageProvider.
type UploadRequest struct {
	// ContentType refers to the content type that this file is.
	// So it can properly type the `Content-Type` header in
	// /api/v1/storage/file/:user/:project/...:path
	ContentType string

	// Contents refers to the actual content of the file.
	Contents string

	// Project is the project's ID.
	Project string

	// Owner is the project owner's ID.
	Owner string

	// Name is the file name.
	Name string

	// Size is how big the file is. By default,
	// Fubuki will not load the editor if the file is over 1GB.
	Size int64
}
