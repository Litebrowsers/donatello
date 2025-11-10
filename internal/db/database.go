// Package db provides database connectivity and initialization.
package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

// DB is the global database connection.
var DB *gorm.DB

// InitDB initializes the database connection.
func InitDB() error {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "donatello.db"
	}
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
