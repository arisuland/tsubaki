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

// Keys returns a []string of all the keys found in the given map.
func Keys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// Find performs an operation to iterate all the Keys in a given map
// and return the value as `interface{}`, if the key isn't found, it will
// return `nil`.
func Find(m map[string]interface{}, key string) *interface{} {
	keys := Keys(m)
	for _, k := range keys {
		if k == key {
			val := m[k]
			return &val
		}
	}

	return nil
}
