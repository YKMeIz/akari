package akari

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	Name  string
	Token string
}

var (
	db                 *gorm.DB
	databaseConnection bool
)

// initialize sqlite3 database
func InitDatabase(databasePath string) {
	var err error
	db, err = gorm.Open("sqlite3", databasePath)
	if err != nil {
		panic("err: failed to connect database.")
	}
	db.CreateTable(&User{})
	databaseConnection = true
}

func (c Core) OpenDatabase() {
	var err error
	db, err = gorm.Open("sqlite3", c.DatabasePath)
	if err != nil {
		panic("err: failed to connect database.")
	}
	databaseConnection = true
}

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

func IsDatabaseConnected() bool {
	return databaseConnection
}

func checkDatabaseConnection() {
	if !databaseConnection {
		panic("Database is not initialized or opened.")
	}
}
