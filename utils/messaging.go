package utils

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

// sendHistory sends all the previous messages to a clients writer
func sendHistory(writer *bufio.Writer) {
	messageMutex.RLock()
	for _, msg := range messages {
		WriteToClient(msg, writer, true)
	}
	messageMutex.RUnlock()
	writer.Flush()
}

// listenForMessages listens for incoming messages from a client, processes the messages, and broadcasts them as necessary.
func listenForMessages(client *Client, name string, reader *bufio.Reader, writer *bufio.Writer) {
	for {
		// Read a message from the client until a newline character is encountered.
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client "+`"`+name+`"`+" exiting?", err)
			break
		}

		// Remove unwanted characters and trim whitespaces from the message.
		msg = cleanMessage(msg)
		// Skip empty messages.
		if msg == "" {
			continue
		}
		// Check if the message starts with "name=" to handle name changes.
		if len(msg) > 4 && msg[:5] == "name=" {
			newName := strings.TrimSpace(msg[5:])
			if newName == "" || clientNameExists(newName) {
				if newName == "" {
					WriteToClient("Invalid new name. Name not changed.", writer, true)
				} else {
					WriteToClient("Name already taken. Name not changed.", writer, true)
				}
				writer.Flush()
			}
			// Broadcast the name change to all clients.
			broadcast(fmt.Sprintf("%s has changed their name to %s", name, newName))
			name = newName
			clientMutex.Lock()
			client.Name = newName
			clientMutex.Unlock()
			continue
		}
		// Handle client exit commands.
		if strings.ToLower(msg) == "exit" || strings.ToLower(msg) == "quit" {
			fmt.Println("Client exiting")
			WriteToClient("Exiting chat room.", writer, true)
			writer.Flush()
			break
		}
		// Format the message with a timestamp and broadcast it.
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg)
		broadcast(formattedMsg)
	}
}

// cleanMessage deals with backspaces and cleans leading and trailing whitespaces
func cleanMessage(msg string) string {
	cleaned := []rune{}
	for _, r := range msg {
		if r == 8 {
			if len(cleaned) > 0 {
				cleaned = cleaned[:len(cleaned)-1]
			}
		} else {
			cleaned = append(cleaned, r)
		}
	}
	return strings.TrimSpace(string(cleaned))
}

// WriteToClient writes the message to the given client. The bool decides whether or not to add a "\r\n" to the end of the message
func WriteToClient(msg string, writer *bufio.Writer, hasNewline bool) {
	if hasNewline {
		msg = msg + "\r\n"
	}
	writer.WriteString(msg)
}
