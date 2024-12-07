package main

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"net/http"
)

// Main function
func main() {
	// initialize logger
	log.InitializeLogger()
	defer log.Logger.Sync()

	// parse config
	config.ParseConfig()

	// start server
	log.Logger.Infof("Server started at: %s", config.Config.ServerAddress)
	err := http.ListenAndServe(config.Config.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}
}
