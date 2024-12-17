# Net-Cat
Description
This project is a recreation of the NetCat tool in a Server-Client architecture. The program allows for real-time group chat functionality over a TCP connection. It can run in server mode, listening for connections, or in client mode, connecting to the server and transmitting messages.

# Features
The project implements a group chat with the following capabilities:

TCP Connection

A server can manage multiple clients (1-to-many relationship).

Name Requirement
Clients must provide a name when connecting.

Clients can send messages to the group chat.
Empty messages are ignored.
Each message includes a timestamp and the sender’s name.
Example format:

[2024-06-17 12:45:00][Alice]: Hello, world!

New clients receive all previous chat messages upon joining.
Join/Exit Notifications

When a client joins, all other clients are notified.
When a client leaves, the group is informed, but the chat continues without disconnection.
Default Port Handling

If no port is specified, the server listens on port 8989.

Usage message when no port is provided:
[USAGE]: ./TCPChat $port
Connection Management

Supports up to 10 simultaneous client connections.
Proper error handling for both server and client sides.

Clients can disconnect and reconnect without disrupting the server or other clients.

# Instructions
Prerequisites
Go installed on your system.

# Run the server

Start the server on specified port: go run . $port

if no port is specified, then default 8989 is used: go run .

# Run the client

nc localhost $port

# Implementation Notes
Written in Go.
Uses goroutines for concurrent client handling.
Utilizes channels and mutexes for synchronization.
Follows Go best practices for code structure and error handling.

# Authors

Roope & Toft