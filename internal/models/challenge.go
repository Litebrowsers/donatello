/*
# Donatello

Copyright © 2025 Litebrowsers
Licensed under a Proprietary License

This software is the confidential and proprietary information of Litebrowsers
Unauthorized copying, redistribution, or use is prohibited.
For licensing inquiries, contact:
vera cohopie at gmail dot com
thor betson at gmail dot com
*/

package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	ID             string `gorm:"primaryKey"`
	Task           string
	ActualHash     string
	ExpectedHash   string
	ExpiresAt      time.Time
	NoiseDetected  bool
	Fingerprint    string
	NoiseHash      *string
	ProcessingTime int64
	JavaScript     *bool `gorm:"default:null"`
}
