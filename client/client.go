package client

import (
	"fmt"
	"io/ioutil"
	"encoding/binary"
	"bytes"
	"io"
	"crypto/tls"
	Path 		"path/filePath"
	utility 	"go/tcp/carrozzino/utility"
	tls_config 	"go/tcp/carrozzino/tls"
)

var parentDir string = utility.GetCurrentDirectory()
var filePath  string = Path.Join(parentDir, "client", "foto.jpg")

func getFileInBytes() []byte {
    // Read the contents of the file into a byte slice
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return nil
    }
	return data
}

/**
* The header is divided into 2 main pieces of information:
* First 8 bytes - Size of the file data
* Next 8 bytes  - Size of the file name
* 
* This method inserts the file data size and file name size into the provided data slice,
* and returns the actual file name as a byte slice.
*
* @param data []byte: The byte slice where the header information will be inserted.
* @param size int64: The size of the file data.
* @param filePath string: The path to the file to extract the file name from.
* @return []byte: The byte slice representing the file name.
*/
func getHeader(data []byte, size int64) []byte {
	
	// File size
	binary.LittleEndian.PutUint64(data[:8], uint64(size))

	// Header size
	fileName := Path.Base(filePath)
	binary.LittleEndian.PutUint64(data[8:], uint64(len(fileName)))

	return []byte(fileName)
}

func SendFile() error {
	//|***************
	//| 	TLS 	 |
	//****************
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
	headerSize 	:= make([]byte, 16)
	header 		:= getHeader(headerSize, int64(size))
	binary.Write(conn, binary.LittleEndian, headerSize)

	// Insert header to the file slice
	file = append(header, file...)
	size = len(file)

	// Send
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes over the network\n", n)
	return nil
}