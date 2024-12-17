package utils

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

func AcceptConnections(shutdown chan struct{}, listener net.Listener, wg *sync.WaitGroup) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				fmt.Printf("Error accepting connection: %v\n", err)
			}
			continue
		}
		wg.Add(1)
		go HandleClient(conn, shutdown, wg)

	}
}
