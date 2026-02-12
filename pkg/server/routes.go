package server

import (
	"net/http"
)

func newRouter(apiMux http.Handler) *http.ServeMux {
	rootMux := http.NewServeMux()

	rootMux.Handle("/api/", apiMux)

	opsMux := http.NewServeMux()
	opsMux.HandleFunc("/health", HealthCheckHandler)
	rootMux.Handle("/", opsMux)

	return rootMux
}
