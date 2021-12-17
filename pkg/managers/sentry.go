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

package managers

import (
	"arisu.land/tsubaki/pkg"
	"fmt"
	"github.com/getsentry/sentry-go"
)

// SentryManager is a struct that provides error handling and outputs
// it to Sentry. This is used in the error handling middleware.
type SentryManager struct {
	// Client is the Sentry client if the DSN is provided,
	// or `nil` if no Sentry DSN is present.
	Client *sentry.Client
}

func NewSentryManager(config *Config) (SentryManager, error) {
	if config.SentryDsn == nil {
		return SentryManager{
			Client: nil,
		}, nil
	}

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              *config.SentryDsn,
		AttachStacktrace: true,
		SampleRate:       1.0,
		ServerName:       fmt.Sprintf("arisu.tsubaki v%s", pkg.Version),
	})

	if err != nil {
		return SentryManager{
			Client: nil,
		}, err
	}

	return SentryManager{
		Client: client,
	}, nil
}
