package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
)

// Map to store urls
var m map[string]string = make(map[string]string)

// Handle POST requests
func postURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" && req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Unable to read body", http.StatusBadRequest)
			return
		}
		bodyString := string(bodyBytes)
		hash := getHash(bodyString)
		m[hash] = bodyString
		log.Printf(
			"Url received and put to map: url=%s; hash=%s\n",
			bodyString,
			hash,
		)
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080/" + hash))
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

// Handle GET requests
func getURL(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if req.Method == "GET" && len(parts) > 1 {
		// return 307 status and Location header
		id := parts[1]
		log.Println("Url shortcut: ", id)

		// return 404 if id not found
		url, ok := m[id]
		if !ok {
			log.Println("Url not found")
			res.WriteHeader(http.StatusNotFound)
			return
		}

		// return result
		log.Println("Url found: ", url)
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

// Main function
func main() {
	http.HandleFunc(`/`, postURL)
	http.HandleFunc(`/{id}`, getURL)

	log.Println("Server started at :8080")
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
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
