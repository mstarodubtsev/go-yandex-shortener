package config

// AppConfig struct
type AppConfig struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

// Config variable
var Config AppConfig

// ParseConfig parses and returns application configuration
func ParseConfig() {
	env := GetEnvConfig() // Retrieve environment variables
	parseFlags()          // Parse CLI flags

	// Assign default values if environment variables are empty
	Config.ServerAddress = chooseNonEmpty(env.ServerAddress, flagRunAddr)
	Config.BaseURL = chooseNonEmpty(env.BaseURL, flagBaseURL)
	Config.FileStoragePath = chooseNonEmpty(env.FileStoragePath, flagFileStoragePath)
}

// chooseNonEmpty returns the first non-empty string from the arguments
func chooseNonEmpty(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}
