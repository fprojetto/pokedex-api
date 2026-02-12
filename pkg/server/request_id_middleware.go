package server

import (
	"context"
	"crypto/rand"
	"net/http"
)

// Define a custom type for context keys to prevent collisions
type contextKey string

const requestIDKey contextKey = "request_id"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Generate ID
		id := rand.Text()

		// 2. Set in Response Header (Standard practice)
		w.Header().Set("X-Request-ID", id)

		// 3. Store in Context
		ctx := context.WithValue(r.Context(), requestIDKey, id)

		// 4. Pass the modified context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to retrieve the ID later
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
