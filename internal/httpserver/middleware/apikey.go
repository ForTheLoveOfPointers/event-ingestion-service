package middleware

import (
	"errors"
	"event-ingestion-service/internal/storage"
	"log"
	"net/http"
	"strings"
)

func APIKeyMiddleware(next http.Handler) http.HandlerFunc {

	reverseKeyIndex := make(map[string]string)
	for name, key := range storage.APIKeys {
		reverseKeyIndex[key.Key] = name
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiToken, err := bearerToken(r)
		if err != nil {
			log.Println("Could not authenticate properly")
			http.Error(w, "Could not authenticate", http.StatusUnauthorized)
			return
		}

		if _, found := storage.APIKeys[apiToken]; !found {
			http.Error(w, "Not a valid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerToken(r *http.Request) (string, error) {
	rawToken := r.Header.Get("Authorization")
	pieces := strings.SplitN(rawToken, " ", 2)

	if len(pieces) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}

	token := strings.TrimSpace(pieces[1])

	return token, nil
}
