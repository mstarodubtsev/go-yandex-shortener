package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddURL tests the AddURL function
func TestAddURL(t *testing.T) {
	m := NewMap()
	m.AddURL("hash1", "https://example.com")
	assert.Equal(t, "https://example.com", m.m["hash1"])
}

// TestGetURL tests the GetURL function
func TestGetURL(t *testing.T) {
	m := NewMap()
	m.AddURL("hash1", "https://example.com")
	url, ok := m.GetURL("hash1")
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", url)

	_, ok = m.GetURL("nonexistent")
	assert.False(t, ok)
}

// TestGetAll tests the GetAll function
func TestGetAll(t *testing.T) {
	m := NewMap()
	m.AddURL("hash1", "https://example.com")
	m.AddURL("hash2", "https://example.org")
	all := m.GetAll()
	assert.Equal(t, 2, len(all))
	assert.Equal(t, "https://example.com", all["hash1"])
	assert.Equal(t, "https://example.org", all["hash2"])
}

// TestNewMap tests the NewMap function
func TestNewMap(t *testing.T) {
	m := NewMap()
	assert.NotNil(t, m)
	assert.Equal(t, 0, len(m.m))
}
