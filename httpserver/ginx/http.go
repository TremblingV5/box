package ginx

import (
	"net/http"
	"time"
)

func defaultHttpServer() *http.Server {
	return &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Minute,
	}
}
