package utils

import (
	"log"
	"os"
	"sync"
)

var (
	chatHistory      []string
	chatHistoryMutex sync.RWMutex
	defaultFileName  = "chat_history.txt" // Default for storing messages and to run test.
)

// StoreChat stores a single message in memory and appends it to a file.
func StoreChat(message string, fileName ...string) error {
	// Determine file name to use
	targetFile := defaultFileName
	if len(fileName) > 0 && fileName[0] != "" {
		targetFile = fileName[0]
	}

	chatHistoryMutex.Lock()
	chatHistory = append(chatHistory, message)
	chatHistoryMutex.Unlock()

	// Log the message
	log.Println("Message received:", message)

	// Append the message to the file
	file, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error writing to chat history file:", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	if err != nil {
		log.Println("Error writing message to file:", err)
		return err
	}
	return nil
}
