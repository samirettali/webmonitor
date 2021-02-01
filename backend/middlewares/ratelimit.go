package middlewares

import "net/http"

// Usage:
// http.HandleFunc("/route", limitNumClients(handler, 10))
func limitNumClients(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	sem := make(chan struct{}, maxClients)

	return func(w http.ResponseWriter, req *http.Request) {
		sem <- struct{}{}
		defer func() { <-sem }()
		f(w, req)
	}
}
