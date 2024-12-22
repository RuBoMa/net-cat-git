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
		log.Fatalln("coulden't listen to network")
	}
	defer listener.Close()

	log.Println("Server started on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("err while accept", err)
		}

		log.Println("New client connected")
		go utils.HandleClientConnection(conn)
	}

}