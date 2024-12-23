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

// getClientName prompts the client to enter their name, ensuring it's not empty or already taken.
// It returns a Client object with the provided name and connection details.
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
		activeclientMutex.Lock()
		if isNameTaken(name) {
			activeclientMutex.Unlock()
			writer.WriteString("Name is already taken. Try another.\n")
			writer.Flush()
			continue
		}
		activeclientMutex.Unlock()
		break
	}

	client := &Client{Name: name, Conn: conn, Writer: writer}
	activeclientMutex.Lock()
	activeclients[client] = true
	activeclientMutex.Unlock()

	return client, nil
}

// isNameTaken checks if the given name is already taken by another client.
// It returns true if the name is taken, otherwise false.
func isNameTaken(name string) bool {
	for client := range activeclients {
		if client.Name == name {
			return true
		}
	}
	return false
}

// HandleNameChange allows a client to change their name, ensuring the new name is not already taken.
// It broadcasts the name change to all other clients.
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
