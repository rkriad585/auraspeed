package http

import (
	"net/http"
	"time"
)

// Client is a shared HTTP client with connection pooling for future API calls.
// Currently the app doesn't make many HTTP calls, but as features are added
// (e.g., additional speed test servers, API integrations, telemetry),
// this pooled client will improve performance by reusing connections.
var Client *http.Client

func init() {
	Client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}
