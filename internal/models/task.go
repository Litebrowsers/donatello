package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Value string
	Name  string
}
