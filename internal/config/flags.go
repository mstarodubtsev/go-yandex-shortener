package config

import (
	"flag"
)

// FlagRunAddr address and port to run server
var FlagRunAddr string

// FlagResultUrl address and port to result url
var FlagResultURL string

// ParseFlags parses flags
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagResultURL, "b", "localhost:8080", "address and port for result URL")
	flag.Parse()
}
