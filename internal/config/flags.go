package config

import (
	"flag"
)

// FlagRunAddr address and port to run server
var flagRunAddr string

// FlagResultUrl address and port to result url
var flagBaseURL string

// ParseFlags parses flags
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080", "base URL for result")
	flag.Parse()
}
