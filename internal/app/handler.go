package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/middleware"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/storage"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Store URL storage
var store storage.Storage

func SetStore(s storage.Storage) {
	store = s
}

// POST structure of the request body
type shortenRequest struct {
	URL string `json:"url"`
}

// POST structure of the response body
type shortenResponse struct {
	Result string `json:"result"`
}

// Router
func Router() chi.Router {
	r := chi.NewRouter()

	// Apply the WithLogging middleware with the logger
	r.Use(middleware.WithLogging, middleware.WithCompressing)

	// Routes
	r.Post("/api/shorten", PostURLHandlerJSON)
	r.Post("/", PostURLHandler)
	r.Get("/{id}", GetURLHandler)
	r.Get("/list", ListURLHandler)
	return r
}

// PostURLHandlerJSON Handle POST requests with JSON body
func PostURLHandlerJSON(res http.ResponseWriter, req *http.Request) {
	log.Infof("POST /api/shorten")
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
	// decode request JSON body
	var request shortenRequest
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("URL received: %s", request.URL)
	url := string(request.URL)
	hash := getHash(url)
	// check if the URL already exists in the map
	_, ok, err := store.GetURL(hash)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok {
		log.Infof("URL already exists in the map: url=%s; hash=%s", url, hash)
	} else {
		// Add new URL to the map
		store.AddURL(hash, url)
		log.Infof("URL received and added to the map: url=%s; hash=%s", url, hash)
	}
	response := shortenResponse{Result: config.Config.BaseURL + "/" + hash}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Unable to marshal response", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(responseBytes)
}

// Validate checks if the required fields are present
func (r *shortenRequest) Validate() error {
	return ValidateURL(r.URL)
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
	if err := ValidateURL(bodyString); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	hash := getHash(bodyString)
	// check if the URL already exists in the map
	_, ok, err := store.GetURL(hash)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok {
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
	url, ok, err := store.GetURL(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
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
			rows, err := store.GetAll()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			for k, v := range rows {
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

// ValidateURL Check if the URL is valid
func ValidateURL(value string) error {
	log.Infof("Validating URL: %s", value)
	// check if URL is empty
	if value == "" {
		return errors.New("URL cannot be empty")
	}
	// parse the URL
	parsedURL, err := url.Parse(value)
	if err != nil {
		log.Infof("Invalid URL format: %s", value)
		return fmt.Errorf("invalid URL format: %v", err)
	}
	// check if scheme is HTTP or HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must start with http:// or https://")
	}
	// check if host is present
	if parsedURL.Host == "" {
		return errors.New("URL must contain a host")
	}
	log.Infof("Valid URL: %s", value)
	return nil
}
