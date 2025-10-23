// Package db provides database connectivity and initialization.
package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global database connection.
var DB *gorm.DB

// InitDB initializes the database connection.
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("donatello.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
