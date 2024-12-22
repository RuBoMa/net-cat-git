package main

import (
	"TCPChat/utils"
	"fmt"
	"net"
	"sync"
)

func main() {
	port := utils.GetPort()

	shutdown := make(chan struct{}) // Used to signal server shutdown.
	var wg sync.WaitGroup

	// Start listening for incoming TCP connections.
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on the port :%s\n", port)
	fmt.Println("Press Ctrl+C to shutdown the server")

	// Start goroutines for broadcasting messages and accepting connections.
	go utils.BroadcastMessages()
	go utils.AcceptConnections(shutdown, listener, &wg)

	<-shutdown // Block until a shutdown signal is received.

	wg.Wait() // Wait for all client handlers to finish.
}
