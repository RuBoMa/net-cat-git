package utils

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	Name   string
	Writer *bufio.Writer
}

var (
	messageBuffer = make(chan string, 10)
	messageMutex  sync.RWMutex
	clientMutex   sync.Mutex
	messages      []string
	clients       = make(map[*Client]bool)
)

// clientNameExists checks if a name is already in use by another client.
func clientNameExists(name string) bool {
	for client := range clients {
		if client.Name == name {
			return true
		}
	}
	return false
}

// getClientName validates and retrieves a unique name for the new client.
func getClientName(reader *bufio.Reader, writer *bufio.Writer) (*Client, string, bool) {
	var client *Client
	var name string
	open := true
	for {
		// Get proposed name
		var err error
		name, err = reader.ReadString('\n')
		name = cleanMessage(name)
		if err != nil {
			open = false
			break
		}

		// Ensure the name is not empty.
		if name == "" {
			WriteToClient("Invalid name.\r\n[ENTER YOUR NAME]: ", writer, false)
			writer.Flush()
		}

		// Add the client to the map if the name is unique and the chat room is not full.
		clientMutex.Lock()

		if clientNameExists(name) {
			WriteToClient("Name already taken.\r\n[ENTER YOUR NAME]: ", writer, false)
			writer.Flush()
			clientMutex.Unlock()
			continue
		}

		if len(clients) >= 10 {
			WriteToClient("Chat room is full. Connection closed.", writer, true)
			writer.Flush()
			clientMutex.Unlock()
			open = false
			return nil, "", false
		}

		client = &Client{name, writer}
		clients[client] = true
		clientMutex.Unlock()
		break
	}
	return client, name, open
}

// HandleClient manages a single client connection and processes their messages.
func HandleClient(conn net.Conn, shutdown chan struct{}, wg *sync.WaitGroup) {
	defer conn.Close() // Ensure the connection is closed when the function exits.
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send welcome message to the new connection
	WriteToClient("Welcome to TCP-Chat!", writer, true)
	WriteToClient(linuxLogo()+"\r\n[ENTER YOUR NAME]: ", writer, false)
	writer.Flush()

	// Gracefully handle server shutdown for the client.
	go func() {
		<-shutdown
		conn.Write([]byte("Server shutting down\n"))
		conn.Close()
	}()

	// Get the client's name and ensure the connection remains open.
	client, name, open := getClientName(reader, writer)
	if open {
		sendHistory(writer)
		broadcast(name + " has joined the chat...")
		listenForMessages(client, name, reader, writer)
	}

	// Handle client disconnection by removing them from the client map and broadcasting their departure.
	clientMutex.Lock()
	delete(clients, client)
	clientMutex.Unlock()
	broadcast(fmt.Sprintf("%s has left the chat...", name))
	wg.Done()
}
