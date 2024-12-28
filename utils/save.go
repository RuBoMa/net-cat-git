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
	defaultFileName  = "chat_history.txt" // Default for storing messages and to run test.
)

// StoreChat stores a single message in memory and appends it to a file.
func StoreChat(message string, fileName ...string) {
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
		return
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	if err != nil {
		log.Println("Error writing message to file:", err)
	}
}

// SaveHistoryToFile writes the entire in-memory chat history to a file.
func SaveHistoryToFile(fileName string) {
	file, err := os.Create("chat_history.txt")
	if err != nil {
		log.Println("Error creating chat history file:", err)
		return
	}
	defer file.Close()

	chatHistoryMutex.RLock()
	defer chatHistoryMutex.RUnlock()

	for _, message := range chatHistory {
		_, err := file.WriteString(message + "\n")
		if err != nil {
			log.Println("Error writing message to file:", err)
		}
	}

	log.Println("Chat history saved to chat_history.txt")
}

// StartPeriodicSaving periodically saves the chat history to a file in a separate goroutine.
func StartPeriodicSaving(interval time.Duration, fileName string) {
	go func() {
		for {
			time.Sleep(interval)
			SaveHistoryToFile(fileName)
		}
	}()
}
