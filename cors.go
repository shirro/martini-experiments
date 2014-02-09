package main

import (
	"net/http"
	"shirro.com/martini"
	"strings"
	"sync"
)

type Cors struct {
	OriginsMutex sync.RWMutex
	Origins      map[string]struct{}
	HeadersMutex sync.RWMutex
	Headers      map[string]string
	Martini      martini.Router
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

func stringMethods(methods []string) string {
	methods = append(methods, "OPTIONS")
	return strings.Join(methods, ",")
}

func (cors *Cors) NotFound(w http.ResponseWriter, r *http.Request) {
	// For if we have a patched Martini and can find all methods for path
	if cors.Martini != nil {
		methods := cors.Martini.Methods(r.URL.Path)

		// No methods is a 404 - don't handle here
		if len(methods) == 0 {
			return
		}

		h := w.Header()

		// Add cors headers if this is an OPTIONS request
		if r.Method == "OPTIONS" {
			cors.setHeaders(h)
			h.Set("Access-Control-Allow-Methods", stringMethods(methods))
			w.WriteHeader(http.StatusOK)
			return
		}

		// If we have a path with some methods, and our method isn't in it
		// and isn't OPTIONS return a 405 with Allow header
		h.Set("Allow", stringMethods(methods))
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func (cors *Cors) setOrigin(h http.Header, origin string) bool {

	// Block empty or nonexistent Origin headers
	if origin == "" {
		return false
	}

	// Reader lock so we can change the map dynamically
	cors.OriginsMutex.RLock()
	defer cors.OriginsMutex.RUnlock()

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

// Non-DRY version is to use a fallback All() on a route
func (cors *Cors) Handler(methods string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		if r.Method == "OPTIONS" {
			cors.setHeaders(h)
			h.Set("Access-Control-Allow-Methods", methods)
		} else {
			h.Set("Allow", methods)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
