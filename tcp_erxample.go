package main

import (
	"log"
	"net"
	"fmt"
	"io"
	"time"
	"bytes"
	"encoding/binary"
	"crypto/tls"
	"os"
    "crypto/x509"
    "errors"
    "io/ioutil"
)

type FileServer struct {

}

func (fs *FileServer) start() {

	fmt.Println("> Start FileServer")

	//|***************
    //| 	TLS		 |
    //****************
    // Load the server's TLS certificate and private key.
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        fmt.Println("Error loading server certificates:", err)
        return
    }

    // Load the CA certificate to verify client certificates.
    caCert, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        fmt.Println("Error loading CA certificate:", err)
        return
    }

    // Create a certificate pool with the CA certificate.
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Create a TLS configuration.
    config := tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caCertPool,
        ClientAuth:   tls.RequireAndVerifyClientCert, // Require client certificate and verify it
    }

	fmt.Println(">> Define the certs")

	//|***************
	//| 	TCP 	 |
	//****************
	ln, err := tls.Listen("tcp", ":3000", &config)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Println(">>> Starting read")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go fs.readLoop(conn)
	}

} 

func writeFile(buffer *bytes.Buffer) error {
	filePath := "output.txt"

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
	
	
	

func (fs *FileServer) readLoop(conn net.Conn){

	buf := new(bytes.Buffer)

	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf.Bytes())
		fmt.Printf("Received %d bytes over the network\n", n)
		err2 := writeFile(buf)
		if err2 != nil {
			fmt.Println(err2)
			os.Exit(-1)
		} 
		fmt.Printf("File written", n)
	}
}

func getFileInBytes() []byte {
    filePath := "example.txt"

    // Read the contents of the file into a byte slice
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return nil
    }
	return data
}

func sendFile() error {

	fmt.Println("> Sending file")

	//|***************
	//| 	TLS		 |
	//****************
	// Load the CA certificate for verifying the server's certificate.
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		return err
	}
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		return err
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caCert) {
		return errors.New("failed to append CA certificate")
	}
	config := tls.Config{
		RootCAs:      	rootCAs,
		Certificates: 	[]tls.Certificate{clientCert},
		ServerName: 	"daniele.carrozzino.server",
	}

	fmt.Println(">> Define the certs")

	//|***************
	//| 	TCP 	 |
	//****************
	// File example.txt
	file := getFileInBytes()
	size := len(file)

	// Open the tcp tunnel
	conn, err3 := tls.Dial("tcp", ":3000", &config)
	if err3 != nil {
		return err3
	}

	fmt.Println(">>> Open the connection")

	// Define the dimension of the file
	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err4 := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err4 != nil {
		return err4
	}
	fmt.Printf("Written %d bytes over the network\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(1 * time.Second)
		err := sendFile()
		if err != nil {
			fmt.Println("Got this error : ", err)
			os.Exit(-1)
		}
	}()
	server := &FileServer{}
	server.start()
}