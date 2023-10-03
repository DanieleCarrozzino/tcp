# Go TCP Server and Client with TLS
## Overview
This project demonstrates a simple TCP server and client implemented in Go, using the Transport Layer Security (TLS) protocol for secure communication. The project is structured into multiple folders, including client, server, utility, structure, and tlsconfigurator, to maintain a clean and organized codebase.

### Features
. Server Implementation: The project provides a flexible server implementation with an interface, ServerInterface, defining methods like GetDns(), Start(), Read(conn net.Conn), and WriteFile(buffer *bytes.Buffer, fileExt string) error. This allows for easy customization and extension of the server's functionality.
. TLS Encryption: TLS is used to encrypt the communication between the server and client, ensuring data confidentiality and security.
. Client Implementation: The client is designed for simplicity, with a single main method, SendFile, responsible for sending files over the TCP connection.

### Project Structure
- client: Contains the client-side code.
- server: Houses the server-side code.
- utility: Provides utility functions and helper methods.
- structure: Defines data structures and interfaces.
- tlsconfigurator: Contains TLS configuration related code.

### Usage
#### Server
To run the server, navigate to the server directory and execute the appropriate commands. Customize the server implementation as needed by implementing the ServerInterface methods in your chosen server module.

#### Client
To use the client, navigate to the client directory and execute the SendFile command. You can modify this client to suit your specific use case if needed.

#### Getting Started
Clone this repository to your local machine.
Customize the server and client modules to meet your project requirements.
Build and run the server and client components as needed.
#### Contributors
Daniele Carrozzino - carroch97@outlook.it
#### License
This project is licensed under the MIT License - see the LICENSE file for details.
