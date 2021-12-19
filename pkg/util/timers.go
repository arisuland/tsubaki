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

package util

import "time"

// SetInterval creates a new timer to run code in a goroutine
// based off this Stackoverflow thread: https://stackoverflow.com/a/16466581
func SetInterval(run func(), t time.Duration) chan struct{} {
	ticker := time.NewTicker(t)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				run()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
