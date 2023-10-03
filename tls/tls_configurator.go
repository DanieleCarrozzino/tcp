package tls

import (
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"path/filePath"
	"errors"
	utility "go/tcp/carrozzino/utility"
)

func GetServerTlsConfig() (*tls.Config, error) {

	// Get the server's TLS certificate and private key.
	// Get the ca's certificate
	parentDir 	:= utility.GetCurrentDirectory()
	certPath 	:= filepath.Join(parentDir, "server", "server.crt")
	keyPath 	:= filepath.Join(parentDir, "server", "server.key")
	caPath 		:= filepath.Join(parentDir, "ca", "ca.crt")

	// Load the server's TLS certificate and private key.
    cert, err := tls.LoadX509KeyPair(certPath, keyPath)
    if err != nil {
        fmt.Println("Error loading server certificates:", err)
        return nil, err
    }

    // Load the CA certificate to verify client certificates.
    caCert, err := ioutil.ReadFile(caPath)
    if err != nil {
        fmt.Println("Error loading CA certificate:", err)
        return nil, err
    }

    // Create a certificate pool with the CA certificate.
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Create a TLS configuration.
    config := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caCertPool,
        ClientAuth:   tls.RequireAndVerifyClientCert, // Require client certificate and verify it
    }
	return config, nil
}

func GetClientTlsConfig() (*tls.Config, error) {
	// Get the server's TLS certificate and private key.
	// Get the ca's certificate
	parentDir 	:= utility.GetCurrentDirectory()
	certPath 	:= filepath.Join(parentDir, "client", "client.crt")
	keyPath 	:= filepath.Join(parentDir, "client", "client.key")
	caPath 		:= filepath.Join(parentDir, "ca", "ca.crt")

	// Load the CA certificate for verifying the server's certificate.
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append CA certificate")
	}
	config := &tls.Config{
		RootCAs:      	rootCAs,
		Certificates: 	[]tls.Certificate{clientCert},
		ServerName: 	"daniele.carrozzino.server",
	}
	return config, nil
}