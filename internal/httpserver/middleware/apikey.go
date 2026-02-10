package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"event-ingestion-service/internal/storage"
	"log"
	"net/http"
	"strings"
)

func APIKeyMiddleware(next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiToken, err := bearerToken(r)
		if err != nil {
			log.Println("Could not authenticate properly")
			http.Error(w, "Could not authenticate", http.StatusUnauthorized)
			return
		}

		found := false

		for _, key := range storage.APIKeys {
			if subtle.ConstantTimeCompare(key.KeyHash, apiToken) == 1 {
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Not a valid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerToken(r *http.Request) ([]byte, error) {
	rawToken := r.Header.Get("Authorization")
	pieces := strings.SplitN(rawToken, " ", 2)
	if len(pieces) < 2 {
		return []byte{}, errors.New("token with incorrect bearer format")
	}
	token := strings.TrimSpace(pieces[1])
	h := sha256.Sum256([]byte(token))
	hash := h[:]
	return hash, nil
}
