package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Log levels
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logger    *log.Logger
	logLevel  int
	logFile   *os.File
	levelText = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

// Setup initializes the logger
func Setup(level int) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Create log file with current date
	fileName := fmt.Sprintf("logs/bot_%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Set up logger
	logFile = file
	logger = log.New(file, "", log.Ldate|log.Ltime)
	logLevel = level

	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// Debug logs debug information
func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		logMessage(DEBUG, format, v...)
	}
}

// Info logs general information
func Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		logMessage(INFO, format, v...)
	}
}

// Warn logs warning information
func Warn(format string, v ...interface{}) {
	if logLevel <= WARN {
		logMessage(WARN, format, v...)
	}
}

// Error logs error information
func Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		logMessage(ERROR, format, v...)
	}
}

// Fatal logs fatal error and exits
func Fatal(format string, v ...interface{}) {
	if logLevel <= FATAL {
		logMessage(FATAL, format, v...)
		os.Exit(1)
	}
}

// Internal function to log messages
func logMessage(level int, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	prefix := fmt.Sprintf("[%s] ", levelText[level])
	logger.Print(prefix + message)

	// Also output to console
	fmt.Println(prefix + message)
}
