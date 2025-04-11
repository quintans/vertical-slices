package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/quintans/vertical-slices/internal/config"
)

func main() {
	// Configure the API routes
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	// Configure the application
	c := &config.Config{}
	config.WireInfra(c)
	config.WireRepositories(c)
	config.WireProductEventHandlers(c)
	config.WireProductAPI(c, api)
	config.WireOrderAPI(c, api)

	// Start the server!
	http.ListenAndServe("127.0.0.1:8888", router)
}
