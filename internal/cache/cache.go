package cache

import "time"

type Challenge struct {
	Task         string
	ExpectedHash string
	ExpiresAt    time.Time
}

type Cache interface {
	Set(key string, ch Challenge) error
	Get(key string) (Challenge, bool, error)
	Delete(key string) error
}
