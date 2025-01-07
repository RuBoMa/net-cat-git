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
	activeclients     = make(map[*Client]bool)
	activeclientMutex sync.Mutex
	maxConnections    = 10
)

// HandleClientConnection handles the communication between the server and the connected client.
// It manages the connection, handles client messages, and sends message history.
func HandleClientConnection(conn net.Conn) {
	activeclientMutex.Lock()
	if len(activeclients) >= maxConnections {
		// Reject the connection if max active clients are reached
		log.Println("Connection limit reached. Rejecting client")
		conn.Write([]byte("Server is full. Try again later.\n"))
		conn.Close()
		activeclientMutex.Unlock()
		return
	}
	activeclientMutex.Unlock()

	var client *Client
	defer func() {
		activeclientMutex.Lock()
		deleteClientByConn(conn)
		activeclientMutex.Unlock()
		conn.Close()

		//leftMessage := fmt.Sprintf("%s has left the chat...", client.Name)
		leftMessage := client.Name + " has left the chat..."

		if err := StoreChat(leftMessage); err != nil {
			log.Printf("Failed to store left message for %s: %v", client.Name, err)
		}
		broadcastMessage(leftMessage, client)

		log.Printf("%v disconnected\n", client.Name) // logging who left
	}()

	connWriter := bufio.NewWriter(conn)
	connReader := bufio.NewReader(conn)

	client, err := getClientName(connReader, connWriter, conn)
	if err != nil {
		log.Println("Error getting client name:", err)
		return
	}

	sendHistory(connWriter)

	joinMessage := fmt.Sprintf("%s has joined the chat...", client.Name)

	if err := StoreChat(joinMessage); err != nil {
		log.Printf("Failed to store join message for %s: %v", client.Name, err)
	}
	broadcastMessage(joinMessage, client)

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
			newName := strings.TrimSpace(strings.TrimPrefix(message, "name=")) // Is this in the README or are there instructions for the user?
			HandleNameChange(client, newName)
		} else {

			// Add timestamp to the message
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			formattedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, client.Name, message)

			if err := StoreChat(formattedMessage); err != nil {
				log.Printf("Failed to store formatted message for %s: %v", client.Name, err)
			}
			broadcastMessage(formattedMessage, client)
		}
	}
}

// broadcastMessage sends the provided message to all connected clients.
func broadcastMessage(message string, sender *Client) {
	activeclientMutex.Lock()
	defer activeclientMutex.Unlock()

	for client := range activeclients {
		if client == sender {
			continue
		}
		_, err := client.Writer.WriteString(message + "\n")
		if err != nil {
			log.Println("Error broadcasting message to client:", err) //updated error message
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
