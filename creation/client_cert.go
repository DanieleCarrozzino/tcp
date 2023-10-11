package creation

import (
	"crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
	"path/filePath"
    "fmt"
    "os"
	"io/ioutil"
	"crypto/x509/pkix"
	"math/big"
	"time"
	"encoding/asn1"
	utility   "go/tcp/carrozzino/utility"
)

var parentDir string 		= utility.GetCurrentDirectory()
var caCrtPath string 		= filepath.Join(parentDir, "ca", "ca.crt")
var caKeyPath string		= filepath.Join(parentDir, "ca", "ca.key")
var clientKeyPath string	= filepath.Join(parentDir, "client", "client.key")
var clientCsrPath string	= filepath.Join(parentDir, "client", "client.csr")
var clientCrtPath string	= filepath.Join(parentDir, "client", "client.crt")

type Gen struct {
	template *x509.Certificate
	templateRequest *x509.CertificateRequest
} 

func (g *Gen) CreateKey(){
	// Generate an RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA private key:", err)
		return
	}

	// Create a PKCS#1 private key DER encoding
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// Create a PEM block for the private key
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}

	// Write the private key to a file
	privateKeyFile, err := os.Create(clientKeyPath)
	if err != nil {
		fmt.Println("Error creating private key file:", err)
		return
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		fmt.Println("Error writing private key to file:", err)
		return
	}

	fmt.Println("Generate client.key")
}

func (g *Gen) CreateCSR(){
	// Load the private key from the file
	privateKeyFile, err := os.ReadFile(clientKeyPath)
	if err != nil {
		fmt.Println("Error reading private key file:", err)
		return
	}

	block, _ := pem.Decode(privateKeyFile)
	if block == nil {
		fmt.Println("Error decoding private key PEM block")
		return
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return
	}

	//define template
	g.defineTemplateRequest()

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, g.templateRequest, privateKey)
	if err != nil {
		fmt.Println("Error creating CSR:", err)
		return
	}

	// Create a PEM encoding of the CSR
	csrPEM := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	}

	// Write the CSR to a file
	csrFile, err := os.Create(clientCsrPath)
	if err != nil {
		fmt.Println("Error creating CSR file:", err)
		return
	}
	defer csrFile.Close()

	err = pem.Encode(csrFile, &csrPEM)
	if err != nil {
		fmt.Println("Error writing CSR to file:", err)
		return
	}

	fmt.Println("Generate client.csr")
}

func (g *Gen) CreateCRT(){
	// Load the CA certificate and private key
	caCertFile, err := ioutil.ReadFile(caCrtPath)
	if err != nil {
		fmt.Println("Error reading CA certificate file:", err)
		return
	}
	
	caKeyFile, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		fmt.Println("Error reading CA private key file:", err)
		return
	}

	caCertBlock, _ := pem.Decode(caCertFile)
	if caCertBlock == nil {
		fmt.Println("Error decoding CA certificate PEM block")
		return
	}

	caKeyBlock, _ := pem.Decode(caKeyFile)
	if caKeyBlock == nil {
		fmt.Println("Error decoding CA private key PEM block")
		return
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing CA certificate:", err)
		return
	}

	caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing CA private key:", err)
		return
	}

	// Load the CSR from the file
	csrFile, err := ioutil.ReadFile(clientCsrPath)
	if err != nil {
		fmt.Println("Error reading CSR file:", err)
		return
	}

	csrBlock, _ := pem.Decode(csrFile)
	if csrBlock == nil {
		fmt.Println("Error decoding CSR PEM block")
		return
	}
	
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing CSR:", err)
		return
	}

	// define template
	g.defineTemplate(csr.Subject)

	// Sign the certificate with the CA private key
	certDER, err := x509.CreateCertificate(rand.Reader, g.template, caCert, csr.PublicKey, caKey)
	if err != nil {
		fmt.Println("Error creating certificate:", err)
		return
	}

	// Create a PEM encoding of the certificate
	certPEM := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}

	// Write the certificate to a file
	certFile, err := os.Create(clientCrtPath)
	if err != nil {
		fmt.Println("Error creating certificate file:", err)
		return
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &certPEM)
	if err != nil {
		fmt.Println("Error writing certificate to file:", err)
		return
	}

	fmt.Println("Generate client.crt")
}

func (g *Gen) defineTemplate(subject pkix.Name){

	g.template = &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"Continent"},
			Province:           []string{"State"},
			Locality:           []string{"Locality"},
			Organization:       []string{"Organization"},
			OrganizationalUnit: []string{"Organizational Unit"},
			CommonName:         "Common Name",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:	false,
		DNSNames: []string{"daniele.carrozzino.server, daniele.carrozzino.server_support"},
	}
} 

func (g *Gen) defineTemplateRequest(){

	// Create a new certificate template
	g.templateRequest = &x509.CertificateRequest{}

	// Set the Distinguished Name (DN) fields
    g.templateRequest.Subject = pkix.Name{
        Country:            []string{"EU"},
        Province:           []string{"Italy"},
        Locality:           []string{"Piedmont"},
        Organization:       []string{"carrozzino"},
        OrganizationalUnit: []string{"IT"},
        CommonName:         "Daniele auto gen 2",
    }

    // Key Usage extension
    keyUsageExtension := pkix.Extension{
        Id:       asn1.ObjectIdentifier{2, 5, 29, 15}, // Key Usage OID
        Critical: true,                               // Set to true if Key Usage is critical
        Value:    []byte{0x03, 0x02, 0x04},            // Key Usage flags (keyEncipherment, dataEncipherment)
    }
    g.templateRequest.ExtraExtensions = append(g.templateRequest.ExtraExtensions, keyUsageExtension)

    // Extended Key Usage extension
    extKeyUsageExtension := pkix.Extension{
        Id:       asn1.ObjectIdentifier{2, 5, 29, 37}, // Extended Key Usage OID
        Critical: true,                               // Set to true if Extended Key Usage is critical
        Value:    []byte{0x03, 0x02, 0x01, 0x01},      // Extended Key Usage flags (serverAuth, clientAuth)
    }
    g.templateRequest.ExtraExtensions = append(g.templateRequest.ExtraExtensions, extKeyUsageExtension)

	g.templateRequest.DNSNames = []string{
        "daniele.carrozzino.server",
        "daniele.carrozzino.server_support",
    }
} 