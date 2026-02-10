package handlers

import (
	"encoding/json"
	"event-ingestion-service/internal/httpserver/types"
	"event-ingestion-service/internal/ingest"
	"fmt"
	"io"
	"net/http"
)

func EventHandler(ingestor *ingest.Ingest) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method supported", http.StatusMethodNotAllowed)
			return
		}
		var e types.Event
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Printf("Error reading request body: %s\r\n", err.Error())
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &e); err != nil {
			fmt.Printf("Unable to process JSON data: %s\r\n", err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if res := ingestor.Enqueue(e); !res {
			w.Header().Add("Retry-After", "60")
			http.Error(w, "Event limit reached in queue.", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
