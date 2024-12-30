package utils

import (
	"fmt"
	"os"
)

const defaultPort = "8989"

// GetPort checks if a port is provided as a command-line argument.
// If no argument is provided, it returns the default port.
func GetPort() string {
	port := defaultPort
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port") //not accurate when running the program first time
		os.Exit(1)
	}
	return port
}
