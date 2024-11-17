package app

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-chi/chi/v5"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/middleware"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/storage"
	"io"
	"net/http"
	"strings"
)

// POST structure of the request body
type shortenRequest struct {
	URL string `json:"url"`
}

// POST structure of the response body
type shortenResponse struct {
	Result string `json:"result"`
}

// Store URL storage
var store storage.Storage = storage.NewMap()

// Router
func Router() chi.Router {
	r := chi.NewRouter()

	// Apply the WithLogging middleware with the logger
	r.Use(func(next http.Handler) http.Handler {
		return middleware.WithLogging(next)
	})

	r.Post("/", PostURLHandler)
	r.Get("/{id}", GetURLHandler)
	r.Get("/list", ListURLHandler)
	return r
}

// PostURLHandler Handle POST requests
func PostURLHandler(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(res, "Empty body", http.StatusBadRequest)
		return
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Unable to read body", http.StatusBadRequest)
		return
	}
	if len(bodyBytes) == 0 {
		http.Error(res, "Empty body", http.StatusBadRequest)
		return
	}
	bodyString := string(bodyBytes)
	hash := getHash(bodyString)
	// check if the URL already exists in the map
	if _, ok := store.GetURL(hash); ok {
		log.Infof("URL already exists in the map: url=%s; hash=%s", bodyString, hash)
	} else {
		// Add new URL to the map
		store.AddURL(hash, bodyString)
		log.Infof("URL received and added to the map: url=%s; hash=%s", bodyString, hash)
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Config.BaseURL + "/" + hash))
}

// GetURLHandler Handle GET requests
func GetURLHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if !(len(parts) > 1 && len(parts[1]) > 0) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// handle redirect request
	id := parts[1]
	log.Infof("Get Url shortcut: %s", id)
	// return 404 if id not found
	url, ok := store.GetURL(id)
	if !ok {
		log.Infof("Url not found: %s", id)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	// return 307 status and Location header
	log.Infof("Url found: %s", url)
	//res.Header().Set("Location", url)
	//res.WriteHeader(http.StatusTemporaryRedirect)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// ListURLHandler Handle list URL requests
func ListURLHandler(res http.ResponseWriter, req *http.Request) {
	log.Infof("List Url shortcuts")
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if req.Method == "GET" && len(parts) > 1 && len(parts[1]) > 0 {
		// handle list request
		id := parts[1]
		if id == "list" {
			res.Header().Set("Content-Type", "text/plain")
			res.WriteHeader(http.StatusOK)
			for k, v := range store.GetAll() {
				res.Write([]byte(k + " -> " + v + "\n"))
			}
			return
		}
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

// Compute SHA-256 hash of the body string
func getHash(bodyString string) string {
	hash := sha256.New()
	hash.Write([]byte(bodyString))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	if len(hashString) > 8 {
		hashString = hashString[:8]
	}
	return hashString
}
