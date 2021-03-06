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
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User defines user information.
type User struct {
	// Name of a user.
	Name string

	// Token is usually generated by system.
	Token string
}

var (
	db                 *gorm.DB
	databaseConnection bool
)

// InitDatabase Initializes a new SQLite database file.
// It is only utilized for first time initialization.
// If you have already initialized a database file,
// please run OpenDatabase() instead.
func InitDatabase(databasePath string) {
	var err error
	db, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		panic("err: failed to connect database.")
	}
	db.AutoMigrate(&User{})
	databaseConnection = true
}

// OpenDatabase opens an exist SQLite database file.
func (c Core) OpenDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open(c.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("err: failed to connect database.")
	}
	databaseConnection = true
}

// IsUser checks if given user information is exist in database/system.
// It accepts one of user information, or both "Name" and "Token".
// Return true for exist; false for not exist.
func (u User) IsUser() bool {
	checkDatabaseConnection()
	if u.Name != "" && u.Token != "" {
		r := User{Name: u.Name}
		db.Where(&r).First(&r)
		if u.Token == r.Token {
			return true
		}
		return false
	}
	if u.Name != "" {
		r := User{Name: u.Name}
		db.Where(&r).First(&r)
		if r.Token != "" {
			return true
		}
		return false
	}
	if u.Token != "" {
		r := User{Token: u.Token}
		db.Where(&r).First(&r)
		if r.Name != "" {
			return true
		}
		return false
	}
	return false
}

// UserCompletion complets user information.
// Give one of user information, it will fill the missing value.
// Return error if giving both "Name" and "Token", or user is not found.
func (u *User) UserCompletion() error {
	checkDatabaseConnection()
	if u.Name != "" && u.Token != "" {
		return errors.New("err: Given user information is already completed.")
	}
	db.Where(&u).First(&u)
	if u.Name != "" && u.Token != "" {
		return nil
	}
	return errors.New("err: User is not found.")
}

// IsDatabaseConnected checks if database is opened.
func IsDatabaseConnected() bool {
	return databaseConnection
}

func checkDatabaseConnection() {
	if !databaseConnection {
		panic("Database is not initialized or opened.")
	}
}
