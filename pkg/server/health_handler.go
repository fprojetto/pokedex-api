package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		log.Printf("failed to write json: %v", err)
	}
}
