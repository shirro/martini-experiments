package main

import (
	//	"github.com/codegangsta/martini"
	"net/http"
	"shirro.com/martini"
	"strings"
)

func AuthController(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-API-KEY") != "secret123" {
		w.WriteHeader(http.StatusUnauthorized)
	}
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
	m.NotFound(MethodNotFound, BasicNotFound)

	m.Get("/hello/:name", func(params martini.Params, route martini.Routes, r *http.Request) string {
		return "Hello " + params["name"] + strings.Join(route.MethodsFor(r.URL.Path), ",")
	})

	m.Post("/hello/:name", func(params martini.Params) string {
		return "Posting"
	})

	m.Run()
}
