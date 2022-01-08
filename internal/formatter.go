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

package internal

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Formatter struct {
	// DisableColors is when we need to disable colour output.
	//
	// This can be overrided using the `TSUBAKI_DISABLE_COLORS` environment
	// variable.
	DisableColors bool
}

var format = "Jan 02, 2006 - 15:04:05 MST"

// NewFormatter creates a new Formatter instance.
func NewFormatter() *Formatter {
	var colorsEnabled = true
	if os.Getenv("TSUBAKI_DISABLE_COLORS") != "" {
		colorsEnabled = false
	}

	return &Formatter{
		DisableColors: colorsEnabled,
	}
}

// Format renders a single log entry for logrus.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := make(logrus.Fields)
	for k, v := range entry.Data {
		fields[k] = v
	}

	level := f.getColourForLevel(entry.Level)
	b := &bytes.Buffer{}

	if f.DisableColors {
		_, _ = fmt.Fprintf(b, "[%s] ", entry.Time.Format(format))
	} else {
		_, _ = fmt.Fprintf(b, "\x1b[38;2;134;134;134m%s \x1b[0m", entry.Time.Format(format))
	}

	l := strings.ToUpper(entry.Level.String())
	if f.DisableColors {
		b.WriteString(level)
		b.WriteString("[" + l[:4] + "] ")
		b.WriteString("\x1b[0m")
	} else {
		b.WriteString("[" + l[:4] + "] ")
	}

	if len(fields) != 0 {
		for f, v := range fields {
			_, _ = fmt.Fprintf(b, "[%s=%v] ", f, v)
		}
	}

	b.WriteString(" ")
	if entry.HasCaller() {
		var pkg string
		if strings.HasPrefix(entry.Caller.Function, "arisu.land/tsubaki") {
			pkg = strings.TrimPrefix(entry.Caller.Function, "arisu.land/tsubaki/")
		} else {
			pkg = entry.Caller.Function
		}

		if f.DisableColors {
			_, _ = fmt.Fprintf(b, "[%s (%s:%d)] ", pkg, entry.Caller.File, entry.Caller.Line)
		} else {
			_, _ = fmt.Fprintf(b, "\x1b[38;2;134;134;134m[%s (%s:%d)]\x1b[0m", pkg, entry.Caller.File, entry.Caller.Line)
		}

		b.WriteString(" ")
	}

	b.WriteString(strings.TrimSpace(entry.Message))
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) getColourForLevel(level logrus.Level) string {
	if f.DisableColors {
		return ""
	}

	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		// #A3B68A
		return "\x1b[1m\x1b[38;2;163;182;138m"

	case logrus.ErrorLevel, logrus.FatalLevel:
		// #994B68
		return "\x1b[1m\x1b[38;2;153;75;104m"

	case logrus.WarnLevel:
		// #F3F386
		return "\x1b[1m\x1b[38;2;243;243;134m"

	case logrus.InfoLevel:
		// #B29DF3
		return "\x1b[1m\x1b[38;2;178;157;243m"

	default:
		// #2f2f2f
		return "\x1b[1m\x1b[38;2;47;47;47m"
	}
}
