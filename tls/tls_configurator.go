package tls

import (
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"path/filePath"
	"errors"
	"time"
	"encoding/pem"
	utility   "go/tcp/carrozzino/utility"
)

func GetServerTlsConfig() (*tls.Config, error) {

	// Get the server's TLS certificate and private key.
	// Get the ca's certificate
	parentDir 	:= utility.GetCurrentDirectory()
	certPath 	:= filepath.Join(parentDir, "server", "server.crt")
	keyPath 	:= filepath.Join(parentDir, "server", "server.key")
	caPath 		:= filepath.Join(parentDir, "ca", "ca.crt")

	// Check validity
	checkValidity(certPath, utility.SERVER)

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

	// Check validity
	checkValidity(certPath, utility.CLIENT)

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

func checkValidity(certPath string, certType utility.CertType) utility.Result {
	// Load the certificate file
    certPEM, err := ioutil.ReadFile(certPath)
    if err != nil {
        fmt.Println("Error reading certificate:", err)
        return utility.FAILED
    }

    // Parse the certificate PEM block
    block, _ := pem.Decode(certPEM)
    if block == nil {
        fmt.Println("Error decoding certificate PEM")
        return utility.FAILED
    }

    // Parse the certificate
    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        fmt.Println("Error parsing certificate:", err)
        return utility.FAILED
    }


	/*******************
	 * DATA VALIDATION *
	 *******************/

    // Get the current time
    currentTime := time.Now()

    // Check if the certificate is expired
    if currentTime.After(cert.NotAfter) {
        fmt.Println("Certificate has expired")
		return utility.FAILED
    } else {
        fmt.Println("Certificate is still valid until", cert.NotAfter)
    }

	/*******************
	 *    SAME TYPE    *
	 *******************/

	if certType == utility.CLIENT {
		// Check if it's a client certificate
		if cert.KeyUsage&x509.KeyUsageKeyEncipherment == 0 || !containsClientAuthEKU(cert.ExtKeyUsage) {
			return utility.FAILED
		}
	} else if certType == utility.SERVER {
		// Check if it's a server certificate
		if cert.KeyUsage&x509.KeyUsageKeyEncipherment == 0 || !containsServerAuthEKU(cert.ExtKeyUsage) {
			return utility.FAILED
		}
	}
	
	return utility.SUCCESS
}


// Check if ExtendedKeyUsage contains server authentication OID
func containsServerAuthEKU(eku []x509.ExtKeyUsage) bool {
	for _, usage := range eku {
		if usage == x509.ExtKeyUsageServerAuth {
			return true
		}
	}
	return false
}

// Check if ExtendedKeyUsage contains client authentication OID
func containsClientAuthEKU(eku []x509.ExtKeyUsage) bool {
	for _, usage := range eku {
		if usage == x509.ExtKeyUsageClientAuth {
			return true
		}
	}
	return false
}