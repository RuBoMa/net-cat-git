package main

import (
	"TCPChat/utils"
	"errors"
	"fmt"
	"net"
	"sync"
)

// Handle user input for shutdown
// func shutDowner(shutdown *chan struct{}) {
// 	reader := bufio.NewReader(os.Stdin)
// 	for {
// 		char, _, err := reader.ReadRune()
// 		if err != nil {
// 			fmt.Println("Error reading input:", err)
// 			continue
// 		}
// 		if char == 'q' || char == 'Q' {
// 			fmt.Println("\nInitiating shutdown...")
// 			close(*shutdown)
// 			return
// 		}
// 	}

// }

func acceptConnections(shutdown chan struct{}, listener net.Listener, wg *sync.WaitGroup) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				fmt.Printf("Error accepting connection: %v\n", err)
			}
			continue
		}
		wg.Add(1)
		go utils.HandleClient(conn, shutdown, wg)

	}
}

func main() {
	// logging.InitLogger()
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
	// go shutDowner(&shutdown)
	go acceptConnections(shutdown, listener, &wg)

	<-shutdown
	// fmt.Println("Waiting for all connections to close...")

	wg.Wait()
	// fmt.Println("Server shutdown complete")
}
