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

package util

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var db = make(map[string]mime)

// mime is represented as the mime type from the mime database.
type mime struct {
}

func init() {
	if err := populateDb(); err != nil {
		logrus.Fatalf("Unable to populate mime type database.")
		panic(err)
	}
}

func populateDb() error {
	logrus.Debugf("Now populating mime type database!")

	iana := "iana"
	apache := "apache"
	nginx := "nginx"
	_ = []*string{&nginx, &apache, nil, &iana}

	// Load in the file
	_, err := ioutil.ReadFile("./assets/mime-db.json")
	if err != nil {
		return err
	}

	return nil
}
