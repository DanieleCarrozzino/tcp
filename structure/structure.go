package structure

import (
	"net"
	"bytes"
)

type ServerInterface interface {
	GetDns() string
	Start()
	Read(conn net.Conn)
	WriteFile(buffer *bytes.Buffer, fileExt string) error
}

type ClientInterface interface {
	SendFile()
}

type ClientGenerator interface {
	CreateKey()
	CreateCSR()
	CreateCRT()
}