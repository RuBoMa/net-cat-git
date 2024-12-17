package utils

import (
	"fmt"
	"net"
	"sync"
)
// listens for incoming client connections on the given listener.
func AcceptConnections(shutdown chan struct{}, listener net.Listener, wg *sync.WaitGroup) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			return // Exit the loop on any error (including listner close)
		}
		wg.Add(1)
		go HandleClient(conn, shutdown, wg)

	}
}
