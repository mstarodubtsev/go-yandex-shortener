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
	// Initialize logger
	log.InitializeLogger()
	defer log.Logger.Sync()
}

// mockHandler is a test handler that writes a simple response
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

// TestWithCompressingPost tests the WithCompressingPost middleware
func TestWithCompressingPost(t *testing.T) {
	setup()
	testCases := []struct {
		name            string
		contentEncoding string
		requestBody     []byte
		expectedBody    []byte
	}{
		{
			name:            "Uncompressed Request",
			contentEncoding: "",
			requestBody:     []byte("test body"),
			expectedBody:    []byte("test body"),
		},
		{
			name:            "Gzipped Request",
			contentEncoding: "gzip",
			requestBody:     compressBody([]byte("test body")),
			expectedBody:    []byte("test body"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var receivedBody []byte
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}
				receivedBody = body
			})

			wrappedHandler := WithCompressingPost(handler)

			req, err := http.NewRequest("POST", "/", bytes.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tc.contentEncoding != "" {
				req.Header.Set("Content-Encoding", tc.contentEncoding)
			}

			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			if !bytes.Equal(receivedBody, tc.expectedBody) {
				t.Errorf("Unexpected request body. Got %v, want %v", string(receivedBody), string(tc.expectedBody))
			}
		})
	}
}

// TestWithCompressingGet tests the WithCompressingGet middleware
func TestWithCompressingGet(t *testing.T) {
	setup()
	testCases := []struct {
		name            string
		acceptEncoding  string
		expectedEncoded bool
	}{
		{
			name:            "Gzip Supported",
			acceptEncoding:  "gzip",
			expectedEncoded: true,
		},
		{
			name:            "Gzip Not Supported",
			acceptEncoding:  "",
			expectedEncoded: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(mockHandler)
			wrappedHandler := WithCompressingGet(handler)

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tc.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tc.acceptEncoding)
			}

			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			// Check content encoding header
			if tc.expectedEncoded {
				if rr.Header().Get("Content-Encoding") != "gzip" {
					t.Errorf("Expected gzip encoding, got none")
				}

				// Verify the response can be decompressed
				gz, err := gzip.NewReader(bytes.NewReader(rr.Body.Bytes()))
				if err != nil {
					t.Fatalf("Failed to create gzip reader: %v", err)
				}
				defer gz.Close()

				decompressed, err := io.ReadAll(gz)
				if err != nil {
					t.Fatalf("Failed to read decompressed content: %v", err)
				}

				if string(decompressed) != "Hello, World!" {
					t.Errorf("Unexpected decompressed content. Got %s, want Hello, World!", string(decompressed))
				}
			} else {
				if rr.Header().Get("Content-Encoding") == "gzip" {
					t.Errorf("Unexpected gzip encoding")
				}
				if rr.Body.String() != "Hello, World!" {
					t.Errorf("Unexpected response content: %s", rr.Body.String())
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
