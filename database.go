package akari

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Device struct {
	ID    int
	Name  string
	Token string
}

var db *gorm.DB

// initialize sqlite3 database
func initDatabase(databasePath string) {
	var err error
	db, err = gorm.Open("sqlite3", databasePath)
	if err != nil {
		panic("err: failed to connect database.")
	}
	// db.CreateTable(&Device{})
	// new := Device{Name: "test", Token: "c8424762ae5148954e2e640c389c31c5"}
	// db.Save(&new)
}

// IsName checks if given name appears in database.
//
// It returns true if name appears in the database; returns false
// if name does not appear in the database
func isName(name string) bool {
	d := Device{Name: name}
	db.Where(&d).First(&d)
	if d.ID != 0 {
		return true
	}
	return false
}

// IsName checks if given token appears in database.
//
// It returns true if token appears in the database; returns false
// if token does not appear in the database
func isToken(token string) bool {
	d := Device{Token: token}
	db.Where(&d).First(&d)
	if d.ID != 0 {
		return true
	}
	return false
}

// CompareToken compares given token with the token stored in database.
//
// It returns true if both tokens are same; returns false
// if tokens are different.
func compareToken(name, token string) bool {
	d := Device{Name: name}
	db.Where(&d).First(&d)
	if d.ID != 0 && token == d.Token {
		return true
	}
	return false
}

func getName(token string) string {
	d := Device{Token: token}
	db.Where(&d).First(&d)
	if d.ID != 0 {
		return d.Name
	}
	return ""
}
