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
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
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

func (fs FilesystemProvider) GetMetadata(id string, project string) *ProjectMetadata {
	return nil
}

func (fs FilesystemProvider) HandleUpload(files []UploadRequest) error {
	return nil
}
