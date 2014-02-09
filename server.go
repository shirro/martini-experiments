package main

import (
	//	"github.com/codegangsta/martini"
	"net/http"
	"shirro.com/martini"
)

func AuthController(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-API-KEY") != "secret123" {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func BasicNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {

	m := martini.Classic()
	cors := &Cors{Headers: StandardHeaders}
	/*
		cors.Origins = map[string]struct{}{
			"http://127.0.0.1": struct{}{},
		}
	// */
	m.Use(AuthController)
	m.Use(cors.Middleware)
	m.NotFound(cors.NotFound, BasicNotFound)

	m.Get("/hello/:name", func(params martini.Params) string {
		return "Hello " + params["name"]
	})

	m.Post("/hello/:name", func(params martini.Params) string {
		return "Posting"
	})

	m.Run()
}
