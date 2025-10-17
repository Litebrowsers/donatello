package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("donatello.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
