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
	"arisu.land/tsubaki/util"
	logustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
)

func init() {
	formatter := internal.NewFormatter()
	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(true)

	// Since logrus is usually initialized here, we can only add Logstash
	// support using the `TSUBAKI_LOGGING_APPENDERS` environment variable.
	//
	// We don't prelude the configuration before executing the command because
	// it's pretty redundant.
	if appenders, ok := os.LookupEnv("TSUBAKI_LOGGING_APPENDERS"); ok {
		actual := strings.Split(appenders, ",")
		if util.Contains(actual, "logstash") {
			if uri, ok := os.LookupEnv("TSUBAKI_LOGGING_LOGSTASH_URI"); ok {
				conn, err := net.Dial("tcp", uri)
				if err != nil {
					panic(err)
				}

				hook := logustash.New(conn, logustash.DefaultFormatter(logrus.Fields{
					"server": "tsubaki",
				}))

				logrus.AddHook(hook)
			}
		}
	}
}

func main() {
	code := tsubaki.Execute()
	os.Exit(code)
}
