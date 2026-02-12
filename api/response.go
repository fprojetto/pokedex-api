package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fprojetto/pokedex-api/pkg/server"
)

// Envelope is the generic wrapper for all responses
type Envelope struct {
	Data  any    `json:"data,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

const (
	ErrCodeInternal = "INTERNAL_ERROR"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, data any, status int) {
	requestID := server.GetRequestID(r.Context())
	resp := Envelope{
		Data: data,
		Meta: &Meta{RequestID: requestID},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("failed to write json: %v", err)
	}
}

// WriteError sends a structured error response
func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	requestID := server.GetRequestID(r.Context())
	resp := Envelope{
		Error: &Error{
			Code:    code,
			Message: message,
		},
		Meta: &Meta{RequestID: requestID},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("failed to write json: %v", err)
	}
}
