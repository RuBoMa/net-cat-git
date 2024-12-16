# NetCat

Something much like the original NetCat, a TCP group chat.

TCPChat is a simple, server-client chat application inspired by NetCat. It allows multiple clients to connect to a TCP server and exchange messages in real-time. It includes name assignment, connection controls, and message logging—all while adhering to good coding practices in Go.

# Features

TCP Server: Accepts up to 10 client connections on a specified port (default 8989).

User Naming: Clients are prompted to choose a unique name. If they wish to change name they can type "name=" followed by their new name of choice.

Broadcast Messaging: All messages are time-stamped and sent to every connected client.

Message History: New clients receive the full chat history upon joining.

Join/Leave Notifications: The group is informed when clients connect or disconnect.

Graceful Shutdown: Server can be stopped gracefully, notifying all clients.

Logging: Events and messages are logged to a file.

# Start the server with default port 8989

To start the server with the default port, run the command "go run ." in the terminal while in the folder that contains "netcat.go"

# Start the server on a custom port

To start the server with a custom port, run the command "go run ." followed by your port number of choice in the terminal while in the folder that contains "netcat.go"

For instance: 
go run . 2525

# Connect a client using netcat

Use the command: 
nc -host -port

For instance:
nc localhost 8989

When a client connects, a Linux logo is displayed, and the user is prompted for a name. Once named, the client can send and receive messages from all other users connected to the server. To exit, simply type exit or quit.

# Requirements

Go (golang)

nc (netcat) 
or telnet

# Additional Notes

This project is a simplified re-creation of NetCat’s chat features, focusing on TCP communication, concurrency, and proper synchronization (using goroutines, channels, and mutexes). It also demonstrates error handling and logging best practices.