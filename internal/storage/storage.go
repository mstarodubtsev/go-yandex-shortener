package storage

// Storage interface
type Storage interface {
	// AddURL adds url to storage
	AddURL(hash, url string)

	// GetURL gets url from storage
	GetURL(hash string) (string, bool)

	// GetAll gets all urls from storage
	GetAll() map[string]string
}
