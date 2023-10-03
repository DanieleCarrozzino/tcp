package client

import (
	"fmt"
	"io/ioutil"
	Path "path/filePath"
	"encoding/binary"
	"bytes"
	"io"
	"crypto/tls"
	utility 	"go/tcp/carrozzino/utility"
	tls_config 	"go/tcp/carrozzino/tls"
)

var parentDir string = utility.GetCurrentDirectory()
var filePath  string = Path.Join(parentDir, "client", "example.txt")

func getFileInBytes() []byte {
    // Read the contents of the file into a byte slice
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return nil
    }
	return data
}

func getHeader(data []byte, size int64) {
	binary.LittleEndian.PutUint64(data[:8], uint64(size))
	str := Path.Ext(filePath)
	copy(data[8:], []byte(str))
}

func SendFile() error {
	//|***************
	//| 	TLS 	 |
	//****************
	fmt.Println(">Get tls config")
	config, err := tls_config.GetClientTlsConfig()
	if err != nil {
		return err
	}

	//|***************
	//| 	TCP 	 |
	//****************
	file := getFileInBytes()
	size := len(file)

	// Open the tcp tunnel
	conn, err := tls.Dial("tcp", ":3000", config)
	if err != nil {
		return err
	}

	// Header
	header := make([]byte, 16)
	getHeader(header, int64(size))
	binary.Write(conn, binary.LittleEndian, header)

	// Send
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes over the network\n", n)
	return nil
}