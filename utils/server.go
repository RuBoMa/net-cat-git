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
	messageHistory []string
	historyMutex   sync.RWMutex
	clients        = make(map[*Client]bool)
	clientMutex    sync.Mutex
	maxConnections = 10
)

func HandleClientConnection(conn net.Conn) {
	clientMutex.Lock()
	if len(clients) >= maxConnections {
		// Reject the connection if max clients are reached
		log.Println("Connection limit reached. Rejecting client:", conn.RemoteAddr())
		conn.Write([]byte("Server is full. Try again later.\n"))
		conn.Close()
		clientMutex.Unlock()
		return
	}
	clientMutex.Unlock()
	defer func() {
		clientMutex.Lock()
		deleteClientByConn(conn)
		clientMutex.Unlock()
		conn.Close()
		log.Println("Client disconnected")
	}()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	client, err := getClientName(reader, writer, conn)
	if err != nil {
		log.Println("Error getting client name:", err)
		return
	}
	// send message history to the new client
	sendHistory(writer)
	// Announce that the client has joined
	broadcastMessage(fmt.Sprintf("%s has joined the chat...", client.Name))

	for {
		message, err := reader.ReadString('\n')
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

			historyMutex.Lock()
			messageHistory = append(messageHistory, formattedMessage)
			historyMutex.Unlock()

			log.Println("Message recieved:", formattedMessage)

			broadcastMessage(formattedMessage)
		}
	}
}
func broadcastMessage(message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	for client := range clients {
		_, err := client.Writer.WriteString(message + "\n")
		if err != nil {
			log.Println("Error sending message to clinet:", err)
			client.Conn.Close()
			delete(clients, client)
		}
		client.Writer.Flush()
	}
}

func sendHistory(writer *bufio.Writer) {
	historyMutex.RLock()
	defer historyMutex.RUnlock()

	for _, msg := range messageHistory {
		writer.Write([]byte(msg + "\n"))
	}
	writer.Flush()
}

func deleteClientByConn(conn net.Conn) {
	for client := range clients {
		if client.Conn == conn {
			delete(clients, client)
			return
		}
	}
}
