package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	ID            string `gorm:"primaryKey"`
	Task          string
	ActualHash    string
	ExpectedHash  string
	ExpiresAt     time.Time
	NoiseDetected bool
	Fingerprint   string
}
