package server

import (
	"net/http"
)

func newRouter(apiMux *http.ServeMux) *http.ServeMux {
	rootMux := http.NewServeMux()

	rootMux.Handle("/api/", RequestIDMiddleware(apiMux))

	opsMux := http.NewServeMux()
	opsMux.HandleFunc("/health", HealthCheckHandler)
	rootMux.Handle("/", opsMux)

	return rootMux
}
