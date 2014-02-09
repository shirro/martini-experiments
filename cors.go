package main

import (
	"net/http"
)

type Cors struct {
	OriginsMutex sync.RWMutex
	Origins      map[string]struct{}
	HeadersMutex sync.RWMutex
	Headers      map[string]string
}

var StandardHeaders = map[string]string{
	"Access-Control-Max-Age":        "86000",
	"Access-Control-Allow-Headers":  "Content-Type, Origin, Authorization",
	"Access-Control-Expose-Headers": "Content-Length",
}

func (cors *Cors) Middleware(w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")
	h := w.Header()

	if !cors.setOrigin(h, origin) {
		w.WriteHeader(http.StatusForbidden)
	}

}

func (cors *Cors) setOrigin(h http.Header, origin string) bool {

	// Block empty or nonexistent Origin headers
	if origin == "" {
		return false
	}

	// Reader lock so we can change the map dynamically
	OriginsMutex.RLock()
	defer OriginsMutex.RUnlock()

	// Empty Origins map allows all domains
	if len(cors.Origins) == 0 {
		h.Set("Access-Control-Allow-Origin", "*")
		return true
	}

	// Allow request if Origin in map
	if _, ok := cors.Origins[origin]; ok {
		h.Set("Access-Control-Allow-Origin", origin)
		return true
	}

	// Default
	return false
}

func (cors *Cors) setHeaders(h http.Header) {
	// Reader lock so we can change the headers dynamically
	cors.HeadersMutex.RLock()
	defer cors.HeadersMutex.RUnlock()
	for header, value := range cors.Headers {
		h.Set(header, value)
	}
}

func (cors *Cors) Handler(methods string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		if r.Method == "OPTIONS" {
			cors.setHeaders(w)
		} else {
			h.Set("Allow", methods)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
