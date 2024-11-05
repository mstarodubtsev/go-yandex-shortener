package app

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"io"
	"log"
	"net/http"
	"strings"
)

// Map to store urls
var m map[string]string = make(map[string]string)

// PostURLHandler Handle POST requests
func PostURLHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" && req.Body != nil {
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
		if _, ok := m[hash]; ok {
			log.Printf(
				"Url already exists in the map: url=%s; hash=%s\n",
				bodyString,
				hash,
			)
		} else {
			// Add new URL to the map
			m[hash] = bodyString
			log.Printf(
				"Url received and added to the map: url=%s; hash=%s\n",
				bodyString,
				hash,
			)
		}
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://" + config.FlagResultURL + "/" + hash))
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

// GetURLHandler Handle GET requests
func GetURLHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if req.Method == "GET" && len(parts) > 1 && len(parts[1]) > 0 {
		// handle redirect request
		id := parts[1]
		log.Println("Get Url shortcut: ", id)
		// return 404 if id not found
		url, ok := m[id]
		if !ok {
			log.Println("Url not found")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		// return 307 status and Location header
		log.Println("Url found: ", url)
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
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
			for k, v := range m {
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
