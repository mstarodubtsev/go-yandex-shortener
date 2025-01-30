package storage

import (
	"encoding/json"
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"os"
	"sync"
	"sync/atomic"
)

// Line struct to store single URL
type DataRow struct {
	UUID        int64  `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// FileStorage struct to store all URLs
type FileStorage struct {
	mu      sync.RWMutex
	file    *os.File
	counter int64
}

// AddURL adds a URL
func (storage *FileStorage) AddURL(hash, url string) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	uuid := atomic.AddInt64(&storage.counter, 1)
	row := &DataRow{
		UUID:        uuid,
		ShortURL:    hash,
		OriginalURL: url,
	}
	encoder := json.NewEncoder(storage.file)
	if err := encoder.Encode(row); err != nil {
		return err
	}
	return nil
}

// GetURL retrieves a URL
func (storage *FileStorage) GetURL(hash string) (string, bool, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	// Reset file pointer to beginning
	storage.file.Seek(0, 0)
	decoder := json.NewDecoder(storage.file)

	// while the array contains values
	for decoder.More() {
		var row DataRow
		if err := decoder.Decode(&row); err != nil {
			return "", false, err
		}
		if row.ShortURL == hash {
			return row.OriginalURL, true, nil
		}
	}
	return "", false, nil
}

// GetAll retrieves a copy of all URLs
func (storage *FileStorage) GetAll() (map[string]string, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	// Return a copy to avoid exposing internal state
	mCopy := make(map[string]string, 0)

	// Reset file pointer to beginning
	storage.file.Seek(0, 0)
	decoder := json.NewDecoder(storage.file)

	// while the array contains values
	for decoder.More() {
		var line DataRow
		if err := decoder.Decode(&line); err != nil {
			return nil, err
		}
		mCopy[line.ShortURL] = line.OriginalURL
	}
	return mCopy, nil
}

// restoreCounter restores the counter from the file
func (storage *FileStorage) restoreCounter() {
	storage.file.Seek(0, 0)
	decoder := json.NewDecoder(storage.file)

	var lastRow DataRow
	for decoder.More() {
		if err := decoder.Decode(&lastRow); err != nil {
			break
		}
	}
	storage.counter = lastRow.UUID
}

// NewFileStorage creates a new thread-safe file storage
func NewFileStorage(filename string) (*FileStorage, error) {
	log.Infof("Creating file storage: %s", filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	storage := &FileStorage{
		file:    file,
		counter: 0,
	}
	storage.restoreCounter()
	return storage, nil
}
