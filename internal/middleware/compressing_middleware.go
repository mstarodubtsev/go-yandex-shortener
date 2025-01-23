package middleware

import (
	"compress/gzip"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"io"
	"net/http"
	"strings"
)

// gzipWriter is a wrapper around http.ResponseWriter that provides gzip compression
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// WithCompressingPost is a middleware that compresses the request body if it is gzipped
func WithCompressingPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			log.Infof("Request is gzipped")
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			log.Infof("Request is not gzipped")
			reader = r.Body
		}
		// Create a new request with the decompressed body
		r.Body = io.NopCloser(reader)
		// Pass the handler the request with potentially decompressed body
		next.ServeHTTP(w, r)
	})
}

// WithCompressingGet is a middleware that compresses the response body if the client supports gzip
func WithCompressingGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that the client supports gzip compression
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// if gzip is not supported, pass control
			log.Infof("Gzip is not supported")
			next.ServeHTTP(w, r)
			return
		}
		// create gzip.Writer over the current w
		log.Infof("Gzip is supported")
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// pass to handler the variable of type gzipWriter for output data
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
