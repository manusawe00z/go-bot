package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Discord Configuration
	DiscordToken  string
	CommandPrefix string

	// TTS Configuration
	TTSLanguage string
	TTSVoice    string

	// Queue Configuration
	QueueSize       int
	QueueTimeoutSec int

	// Audio Quality
	AudioQuality string

	// Logging Level (0=debug, 1=info, 2=warn, 3=error, 4=fatal)
	LogLevel int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		DiscordToken:    getEnv("DISCORD_TOKEN", ""),
		CommandPrefix:   getEnv("COMMAND_PREFIX", "!"),
		TTSLanguage:     getEnv("TTS_LANGUAGE", "en"),
		TTSVoice:        getEnv("TTS_VOICE", ""),
		QueueSize:       getEnvAsInt("QUEUE_SIZE", 100),
		QueueTimeoutSec: getEnvAsInt("QUEUE_TIMEOUT_SECONDS", 60),
		AudioQuality:    getEnv("AUDIO_QUALITY", "medium"),
		LogLevel:        getEnvAsInt("LOG_LEVEL", 1),
	}

	return cfg
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
