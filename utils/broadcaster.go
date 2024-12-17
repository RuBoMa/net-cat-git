package utils

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

// broadcast sends the string to the messageBuffer channel
func broadcast(msg string) {
	messageBuffer <- msg
}
