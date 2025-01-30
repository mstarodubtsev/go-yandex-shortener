package main

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/app"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/storage"
	"net/http"
)

// Main function
func main() {
	// initialize logger
	log.InitializeLogger()
	defer log.Logger.Sync()

	// parse config
	config.ParseConfig()

	// init storage
	store, err := storage.NewFileStorage(config.Config.FileStoragePath)
	if err != nil {
		panic(err)
	}
	app.SetStore(store)

	// start server
	log.Infof("Server started at: %s", config.Config.ServerAddress)
	err = http.ListenAndServe(config.Config.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}
}
