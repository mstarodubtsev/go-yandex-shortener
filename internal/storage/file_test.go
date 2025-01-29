package storage

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

// setup function to initialize common test data
func setup() {
	// Initialize logger
	log.InitializeLogger()
	defer log.Logger.Sync()
}

func TestFileStorage_AddURL(t *testing.T) {
	setup()
	file, err := os.CreateTemp("", "storage_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage, _ := NewFileStorage(file.Name())
	storage.AddURL("short1", "http://example.com")

	result, found, _ := storage.GetURL("short1")
	if !found {
		t.Fatalf("Expected URL not found")
	}
	if result != "http://example.com" {
		t.Errorf("Expected %s, got %s", "http://example.com", result)
	}
}

func TestFileStorage_GetURLNotFound(t *testing.T) {
	setup()
	file, err := os.CreateTemp("", "storage_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage, _ := NewFileStorage(file.Name())
	_, found, _ := storage.GetURL("nonexistent")
	if found {
		t.Errorf("Did not expect to find URL")
	}
}

func TestFileStorage_GetAll(t *testing.T) {
	setup()
	file, err := os.CreateTemp("", "storage_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage, _ := NewFileStorage(file.Name())
	storage.AddURL("short1", "http://example1.com")
	storage.AddURL("short2", "http://example2.com")

	allURLs, _ := storage.GetAll()
	if len(allURLs) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(allURLs))
	}
	if allURLs["short1"] != "http://example1.com" {
		t.Errorf("Expected %s, got %s", "http://example1.com", allURLs["short1"])
	}
	if allURLs["short2"] != "http://example2.com" {
		t.Errorf("Expected %s, got %s", "http://example2.com", allURLs["short2"])
	}
}

func TestFileStorage_Concurrency(t *testing.T) {
	setup()
	file, err := os.CreateTemp("", "storage_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage, _ := NewFileStorage(file.Name())
	var wg sync.WaitGroup

	numWrites := 100
	for i := 0; i < numWrites; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			url := "http://example.com/" + strconv.Itoa(i)
			storage.AddURL("short"+strconv.Itoa(i), url)
		}(i)
	}

	wg.Wait()

	allURLs, _ := storage.GetAll()
	if len(allURLs) != numWrites {
		t.Errorf("Expected %d URLs, got %d", numWrites, len(allURLs))
	}
}

func TestFileStorage_RestoreCounter(t *testing.T) {
	setup()
	file, err := os.CreateTemp("", "storage_test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage, _ := NewFileStorage(file.Name())
	storage.AddURL("short1", "http://example1.com")
	storage.AddURL("short2", "http://example2.com")

	file.Close()
	newStorage, _ := NewFileStorage(file.Name())

	if atomic.LoadInt64(&newStorage.counter) != 2 {
		t.Errorf("Expected counter to be restored to 2, got %d", newStorage.counter)
	}
}
