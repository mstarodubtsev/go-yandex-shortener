package app

import (
	"bytes"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setup function to initialize common test data
func setup() {
	config.FlagRunAddr = "localhost:8080"
	config.FlagResultURL = "localhost:8080"
}

// TestPostURLHandler tests the PostURLHandler function
func TestPostURLHandler(t *testing.T) {
	setup()

	type want struct {
		code        int
		body        string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "Valid POST request",
			method: "POST",
			body:   "https://example.com",
			want: want{
				code:        http.StatusCreated,
				body:        "http://localhost:8080/",
				contentType: "text/plain",
			},
		},
		{
			name:   "Empty body",
			method: "POST",
			body:   "",
			want: want{
				code: http.StatusBadRequest,
				body: "Empty body\n",
			},
		},
		{
			name:   "Non-POST request method",
			method: "GET",
			body:   "https://example.com",
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/post", bytes.NewBufferString(tt.body))
			res := httptest.NewRecorder()

			PostURLHandler(res, req)

			result := res.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.want.code {
				t.Errorf("Expected status %d, got %d", tt.want.code, result.StatusCode)
			}

			bodyBytes, _ := io.ReadAll(result.Body)
			bodyString := string(bodyBytes)

			if tt.want.code == http.StatusCreated {
				if !strings.HasPrefix(bodyString, tt.want.body) {
					t.Errorf("Expected response body to start with %s, got %s", tt.want.body, bodyString)
				}
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			} else if bodyString != tt.want.body {
				t.Errorf("Expected response body %s, got %s", tt.want.body, bodyString)
			}
		})
	}
}

// TestGetURLHandler tests the GetURLHandler function
func TestGetURLHandler(t *testing.T) {
	setup()

	// Set up test data in the map
	m["12345678"] = "https://example.com"

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "Valid GET request with existing hash",
			method:         "GET",
			path:           "/12345678",
			expectedStatus: http.StatusTemporaryRedirect,
			expectedHeader: "https://example.com",
		},
		{
			name:           "Valid GET request with non-existing hash",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:           "GET request with malformed path",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Non-GET request method",
			method:         "POST",
			path:           "/12345678",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			res := httptest.NewRecorder()

			GetURLHandler(res, req)

			result := res.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, result.StatusCode)
			}

			// Check for redirect Location header if status is 307
			if tt.expectedStatus == http.StatusTemporaryRedirect {
				location := result.Header.Get("Location")
				if location != tt.expectedHeader {
					t.Errorf("Expected Location header %s, got %s", tt.expectedHeader, location)
				}
			}

			// Verify body content for other status cases
			if tt.expectedStatus != http.StatusTemporaryRedirect {
				bodyBytes, _ := io.ReadAll(result.Body)
				bodyString := string(bodyBytes)
				if bodyString != tt.expectedBody {
					t.Errorf("Expected response body %s, got %s", tt.expectedBody, bodyString)
				}
			}
		})
	}
}
