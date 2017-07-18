package main

import (
	"log"
	"net/http"
)

import "github.com/gorilla/mux"

// Router for http api calls
func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		// log.Println("Attaching HTTP handler for route:", route.Method, route.Pattern)
		log.Println("Attaching HTTP handler for route: ", route.Method, " ", route.Pattern)
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
