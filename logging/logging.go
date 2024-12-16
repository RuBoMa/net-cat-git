package logging

import (
	"log"
	"os"
	"time"
)

var logger *log.Logger

// InitLogger sets up a logger to write entries with flags to a log file
func InitLogger() {
	timestamp := time.Now().Format("20060201150405")
	// Create or open the log file
	file, err := os.OpenFile("logs/server"+timestamp+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(file, "", log.LstdFlags)
}

func LogEvent(eventType, clientInfo, details string) {
	logger.Printf("[%s] [%s]: %s", eventType, clientInfo, details)
}
