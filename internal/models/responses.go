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

// ChallengeAnswer represents the answer to a challenge that is sent from the client.
type ChallengeAnswer struct {
	ID                string  `json:"id" binding:"required"`
	FirstTaskHash     string  `json:"totalHash1" binding:"required"`
	DiffTaskHash      *string `json:"diffHash"`
	SecondTaskHash    string  `json:"totalHash2" binding:"required"`
	SecondTaskMetrics string  `json:"metrics2" binding:"required"`
	CopyMismatch      *bool   `json:"copyMismatch"`
}
