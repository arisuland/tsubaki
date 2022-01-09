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

import (
	"arisu.land/tsubaki/util"
	"encoding/json"
	"fmt"
	"github.com/mushroomsir/mimetypes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FilesystemProvider struct {
	Directory string
}

type FilesystemStorageConfig struct {
	Directory string `yaml:"directory"`
}

func NewFilesystemStorageProvider(config FilesystemStorageConfig) BaseStorageProvider {
	return FilesystemProvider{
		Directory: config.Directory,
	}
}

func (fs FilesystemProvider) Init() error {
	logrus.Infof("Checking if directory %s exists...", fs.Directory)

	_, err := os.Stat(fs.Directory)
	if os.IsNotExist(err) {
		logrus.Warnf("Directory %s doesn't exist, creating...", fs.Directory)
		err = os.MkdirAll(filepath.Dir(fs.Directory), 0755)
		if err != nil {
			return err
		}
	}

	logrus.Info("Directory exists! Checking if lockfile exists...")
	path := fs.Directory + "/arisu.lock"

	if _, err = os.Stat(path); err != nil {
		logrus.Warn("Manifest file is missing or corrupted, overriding...")

		file, err := util.CreateFile(path)
		if err != nil {
			return err
		}

		logrus.Infof("Creating file %s!", path)
		err = os.WriteFile(path, []byte("... this file exists to exist ...\n"), 0755)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()
	}

	logrus.Info("Found lockfile! Checking if it was tampered with...")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(contents)
	if content != "... this file exists to exist ...\n" {
		logrus.Warn("Manifest lockfile was tampered or was corrupted! Rewriting...")
		err = os.WriteFile(path, []byte("... this file exists to exist ...\n"), 0755)
		if err != nil {
			return err
		}
	}

	logrus.Info("Everything looks okay~")
	return nil
}

func (fs FilesystemProvider) Name() string {
	return "filesystem"
}

func (fs FilesystemProvider) GetMetadata(id string, project string) (*ProjectMetadata, error) {
	logrus.Infof("Told to grab metadata for project %s/%s", id, project)

	// Check if the directory exists
	dir := fmt.Sprintf("%s/%s/%s", fs.Directory, id, project)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		logrus.Warnf("Project doesn't have a directory for it. Now creating!")
		err = os.MkdirAll(filepath.Dir(dir), 0755)
		if err != nil {
			return nil, err
		}
	}

	// Check if "metadata.lock" exists
	logrus.Infof("Checking if metadata lock exists for project %s/%s", id, project)
	path := dir + "/metadata.lock"

	if _, err = os.Stat(path); err != nil {
		logrus.Warnf("Manifest file is missing for project %s/%s, creating!", id, project)
		file, err := util.CreateFile(path)
		if err != nil {
			return nil, err
		}

		logrus.Debugf("Created file %s!", path)

		metadata := &ProjectMetadata{
			FormatVersion: FormatV1,
			Description:   "",
			Owner:         id,
			Files:         []FileMetadata{},
			Path:          dir,
			Name:          project,
		}

		data, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(path, data, 0755)
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = file.Close()
		}()

		return metadata, nil
	}

	logrus.Infof("Now retrieving content for %s!", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata *ProjectMetadata
	err = json.Unmarshal(contents, metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (fs FilesystemProvider) HandleUpload(files []UploadRequest) error {
	logrus.Infof("Told to handle %d files!", len(files))
	s := time.Now()

	for _, file := range files {
		logrus.Debugf("Taking care of file %s for project %s/%s", file.Name, file.Owner, file.Project)

		// Retrieving the metadata lock will create the directory + file itself.
		m, err := fs.GetMetadata(file.Project, file.Owner)
		if err != nil {
			return err
		}

		logrus.Infof("Using format version %d for project %s/%s", m.FormatVersion.Int(), file.Owner, file.Project)

		// Figure out the mime type
		var mimeType string
		if strings.Contains(file.Name, ".") {
			mimeType = mimetypes.Lookup(file.Name)
		} else {
			// TODO: probably check for shell bangs (!#) to determine
			// it also.
			mimeType = "application/octet-stream"
		}

		logrus.Infof("Figured out that file %s has a mime type of %s.", file.Name, mimeType)

		// Check if the file exists in the directory
		// file.Name should be `folder/file.js` or `file.js` so it can be appended
		// as `<dir>/<owner>/<project>/folder/file.js`
		dir := fmt.Sprintf("%s/%s/%s/%s", fs.Directory, file.Owner, file.Project, file.Name)
		_, err = os.Stat(dir)
		if os.IsNotExist(err) {
			logrus.Warnf("File at %s doesn't exist.", dir)

			f, err := util.CreateFile(dir)
			if err != nil {
				logrus.Warnf("Unable to handle file creation for file %s: %v", dir, err)
				return err
			}

			err = os.WriteFile(dir, []byte(file.Contents), 0755)
			if err != nil {
				logrus.Warnf("Unable to handle file update for file %s: %v", dir, err)
				return err
			}

			// shouldn't care about files closing (for now)
			_ = f.Close()
		} else {
			logrus.Infof("Updating file content for %s...", dir)
			err = os.WriteFile(dir, []byte(file.Contents), 0755)
			if err != nil {
				logrus.Warnf("Unable to handle file update for file %s: %v", dir, err)
				return err
			}
		}

		// find metadata for file
		index := util.FindIndex(m.Files, func(item interface{}) bool {
			casted := item.(FileMetadata)
			return casted.Path == file.Name
		})

		if index != -1 {
			m.Files[index] = FileMetadata{
				Path:        file.Name,
				ContentType: mimeType,
				Size:        file.Size,
			}

			bytes, err := json.Marshal(&m)
			if err != nil {
				logrus.Errorf("Unable to unmarshal metadata.lock file: %v", err)
				return err
			}

			err = os.WriteFile(
				fmt.Sprintf("%s/%s/%s/metadata.lock", fs.Directory, file.Owner, file.Project),
				bytes,
				0755,
			)

			if err != nil {
				logrus.Warnf("Unable to handle file update for file %s: %v", fmt.Sprintf("%s/%s/%s/metadata.lock", fs.Directory, file.Owner, file.Project), err)
				return err
			}
		}
	}

	logrus.Infof("Took %s to complete %d files.", time.Since(s).String(), len(files))
	return nil
}
