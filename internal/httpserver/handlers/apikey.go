package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"event-ingestion-service/internal/storage"
	"net/http"
	"time"
)

type apiKeyRequest struct {
	Producer    string `json:"producer"`
	Environment string `json:"environment"`
	RateLimit   int    `json:"rate_limit"`
}

func APIKeyGen() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req apiKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Unable to parse request body", http.StatusBadRequest)
			return
		}

		raw, hash, err := generateAPIKey()
		if err != nil {
			http.Error(w, "unable to generate api key", http.StatusInternalServerError)
			return
		}
		apiKey := storage.APIKey{
			KeyHash:     hash,
			Producer:    req.Producer,
			Environment: req.Environment,
			RateLimit:   req.RateLimit,
			CreatedAt:   time.Now(),
		}

		storage.APIKeys[string(hash)] = &apiKey

		json.NewEncoder(w).Encode(map[string]string{
			"api_key": raw,
		})

	})
}

func generateAPIKey() (raw string, hash []byte, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}

	raw = "ingest_prod_ak_" + base64.RawURLEncoding.EncodeToString(b)

	h := sha256.Sum256([]byte(raw))
	hash = h[:]

	return
}
