package client

import (
	"net"
	"net/http"
	"time"
)

func HttpClient() *http.Client {
	// Define the Transport (Network Layer)
	t := &http.Transport{
		// 1. Connection Dialing settings
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // Max time to establish a TCP connection
			KeepAlive: 30 * time.Second, // Probe interval for active connections
		}).DialContext,

		// 2. TLS/SSL Handshake
		TLSHandshakeTimeout: 10 * time.Second, // Max time for HTTPS handshake

		// 3. Connection Pooling (Very Important!)
		MaxIdleConns:        100,              // Total max idle connections across all hosts
		MaxIdleConnsPerHost: 100,              // Max idle connections for a SINGLE host (Default is 2!)
		IdleConnTimeout:     90 * time.Second, // How long an idle connection is kept open

		// 4. Response Reading
		ResponseHeaderTimeout: 10 * time.Second, // Max time to wait for server response headers
	}

	// Define the Client (High Level)
	client := &http.Client{
		Transport: t,

		// 5. Total Request Timeout
		// This includes Dial, TLS, sending request, waiting for response, and reading body.
		Timeout: 30 * time.Second,
	}

	return client
}

func BoolPtr(b bool) *bool {
	return &b
}
