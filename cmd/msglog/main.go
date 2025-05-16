package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Define command line flags
	date := flag.String("date", time.Now().Format("2006-01-02"), "Date to search (YYYY-MM-DD)")
	user := flag.String("user", "", "Filter by username")
	guild := flag.String("guild", "", "Filter by guild name")
	channel := flag.String("channel", "", "Filter by channel name")
	content := flag.String("content", "", "Filter by message content")
	flag.Parse()

	// Construct log file path
	logFile := filepath.Join("logs", fmt.Sprintf("messages_%s.log", *date))

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		fmt.Println("Available log files:")
		listAvailableLogs()
		os.Exit(1)
	}
	defer file.Close()

	// Read and filter logs
	scanner := bufio.NewScanner(file)
	lineCount := 0
	matchCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Apply filters
		if *user != "" && !strings.Contains(line, fmt.Sprintf("User: '%s", *user)) {
			continue
		}
		if *guild != "" && !strings.Contains(line, fmt.Sprintf("Guild: '%s", *guild)) {
			continue
		}
		if *channel != "" && !strings.Contains(line, fmt.Sprintf("Channel: '#%s", *channel)) {
			continue
		}
		if *content != "" && !strings.Contains(line, *content) {
			continue
		}

		// Print matching line
		fmt.Println(line)
		matchCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nFound %d matching entries out of %d total log entries.\n", matchCount, lineCount)
}

// List available log files in the logs directory
func listAvailableLogs() {
	files, err := filepath.Glob("logs/messages_*.log")
	if err != nil {
		fmt.Printf("Error listing log files: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No message log files found.")
		return
	}

	for _, file := range files {
		filename := filepath.Base(file)
		parts := strings.Split(filename, "_")
		if len(parts) >= 2 {
			date := strings.TrimSuffix(parts[1], ".log")
			fmt.Printf(" - %s\n", date)
		} else {
			fmt.Printf(" - %s\n", filename)
		}
	}
}
