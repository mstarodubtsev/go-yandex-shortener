package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"log"
	"net/http"
)

// Router
func router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", app.PostURLHandler)
	r.Get("/list", app.ListURLHandler)
	r.Get("/{id}", app.GetURLHandler)
	return r
}

// Main function
func main() {
	// parse command line arguments
	config.ParseFlags()

	// start server
	log.Printf("Server started at: %s\n", config.FlagRunAddr)
	err := http.ListenAndServe(config.FlagRunAddr, router())
	if err != nil {
		panic(err)
	}
}
