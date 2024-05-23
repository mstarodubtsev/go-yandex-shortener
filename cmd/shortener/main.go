package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
)

var m map[string]string = make(map[string]string)

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
		log.Println("Url received and put to map: ", bodyString)
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("http://localhost:8080/" + hash))
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

func getURL(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if req.Method == "GET" && len(parts) > 1 {
		// return 307 status and Location header
		id := parts[1]
		log.Println("Url shortcut: ", id)
		res.Header().Set("Location", "https://practicum.yandex.ru/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc(`/`, postURL)
	http.HandleFunc(`/{id}`, getURL)

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
	return hashString
}
