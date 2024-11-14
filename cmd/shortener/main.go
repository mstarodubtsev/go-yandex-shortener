package main

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"log"
	"net/http"
)

// Main function
func main() {
	// parse config
	config.ParseConfig()

	// start server
	log.Printf("Server started at: %s\n", config.Config.ServerAddress)
	err := http.ListenAndServe(config.Config.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}
}
