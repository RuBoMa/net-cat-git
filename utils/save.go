package utils

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

var (
	chatHistory      []string
	chatHistoryMutex sync.RWMutex
	//defaultFileName  = "chat_history.txt" // Default for storing messages and to run test.
	logFile *os.File
	logger  *bufio.Writer
)

// StoreChat stores a single message in memory and appends it to a file.
func StoreChat(message string, fileName ...string) error {
	// Determine file name to use
	// targetFile := defaultFileName
	// if len(fileName) > 0 && fileName[0] != "" {
	// 	targetFile = fileName[0]
	// }

	chatHistoryMutex.Lock()
	chatHistory = append(chatHistory, message)
	logger.WriteString(message + "\n")
	logger.Flush()
	chatHistoryMutex.Unlock()

	// Log the message
	log.Println("Message received:", message)

	// Append the message to the file
	// file, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("Error writing to chat history file:", err)
	// 	return err
	// }
	// defer file.Close()

	// _, err = file.WriteString(message + "\n")
	// if err != nil {
	// 	log.Println("Error writing message to file:", err)
	// 	return err
	// }
	return nil
}

// Creates a logfile for the current session and writes a start message to the file
func InitializeLog() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logFile, err := os.OpenFile("logs/chat_history"+timestamp+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error creating a chat history file:", err)
	}
	logger = bufio.NewWriter(logFile)
	logger.WriteString(timestamp + " Session started\n")
	logger.Flush()

}

// Writes the shutdown message to the log and closes the file
func CloseLog() {

	if logger != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logger.WriteString(timestamp + " Session terminated\n")
		logger.Flush()
		if logFile != nil {
			logFile.Close()
		}
	}
}
