package utils

import (
	"net/http"
	"time"
)

// NewHTTPClient creates a new HTTP client with custom configuration, to avoid usin default http.Client
func NewHTTPClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 10, // Set a timeout of 10 seconds, enough for the metadata
	}

	return client
}
