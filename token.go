// Copyright © 2016 nrechn <nrechn@gmail.com>
//
// This file is part of akari.
//
// akari is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// akari is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with akari. If not, see <http://www.gnu.org/licenses/>.
//

package akari

import (
	"crypto/rand"
	"errors"
	"fmt"
)

func randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// GenerateToken returns a new token.
func GenerateToken() string {
	return randToken()
}

// RegisterUser registers a new user to database.
// It only accepts username with empty "Token".
// Return username, user's new token, error.
func (u User) RegisterUser() (string, string, error) {
	checkDatabaseConnection()
	if u.Token != "" {
		return "", "", errors.New("err: Custom token is not allowed. Token should be generated by system.")
	}
	if u.Name == "" {
		return "", "", errors.New("err: User name cannot be empty.")
	}
	new := User{Name: u.Name, Token: randToken()}
	db.Save(&new)
	return new.Name, new.Token, nil
}

// RevokeUser revokes a user from database.
func (u User) RevokeUser() error {
	checkDatabaseConnection()
	if !u.IsUser() {
		return errors.New("err: No such a user.")
	}
	if u.Name != "" {
		db.Where(&u).First(&u).Delete(&u)
		return nil
	}
	if u.Token != "" {
		db.Where(&u).First(&u).Delete(&u)
		return nil
	}
	return errors.New("err: Name or Token is required to revoke a user.")
}
