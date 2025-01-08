package utils

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	chatHistory      []string
	chatHistoryMutex sync.RWMutex
	filePath         string
)

// StoreChat stores a single message in memory and appends it to a file.
func StoreChat(message string) error {

	chatHistoryMutex.Lock()
	defer chatHistoryMutex.Unlock()

	chatHistory = append(chatHistory, message)

	// Writing the message to the chat history file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening chat history file:", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	if err != nil {
		log.Println("Error writing message to file:", err)
		return err
	}
	// Log the message
	log.Println("Message received:", message)

	return nil
}

// Creates a logfile for the current session
func InitializeLog() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	filePath = "logs/chat_history" + timestamp + ".txt"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error creating a chat history file:", err)
	}
	defer file.Close()
}
