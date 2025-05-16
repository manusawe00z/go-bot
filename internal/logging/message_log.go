package logging

import (
	"fmt"
	"os"
	"time"
)

var (
	messageLogger *os.File
)

// SetupMessageLogger initializes a separate logger for incoming messages
func SetupMessageLogger() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Create a dedicated log file for messages with current date
	fileName := fmt.Sprintf("logs/messages_%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	messageLogger = file
	return nil
}

// LogMessage logs a message to the dedicated message log file
func LogMessage(guildName, guildID, channelName, channelID, username, discriminator, userID, content string) {
	if messageLogger == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] Guild: '%s' (%s), Channel: '#%s' (%s), User: '%s#%s' (%s): %s\n",
		timestamp, guildName, guildID, channelName, channelID, username, discriminator, userID, content)

	messageLogger.WriteString(logEntry)
}

// CloseMessageLogger closes the message log file
func CloseMessageLogger() {
	if messageLogger != nil {
		messageLogger.Close()
	}
}
