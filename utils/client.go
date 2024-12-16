package utils

import (
	"TCPChat/logging"
	"bufio"
	"fmt"
	"log"
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

func clientNameExists(name string) bool {
	for client := range clients {
		if client.Name == name {
			return true
		}
	}
	return false
}

func getClientName(address *string, reader *bufio.Reader, writer *bufio.Writer) (*Client, string, bool) {
	var client *Client
	var name string
	open := true
	for {
		// Get proposed name
		var err error
		name, err = reader.ReadString('\n')
		name = cleanMessage(name)
		if err != nil {
			logging.LogEvent("INVALID_NAME", *address, "Error reading name request: "+err.Error())
			open = false
			break
		}

		if name == "" {
			logging.LogEvent("INVALID_NAME", *address, "Invalid new name request: "+`""`)
			WriteToClient("Invalid name.\r\n[ENTER YOUR NAME]: ", writer, false)
			writer.Flush()
		}

		// Add client to the map if possible
		clientMutex.Lock()

		if clientNameExists(name) {
			WriteToClient("Name already taken.\r\n[ENTER YOUR NAME]: ", writer, false)
			writer.Flush()
			clientMutex.Unlock()
			logging.LogEvent("NAME_TAKEN", *address, "Taken name requested: "+`"`+name+`"`)
			continue
		}

		if len(clients) >= 10 {
			WriteToClient("Chat room is full. Connection closed.", writer, true)
			writer.Flush()
			clientMutex.Unlock()
			logging.LogEvent("ROOM_FULL", *address, "Chat room full, connection closed")
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

func HandleClient(conn net.Conn, shutdown chan struct{}, wg *sync.WaitGroup) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send welcome message to the new connection
	WriteToClient("Welcome to TCP-Chat!", writer, true)
	WriteToClient(linuxLogo()+"\r\n[ENTER YOUR NAME]: ", writer, false)
	writer.Flush()

	// Use client's IP address in logs
	address, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Printf("Failed to parse client address: %v", err)
		return
	}

	// exit at close of shutdown
	go func() {
		<-shutdown
		conn.Write([]byte("Server shutting down\n"))
		conn.Close()
	}()

	client, name, open := getClientName(&address, reader, writer)

	// Proceed normally if connection didn't close while getting name
	if open {
		logging.LogEvent("CONNECT", address, "Client "+`"`+name+`"`+" joined")
		sendHistory(writer, address, name)
		broadcast(name + " has joined the chat...")
		listenForMessages(client, name, address, reader, writer)
	}

	// Handle client disconnection
	clientMutex.Lock()
	delete(clients, client)
	clientMutex.Unlock()
	broadcast(fmt.Sprintf("%s has left the chat...", name))
	wg.Done()
}
