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

import (
	"arisu.land/tsubaki/pkg/util"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var log = slog.Make(sloghuman.Sink(os.Stdout))

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

func NewFilesystemStorageProvider(config FilesystemStorageConfig) BaseStorageProvider {
	return FilesystemProvider{
		Directory: config.Directory,
	}
}

func (fs FilesystemProvider) Init() error {
	log.Info(context.Background(), "Checking if directory exists...")
	_, err := os.Stat(fs.Directory)

	if os.IsNotExist(err) {
		log.Warn(context.Background(), "Directory doesn't exist, creating...")
		err = os.MkdirAll(filepath.Dir(fs.Directory), 0755)
		if err != nil {
			return err
		}
	}

	//if !stat.IsDir() {
	//	return errors.New(fmt.Sprintf("directory %s was not a valid directory.", fs.Directory))
	//}

	log.Info(context.Background(), "Checking if `arisu.lock` exists...")
	path := fmt.Sprintf("%s/arisu.lock", fs.Directory)

	// this looks ugly, yes.
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Info(context.Background(), "Lockfile is missing or is corrupted, overriding...")
		file, err := util.CreateFile(path)
		if err != nil {
			return nil
		}

		log.Info(context.Background(), fmt.Sprintf("Creating file %s!", path))
		err = os.WriteFile(path, []byte("... this file exists as a lockfile ...\\n"), 0755)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return nil
		}
	}

	log.Info(context.Background(), "Checking if lockfile contents is valid...")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(contents)
	if content != "... this file exists as a lockfile ...\n" {
		log.Info(context.Background(), "Corrupted lockfile (probably was tampered), rewriting...")
		err = os.WriteFile(path, []byte("... this file exists as a lockfile ...\n"), 0755)
		if err != nil {
			return err
		}
	}

	log.Info(context.Background(), "Everything looks alright!")
	return nil
}

func (fs FilesystemProvider) Name() string {
	return "filesystem"
}

func (fs FilesystemProvider) GetMetadata(id string, project string) *ProjectMetadata {
	return nil
}

func (fs FilesystemProvider) HandleUpload() {
	// TODO: this
}
