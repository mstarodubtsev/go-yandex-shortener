package main

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
	"log"
	"net/http"
)

// Main function
func main() {
	http.HandleFunc(`/`, app.PostURLHandler)
	http.HandleFunc(`/list`, app.ListURLHandler)
	http.HandleFunc(`/{id}`, app.GetURLHandler)

	log.Println("Server started at :8080")
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
