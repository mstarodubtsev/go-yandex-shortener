package main

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"log"
	"net/http"
)

// Main function
func main() {
	// parse command line arguments
	config.ParseFlags()

	// start server
	log.Printf("Server started at: %s\n", config.FlagRunAddr)
	err := http.ListenAndServe(config.FlagRunAddr, app.Router())
	if err != nil {
		panic(err)
	}
}
