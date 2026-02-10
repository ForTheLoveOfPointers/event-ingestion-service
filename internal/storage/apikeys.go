package storage

import "time"

type APIKey struct {
	Key       string
	Owner     string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

var APIKeys map[string]APIKey
