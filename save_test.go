package main

import (
	"TCPChat/utils"
	"os"
	"strings"
	"sync"
	"testing"
)

var (
	chatHistory      []string
	chatHistoryMutex sync.RWMutex
)

func setupTest() (string, func()) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "chat_history_test_*.txt")
	if err != nil {
		panic(err)
	}

	// Replace the file used in `StoreChat` and `SaveHistoryToFile`
	chatHistoryMutex.Lock()
	chatHistory = []string{}
	chatHistoryMutex.Unlock()

	// Return cleanup function
	cleanup := func() {
		os.Remove(tmpFile.Name())
	}
	return tmpFile.Name(), cleanup
}

func TestStoreChat(t *testing.T) {
	// Setup test
	fileName := "test_chat_history.txt"
	defer os.Remove(fileName)

	// Simulate storing a message
	message := "Test message"
	utils.StoreChat(message, fileName)

	// Verify the message is in memory
	chatHistoryMutex.RLock()
	defer chatHistoryMutex.RUnlock()
	if len(chatHistory) != 1 || chatHistory[0] != message {
		t.Errorf("Expected message %q in chat history, got %v", message, chatHistory)
	}

	// Verify the message is written to the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read chat history file: %v", err)
	}
	if !strings.Contains(string(data), message) {
		t.Errorf("Expected message %q in chat history file, got %q", message, string(data))
	}
}

func TestSaveHistoryToFile(t *testing.T) {
	// Setup test
	fileName, cleanup := setupTest()
	defer cleanup()

	// Add messages to memory
	messages := []string{"Message 1", "Message 2", "Message 3"}
	chatHistoryMutex.Lock()
	chatHistory = append(chatHistory, messages...)
	chatHistoryMutex.Unlock()

	// Call SaveHistoryToFile
	utils.SaveHistoryToFile(fileName)

	// Verify all messages are written to the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read chat history file: %v", err)
	}
	for _, msg := range messages {
		if !strings.Contains(string(data), msg) {
			t.Errorf("Expected message %q in chat history file, got %q", msg, string(data))
		}
	}
}
