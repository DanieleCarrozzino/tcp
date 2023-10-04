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
* Start is the main method for the Server.
* It initiates the main server's loop,
* starts a TLS protocol communication over port 3000.
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

	bufHeader 	:= new(bytes.Buffer)
	bufFile 	:= new(bytes.Buffer)

	for {

		// Get the header of the communication
		header := make([]byte, 16)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			fmt.Println("Failed to read header:", err)
			os.Exit(-1)
		}

		// Extract header
		sizeFile 	:= int64(binary.LittleEndian.Uint64(header[:8]))
		sizeHeader 	:= int64(binary.LittleEndian.Uint64(header[8:]))
		
		// Extract header
		n, err := io.CopyN(bufHeader, conn, sizeHeader)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received %d bytes - header\n", n)
		fileName := bufHeader.String()

		// Extract file
		nFile, err := io.CopyN(bufFile, conn, sizeFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received %d bytes - file\n", nFile)
		
		// Write file on disk
		err = s.WriteFile(bufFile, fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		} 

		bufHeader.Reset()
        bufFile.Reset()
	}
}

/**
* Write file on disk
*/
func (s *Server)WriteFile(buffer *bytes.Buffer, fileName string) error {
	parentDir 	:= utility.GetCurrentDirectory()
	filePath 	:= filepath.Join(parentDir, "output", fileName)
	fmt.Println(filePath)

    // Open the file for appending (create if it doesn't exist)
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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