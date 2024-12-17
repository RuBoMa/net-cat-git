package main

import (
	"TCPChat/utils"
	"fmt"
	"net"
	"sync"
)

func main() {
	port := utils.GetPort()

	shutdown := make(chan struct{}) // Only thing ever sent is the closing of the channel
	var wg sync.WaitGroup

	listener, err := net.Listen("tcp", ":"+port) // Listening on TCP network
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on the port :%s\n", port)
	fmt.Println("Press Ctrl+C to shutdown the server")

	// Start message broadcaster
	go utils.BroadcastMessages()
	go utils.AcceptConnections(shutdown, listener, &wg)

	<-shutdown

	wg.Wait()
}
