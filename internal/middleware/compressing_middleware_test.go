package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setup function to initialize common test data
func setup() {
	log.InitializeLogger()
	defer log.Logger.Sync()
}

// mockHandler is a test handler that writes a simple response
var payload = []byte("Hello, World!")

func mockHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// TestWithCompressing tests both request decompression and response compression
func TestWithCompressing(t *testing.T) {
	setup()
	testCases := []struct {
		name            string
		method          string
		contentEncoding string
		acceptEncoding  string
		requestBody     []byte
		expectedBody    []byte
	}{
		{
			name:            "POST - Uncompressed Request",
			method:          "POST",
			contentEncoding: "",
			acceptEncoding:  "",
			requestBody:     payload,
			expectedBody:    payload,
		},
		{
			name:            "POST - Gzipped Request",
			method:          "POST",
			contentEncoding: "gzip",
			acceptEncoding:  "",
			requestBody:     compressBody(payload),
			expectedBody:    payload,
		},
		{
			name:            "POST - Gzipped Request and Accept Gzip Response",
			method:          "POST",
			contentEncoding: "gzip",
			acceptEncoding:  "gzip",
			requestBody:     compressBody(payload),
			expectedBody:    payload,
		},
		{
			name:            "GET - Client Accepts Gzip",
			method:          "GET",
			contentEncoding: "",
			acceptEncoding:  "gzip",
			requestBody:     nil,
			expectedBody:    payload,
		},
		{
			name:            "GET - Client Doesn't Accept Gzip",
			method:          "GET",
			contentEncoding: "",
			acceptEncoding:  "",
			requestBody:     nil,
			expectedBody:    payload,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mock handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mockHandler(w, r, tc.expectedBody)
			})

			// wrap handler with compressing middleware
			wrappedHandler := WithCompressing(handler)

			// create request
			req, err := http.NewRequest(tc.method, "/", bytes.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			if tc.contentEncoding != "" {
				req.Header.Set("Content-Encoding", tc.contentEncoding)
			}
			if tc.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tc.acceptEncoding)
			}

			// perform request
			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			// verify response compression and content
			if tc.acceptEncoding == "gzip" {
				if rr.Header().Get("Content-Encoding") != "gzip" {
					t.Errorf("Expected gzip encoding in response, got none")
				}
				decompressed := decompressBody(rr.Body.Bytes())
				if !bytes.Equal(decompressed, tc.expectedBody) {
					t.Errorf("Unexpected decompressed response content. Got %s, want %s",
						string(decompressed), tc.expectedBody)
				}
			} else {
				if rr.Header().Get("Content-Encoding") == "gzip" {
					t.Errorf("Unexpected gzip encoding in response")
				}
				if !bytes.Equal(rr.Body.Bytes(), tc.expectedBody) {
					t.Errorf("Unexpected response content. Got %s, want %s",
						rr.Body.String(), string(tc.expectedBody))
				}
			}
		})
	}
}

// compressBody is a helper function to gzip compress a body
func compressBody(body []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(body)
	gz.Close()
	return buf.Bytes()
}

// decompressBody is a helper function to gzip compress a body
func decompressBody(body []byte) []byte {
	gz, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		panic("Failed to create gzip reader")
	}
	defer gz.Close()
	decompressed, err := io.ReadAll(gz)
	if err != nil {
		panic("Failed to read decompressed content")
	}
	return decompressed
}
