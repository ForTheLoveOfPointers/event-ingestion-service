package handlers

import (
	"encoding/json"
	"event-ingestion-service/internal/httpserver/types"
	"fmt"
	"io"
	"net/http"
)

func EventHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var e types.Event
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading request body: %s\r\n", err.Error())
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
		if err := json.Unmarshal(body, &e); err != nil {
			fmt.Printf("Unable to process JSON data: %s\r\n", err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
		}

	})
}
