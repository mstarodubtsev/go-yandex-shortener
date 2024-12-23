package app

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-chi/chi/v5"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
)

// POST структура тела запроса
type shortenRequest struct {
	URL string `json:"url"`
}

// POST структура тела ответа
type shortenResponse struct {
	Result string `json:"result"`
}

// Store URL storage
var store storage.Storage = storage.NewMap()

// Router
func Router() chi.Router {
	r := chi.NewRouter()
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
		log.Printf(
			"Url already exists in the map: url=%s; hash=%s\n",
			bodyString,
			hash,
		)
	} else {
		// Add new URL to the map
		store.AddURL(hash, bodyString)
		log.Printf(
			"Url received and added to the map: url=%s; hash=%s\n",
			bodyString,
			hash,
		)
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
	log.Println("Get Url shortcut: ", id)
	// return 404 if id not found
	url, ok := store.GetURL(id)
	if !ok {
		log.Println("Url not found")
		res.WriteHeader(http.StatusNotFound)
		return
	}
	// return 307 status and Location header
	log.Println("Url found: ", url)
	//res.Header().Set("Location", url)
	//res.WriteHeader(http.StatusTemporaryRedirect)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// ListURLHandler Handle list URL requests
func ListURLHandler(res http.ResponseWriter, req *http.Request) {
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
