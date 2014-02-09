package main

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func AuthController(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-API-KEY") != "secret123" {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func main() {
	cors := &Cors{Headers: StandardHeaders}
	m := martini.Classic()
	m.Use(AuthController)
	m.Use(cors.Middleware)

	m.Get("/hello/:name", func(params martini.Params) string {
		return "Hello " + params["name"]
	})
	m.Any("/hello/**", cors.Handler("GET,OPTIONS"))

	m.Run()
}
