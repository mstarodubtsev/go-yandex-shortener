package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
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
	log.Println("Server started at :8080")
	err := http.ListenAndServe(`:8080`, router())
	if err != nil {
		panic(err)
	}
}
