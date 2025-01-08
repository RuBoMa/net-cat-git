package main

import (
	"TCPChat/utils"
	"log"
	"net"
)

func main() {
	port := utils.GetPort()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("Couldn't listen to network:", err) // adding error to logging
	}
	defer listener.Close()

	log.Println("Server started on port", port)

	utils.InitializeLog()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error while accepting connection:", err) //no fatal error when one connection fails
			continue
		}

		log.Println("New client connected")
		go utils.HandleClientConnection(conn)
	}
}
