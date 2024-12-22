package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	chatHistory       []string
	chatHistoryMutex  sync.RWMutex
	activeclients     = make(map[*Client]bool)
	activeclientMutex sync.Mutex
	maxConnections    = 10
)

// HandleClientConnection handles the communication between the server and the connected client.
// It manages the connection, handles client messages, and sends message history.
func HandleClientConnection(conn net.Conn) {
	activeclientMutex.Lock()
	if len(activeclients) >= maxConnections {
		// Reject the connection if max activeclients are reached
		log.Println("Connection limit reached. Rejecting client")
		conn.Write([]byte("Server is full. Try again later.\n"))
		conn.Close()
		activeclientMutex.Unlock()
		return
	}
	activeclientMutex.Unlock()
	defer func() {
		activeclientMutex.Lock()
		deleteClientByConn(conn)
		activeclientMutex.Unlock()
		conn.Close()
		log.Println("Client disconnected")
	}()

	connWriter := bufio.NewWriter(conn)
	connReader := bufio.NewReader(conn)

	client, err := getClientName(connReader, connWriter, conn)
	if err != nil {
		log.Println("Error getting client name:", err)
		return
	}
	// send message history to the new client
	sendHistory(connWriter)
	// Announce that the client has joined
	broadcastMessage(fmt.Sprintf("%s has joined the chat...", client.Name))

	for {
		message, err := connReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from client:", err)
			return
		}
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if strings.HasPrefix(message, "name=") {
			// Handle name change
			newName := strings.TrimSpace(strings.TrimPrefix(message, "name="))
			HandleNameChange(client, newName)
		} else {

			// Add timestamp to the message
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			formattedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, client.Name, message)

			chatHistoryMutex.Lock()
			chatHistory = append(chatHistory, formattedMessage)
			chatHistoryMutex.Unlock()

			log.Println("Message recieved:", formattedMessage)

			broadcastMessage(formattedMessage)
		}
	}
}

// broadcastMessage sends the provided message to all connected clients.
func broadcastMessage(message string) {
	activeclientMutex.Lock()
	defer activeclientMutex.Unlock()

	for client := range activeclients {
		_, err := client.Writer.WriteString(message + "\n")
		if err != nil {
			log.Println("Error sending message to clinet:", err)
			client.Conn.Close()
			delete(activeclients, client)
		}
		client.Writer.Flush()
	}
}

// sendHistory sends the entire chat history to the specified writer (client).
func sendHistory(writer *bufio.Writer) {
	chatHistoryMutex.RLock()
	defer chatHistoryMutex.RUnlock()

	for _, msg := range chatHistory {
		writer.Write([]byte(msg + "\n"))
	}
	writer.Flush()
}

// deleteClientByConn removes a client from the activeClients map based on their connection.
// It is called when a client disconnects.
func deleteClientByConn(conn net.Conn) {
	for client := range activeclients {
		if client.Conn == conn {
			delete(activeclients, client)
			return
		}
	}
}
