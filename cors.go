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
}

var StandardHeaders = map[string]string{
	"Access-Control-Max-Age":        "86000",
	"Access-Control-Allow-Headers":  "Content-Type, Origin, Authorization",
	"Access-Control-Expose-Headers": "Content-Length",
}

func (cors *Cors) Middleware(w http.ResponseWriter, r *http.Request, routes martini.Routes) {

	origin := r.Header.Get("Origin")
	h := w.Header()

	if !cors.setOrigin(h, origin) {
		w.WriteHeader(http.StatusForbidden)
	}

	// Add cors headers if this is an OPTIONS request
	if r.Method == "OPTIONS" {
		methods := routes.MethodsFor(r.URL.Path)
		cors.setHeaders(h)
		h.Set("Access-Control-Allow-Methods", stringMethods(methods))
		w.WriteHeader(http.StatusOK)
		return
	}

}

func stringMethods(methods []string) string {
	methods = append(methods, "OPTIONS")
	return strings.Join(methods, ",")
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
