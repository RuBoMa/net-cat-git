package utils

import (
	"fmt"
	"os"
)

const defaultPort = "8989"

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
