package models

type ChallengeAnswer struct {
	ID             string  `json:"id" binding:"required"`
	FirstTaskHash  string  `json:"totalHash1" binding:"required"`
	DiffTaskHash   *string `json:"diffHash"`
	SecondTaskHash string  `json:"totalHash2" binding:"required"`
}
