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

// WithCompressing is a middleware that compresses the request body if Content-Encoding == gzip
// and compresses the response body if Accept-Encoding == gzip
func WithCompressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that the request body is gzipped Content-Encoding == gzip
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			// create gzip.Reader over the current r.Body
			log.Infof("Request is gzipped")
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(gz)
			defer gz.Close()
		} else {
			log.Infof("Request is not gzipped")
		}
		// check that the client supports gzip compression Accept-Encoding == gzip
		var responseWriter http.ResponseWriter
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create gzip.Writer over the current w
			log.Infof("Gzip is supported")
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			responseWriter = gzipWriter{ResponseWriter: w, Writer: gz}
			w.Header().Set("Content-Encoding", "gzip")
		} else {
			log.Infof("Gzip is not supported")
			responseWriter = w
		}
		// pass to handler the variable of type gzipWriter for output data
		next.ServeHTTP(responseWriter, r)
	})
}
