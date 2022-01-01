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
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"math/rand"
	"strings"
)

// GeneratePassword creates a password based off the Argon2 specification using
// the golang.org/x/crypto/argon2 package.
func GeneratePassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, 64*1024, 1, 4, b64Salt, b64Hash), nil
}

// VerifyPassword verifies the password to decode it and check if it's valid
// from the database entry.
func VerifyPassword(password string, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	fmt.Println(hash, parts[3])
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", 64*1024, 1, 4)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decoded, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	keyLen := uint32(len(decoded))
	compare := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLen)

	return subtle.ConstantTimeCompare(decoded, compare) == 1, nil
}
