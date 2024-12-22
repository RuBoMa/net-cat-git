package utils

import (
	"fmt"
	"os"
)

const defaultPort = "8989"

// GetPort retrieves the port number from the command-line arguments or uses the default port.
func GetPort() string {
	port := defaultPort
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	return port
}
