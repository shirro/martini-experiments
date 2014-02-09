package main

import (
	"net/http"
	"shirro.com/martini"
)

func MethodNotFound(w http.ResponseWriter, r *http.Request, routes martini.Routes) {
	// For if we have a patched Martini and can find all methods for path
	methods := routes.MethodsFor(r.URL.Path)

	// No methods is a 404 - don't handle here
	if len(methods) == 0 {
		return
	}

	h := w.Header()

	// If we have a path with some methods, and our method isn't in it
	// return a 405 with Allow header
	h.Set("Allow", stringMethods(methods))
	w.WriteHeader(http.StatusMethodNotAllowed)

}

func BasicNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
