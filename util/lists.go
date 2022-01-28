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

// FindIndexIteratee is a function to call based off the FindIndex function.
type FindIndexIteratee = func(item interface{}) bool

// IndexOf returns the first index of needle in the haystack.
// Or -1 if needle was not in the haystack.
func IndexOf(haystack []interface{}, needle interface{}) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}

	return -1
}

// FindIndex finds the index of a value from the haystack using a FindIndexIteratee.
// This will traverse through the haystack and return the index from the needleIter
// or -1 if the needleIter hasn't concluded anything in the haystack.
func FindIndex(haystack interface{}, needleIter FindIndexIteratee) int {
	for i, v := range haystack.([]interface{}) {
		res := needleIter(v)
		if res {
			return i
		}
	}

	return -1
}

// TODO: use go generics for this function

func Contains(haystack []string, item string) bool {
	for _, v := range haystack {
		if v == item {
			return true
		}
	}

	return false
}
