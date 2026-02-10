package storage

import "time"

type APIKey struct {
	ID          string
	KeyHash     []byte
	Producer    string
	Environment string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	RateLimit   int // events/sec
}

var APIKeys map[string]*APIKey
