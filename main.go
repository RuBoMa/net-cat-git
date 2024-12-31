package main

import (
	"TCPChat/utils"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	port := utils.GetPort()
	utils.InitializeLog()

	// Create a channel to listen for system signals
	quit := make(chan os.Signal, 1)

	// Notify the quit channel when SIGINT (Ctrl+C)
	signal.Notify(quit, os.Interrupt)

	// goroutine waits for termination signal, closes the log and exits.
	go func() {
		<-quit  
		utils.CloseLog()
		os.Exit(0)
	}()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("Couldn't listen to network:", err) // adding error to logging
	}
	defer listener.Close()

	log.Println("Server started on port", port)

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
