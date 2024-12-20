package utils

import (
	"TCPChat/logging"
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// broadcast sends the string to the messageBuffer channel
func broadcast(msg string) {
	messageBuffer <- msg
}

// broadcastMessages stores and sends messages to all clients from the message buffer
func BroadcastMessages() {
	for msg := range messageBuffer { // messageBuffer is a channel so for loop doesn't exit
		messageMutex.Lock()
		messages = append(messages, msg)
		messageMutex.Unlock()

		clientMutex.Lock()
		for client := range clients {
			WriteToClient(msg, client.Writer, true) // write the messsage into the buffer managed by the bufio.Writer
			client.Writer.Flush()                   // force write the buffered data to the underlying connection
		}
		clientMutex.Unlock()
	}
}

// sendHistory sends all the previous messages to a clients writer
func sendHistory(writer *bufio.Writer, address, name string) {
	messageMutex.RLock()
	for _, msg := range messages {
		WriteToClient(msg, writer, true)
	}
	messageMutex.RUnlock()
	writer.Flush()
	logging.LogEvent("MSG_HISTORY", address, "Message history sent to client "+`"`+name+`"`)
}

func listenForMessages(client *Client, name, address string, reader *bufio.Reader, writer *bufio.Writer) {
	// Listen for messages
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client "+`"`+name+`"`+" exiting?", err.Error())
			if errors.Is(err, io.EOF) {
				logging.LogEvent("DISCONNECT", address, `"`+name+`"`+" disconnected (client terminated)")
			} else {
				logging.LogEvent("DISCONNECT", address, `"`+name+`"`+" disconnected")
			}
			break
		}

		//	Removing characters when clients inputs backspace and cleaning leading and trailing whitespaces
		msg = cleanMessage(msg)

		if msg == "" {
			logging.LogEvent("NOT_SENT", address, "Empty message from "+`"`+name+`"`+" not sent")
			continue
		}

		if len(msg) > 4 && msg[:5] == "name=" { // change name if message starts thus
			newName := strings.TrimSpace(msg[5:])
			if newName == "" || clientNameExists(newName) {
				if newName == "" {
					WriteToClient("Invalid new name. Name not changed.", writer, true)
					logging.LogEvent("INVALID_NAME", address, "Invalid new name request: "+`"`+newName+`"`+" from "+`"`+name+`"`)
				} else {
					WriteToClient("Name already taken. Name not changed.", writer, true)
					logging.LogEvent("NAME_TAKEN", address, "Taken name change requested: "+`"`+newName+`"`+" by "+`"`+name+`"`)
				}
				writer.Flush()
				continue
			}

			broadcast(fmt.Sprintf("%s has changed their name to %s", name, newName))
			logging.LogEvent("NEW_NAME", address, `"`+name+`"`+" changed their name to "+`"`+newName+`"`)
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
			logging.LogEvent("DISCONNECT", address, `"`+name+`"`+" left the chat")
			break
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg)
		logging.LogEvent("MESSAGE", name, msg)
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
