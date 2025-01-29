package app

import (
	"bytes"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/config"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setup function to initialize common test data
func setup() {
	config.Config.ServerAddress = "localhost:8080"
	config.Config.BaseURL = "http://localhost:8080"
	config.Config.FileStoragePath = "/tmp/storage.txt"

	// Initialize logger
	log.InitializeLogger()
	defer log.Logger.Sync()

	// init storage
	store, _ := storage.NewFileStorage(config.Config.FileStoragePath)
	SetStore(store)
}

// TestPostURLHandlerJSON tests the PostURLHandlerJSON function
func TestPostURLHandlerJSON(t *testing.T) {
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
			name:   "Valid JSON request",
			method: "POST",
			body:   `{"url": "https://example.com"}`,
			want: want{
				code:        http.StatusCreated,
				body:        "http://localhost:8080/",
				contentType: "application/json",
			},
		},
		{
			name:   "Duplicate URL",
			method: "POST",
			body:   `{"url": "https://example.com"}`,
			want: want{
				code:        http.StatusCreated,
				body:        "http://localhost:8080/",
				contentType: "application/json",
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
			name:   "Malformed JSON",
			method: "POST",
			body:   `{"url": "https://example.com`,
			want: want{
				code: http.StatusBadRequest,
				body: "unexpected EOF\n",
			},
		},
		{
			name:   "Missing URL field",
			method: "POST",
			body:   `{"data": "https://example.com"}`,
			want: want{
				code: http.StatusBadRequest,
				body: "URL cannot be empty\n",
			},
		},
		{
			name:   "Empty URL",
			method: "POST",
			body:   `{"url": ""}`,
			want: want{
				code: http.StatusBadRequest,
				body: "URL cannot be empty\n",
			},
		},
		{
			name:   "Wrong URL",
			method: "POST",
			body:   `{"url": "111"}`,
			want: want{
				code: http.StatusBadRequest,
				body: "URL must start with http:// or https://\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/shorten", bytes.NewBufferString(tt.body))
			res := httptest.NewRecorder()

			PostURLHandlerJSON(res, req)

			result := res.Result()
			defer result.Body.Close()

			bodyBytes, _ := io.ReadAll(result.Body)
			bodyString := string(bodyBytes)

			assert.Equal(t, tt.want.code, result.StatusCode)

			if result.StatusCode == http.StatusCreated {
				assert.True(t, strings.HasPrefix(bodyString, `{"result":"http://localhost:8080/`))
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			} else {
				assert.Equal(t, tt.want.body, bodyString)
			}
		})
	}
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
	store.AddURL("12345678", "https://example.com")

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

// TestRouter tests the Router function
func TestRouter(t *testing.T) {
	setup()
	type want struct {
		code        int
		body        string
		header      string
		contentType string
	}
	tests := []struct {
		name   string
		url    string
		method string
		body   string
		want   want
	}{
		{
			name:   "Valid POST request",
			url:    "/",
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
			url:    "/",
			method: "POST",
			body:   "",
			want: want{
				code: http.StatusBadRequest,
				body: "Empty body\n",
			},
		},
		{
			name:   "Non-POST request method",
			url:    "/",
			method: "GET",
			body:   "https://example.com",
			want: want{
				code: http.StatusMethodNotAllowed,
				body: "",
			},
		},
		{
			url:    "/12345678",
			name:   "Valid GET request with existing hash",
			method: "GET",
			want: want{
				code:   http.StatusTemporaryRedirect,
				header: "https://example.com",
			},
		},
		{
			name:   "Valid GET request with non-existing hash",
			url:    "/nonexistent",
			method: "GET",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "Non-GET request method",
			url:    "/12345678",
			method: "POST",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	ts := httptest.NewServer(Router())
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, bodyString := testRequest(t, ts, tt.method, tt.url, tt.body)
			defer resp.Body.Close()
			if resp.StatusCode != tt.want.code {
				t.Errorf("Expected status %d, got %d", tt.want.code, resp.StatusCode)
			}
			assert.Equal(t, tt.want.code, resp.StatusCode)
			if tt.want.code == http.StatusCreated {
				if !strings.HasPrefix(bodyString, tt.want.body) {
					t.Errorf("Expected response body to start with %s, got %s", tt.want.body, bodyString)
				}
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			} else if tt.want.code == http.StatusTemporaryRedirect {
				location := resp.Header.Get("Location")
				if location != tt.want.header {
					t.Errorf("Expected Location header %s, got %s", tt.want.header, location)
				}
			} else if bodyString != tt.want.body {
				t.Errorf("Expected response body %s, got %s", tt.want.body, bodyString)
			}
		})
	}
}

// testRequest is a helper function to make HTTP requests to the test server
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	// restrict redirects
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	// make request
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
