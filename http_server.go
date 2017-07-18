package main

import (
	"fmt"
	"log"
	"net/http"
)

const DEFAULT_HTTP_PORT = 8080

type HttpServer struct {
	Port int
}

func (self HttpServer) Start() {
	// Attach Http Hanlders
	router := Router()

	// Start server
	log.Println("Magic happens on port ", self.Port)

	bind := fmt.Sprintf(":%v", self.Port)

	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}

}
