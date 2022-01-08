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

package pkg

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

// NewToken creates a new JWT token with the user's ID as the mapped claims.
func NewToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": uid,
	})

	signed, err := token.SignedString([]byte(GlobalContainer.Config.SecretKeyBase))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateToken(token string) (bool, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(GlobalContainer.Config.SecretKeyBase), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return true, nil
	}

	return false, errors.New("unknown error has occurred")
}

func DecodeToken(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(GlobalContainer.Config.SecretKeyBase), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, errors.New("unknown error has occurred")
}
