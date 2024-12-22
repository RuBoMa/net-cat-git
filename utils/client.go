package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Name   string
	Conn   net.Conn
	Writer *bufio.Writer
}

func getClientName(reader *bufio.Reader, writer *bufio.Writer, conn net.Conn) (*Client, error) {
	var name string
	for {
		writer.WriteString(LinuxLogo() + "\n[ENTER YOUR NAME:] ")
		writer.Flush()

		nameInput, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		name = strings.TrimSpace(nameInput)
		if name == "" {
			writer.WriteString("Name cannot be empty.\n")
			writer.Flush()
			continue
		}
		clientMutex.Lock()
		if isNameTaken(name) {
			clientMutex.Unlock()
			writer.WriteString("Name is already taken. Try another.\n")
			writer.Flush()
			continue
		}
		clientMutex.Unlock()
		break
	}

	client := &Client{Name: name, Conn: conn, Writer: writer}
	clientMutex.Lock()
	clients[client] = true
	clientMutex.Unlock()

	return client, nil
}

func isNameTaken(name string) bool {
	for client := range clients {
		if client.Name == name {
			return true
		}
	}
	return false
}
func HandleNameChange(client *Client, newName string) {

	if isNameTaken(newName) {
		client.Writer.WriteString("Name is already taken. Try another.\n")
		client.Writer.Flush()
		return
	}
	oldName := client.Name
	client.Name = newName
	broadcastMessage(fmt.Sprintf("%s has changed their name to %s", oldName, newName))
}
