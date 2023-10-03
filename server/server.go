package server

import (
	"fmt"
	"os"
	"crypto/tls"
	"net"
	"log"
	"io"
	"bytes"
	"encoding/binary"
	"path/filePath"
	tls_config "go/tcp/carrozzino/tls"
	utility    "go/tcp/carrozzino/utility"
)

type Server struct {
	
}

func (s *Server) GetDns() string {
	return "daniele.carrozzino.server"
}

/**
* Main method.
* It starts the main server's loop,
* it starts a tls protocol communication
* over the 3000
*/
func (s *Server) Start() {
	//|***************
	//| 	TLS 	 |
	//****************
	fmt.Println(">Get tls config")
	config, err := tls_config.GetServerTlsConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//|***************
	//| 	TCP 	 |
	//****************
	fmt.Println(">>Start listen tcp")
	ln, err := tls.Listen("tcp", ":3000", config)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(">>> Starting read")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		go s.Read(conn)
	}
}

/**
* Read data method
*/
func (s *Server) Read(conn net.Conn){

	buf := new(bytes.Buffer)

	for {

		// Get the header of the communication
		header := make([]byte, 16)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			fmt.Println("Failed to read header:", err)
			os.Exit(-1)
		}

		// Extract header
		size 	:= int64(binary.LittleEndian.Uint64(header[:8]))
		fileExt := string(header[8:])
		fmt.Println(size, fileExt, header)
		
		// Extract file
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received %d bytes over the network\n", n)
		
		// Write file on disk
		err2 := s.WriteFile(buf, fileExt)
		if err2 != nil {
			fmt.Println(err2)
			os.Exit(-1)
		} 
		fmt.Printf("File written", n)
	}
}

/**
* Write file on disk
*/
func (s *Server)WriteFile(buffer *bytes.Buffer, fileExt string) error {
	parentDir 	:= utility.GetCurrentDirectory()
	filePath 	:= filepath.Join(parentDir, "output", "output" + fileExt)

    // Open the file for appending (create if it doesn't exist)
    file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return err
    }
    defer file.Close()
    
    // Append the data to the file
    _, err = file.Write(buffer.Bytes())
    if err != nil {
        fmt.Println("Error appending data:", err)
        return err
    }
    return nil
}