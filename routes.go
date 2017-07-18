package main

import (
	"net/http"
)

type apiRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type apiRoutes []apiRoute

var routes = apiRoutes{

	// Health check
	apiRoute{"Ping", "GET", "/ping", PingHandler},

	apiRoute{"GetQueues", "GET", "/api/v1/queues", GetQueuesHandler},

	apiRoute{"AddMessage", "POST", "/api/v1/queue/{q}", AddMessageHandler},
	apiRoute{"GetMessage", "GET", "/api/v1/queue/{q}", GetMessageHandler},

	// Web Socket apiRoute
	// apiRoute{"Socket", "GET", "/ws/{ds}", serveWs},
}
