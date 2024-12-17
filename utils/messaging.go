package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

func listenForMessages(client *Client, name string, reader *bufio.Reader, writer *bufio.Writer) {
	// Listen for messages
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client "+`"`+name+`"`+" exiting?", err)
			if errors.Is(err, io.EOF) {
			}
			break
		}

		//	Removing characters when clients inputs backspace and cleaning leading and trailing whitespaces
		msg = cleanMessage(msg)

		if msg == "" {
			continue
		}

		if len(msg) > 4 && msg[:5] == "name=" { // change name if message starts thus
			newName := strings.TrimSpace(msg[5:])
			if newName == "" || clientNameExists(newName) {
				if newName == "" {
					WriteToClient("Invalid new name. Name not changed.", writer, true)
				} else {
					WriteToClient("Name already taken. Name not changed.", writer, true)
				}
				writer.Flush()
				continue
			}

			broadcast(fmt.Sprintf("%s has changed their name to %s", name, newName))
			name = newName
			clientMutex.Lock()
			client.Name = newName
			clientMutex.Unlock()
			continue
		}

		if strings.ToLower(msg) == "exit" || strings.ToLower(msg) == "quit" {
			fmt.Println("Client exiting")
			WriteToClient("Exiting chat room.", writer, true)
			writer.Flush()
			break
		}

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
