package main

import (
	"log"
	"net/http"
	"strings"
)

func postUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" && req.Body != nil {
		// log request and return 200 and message
		log.Println("Url received: ", req.Body)
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("http://localhost:8080/EwHXdJfB"))
	} else {
		// return 400
		res.WriteHeader(http.StatusBadRequest)
	}
}

func getUrl(res http.ResponseWriter, req *http.Request) {
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
	http.HandleFunc(`/`, postUrl)
	http.HandleFunc(`/{id}`, getUrl)

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
