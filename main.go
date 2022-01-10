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

package main

import (
	"arisu.land/tsubaki/cmd/tsubaki"
	"arisu.land/tsubaki/internal"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	formatter := internal.NewFormatter()
	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(true)
}

func main() {
	code := tsubaki.Execute()
	os.Exit(code)
}
