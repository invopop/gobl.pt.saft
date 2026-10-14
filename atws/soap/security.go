package soap

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"time"
)

// SecurityHeader represents a WS-Security header of a SOAP message
type SecurityHeader struct {
	XMLName       struct{}       `xml:"http://schemas.xmlsoap.org/ws/2002/12/secext wss:Security"`
	UsernameToken *UsernameToken `xml:"wss:UsernameToken"`
}

// UsernameToken represents a WS-Security UsernameToken
type UsernameToken struct {
	Username string `xml:"wss:Username"`
	Password string `xml:"wss:Password"`
	Nonce    string `xml:"wss:Nonce"`
	Created  string `xml:"wss:Created"`
}

// NewSecurityHeader creates a new SecurityHeader for the given user, password, and public key
func NewSecurityHeader(user, password string, publicKey []byte) (*SecurityHeader, error) {
	ks, err := randomKey()
	if err != nil {
		return nil, fmt.Errorf("Error generating key: %w", err)
	}

	nonce, err := encryptRSA(ks, publicKey)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting KS: %w", err)
	}

	encryptedPassword, err := encryptAES(password, ks)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting password: %w", err)
	}

	timestamp := currentTimestamp()

	createdField, err := encryptAES(timestamp, ks)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting timestamp: %w", err)
	}

	return &SecurityHeader{
		UsernameToken: &UsernameToken{
			Username: user,
			Password: encryptedPassword,
			Nonce:    nonce,
			Created:  createdField,
		},
	}, nil
}

// randomKey generates a random 128-bit AES key
func randomKey() ([]byte, error) {
	key := make([]byte, 16) // 128-bit key
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// encryptRSA encrypts the symmetric key using RSA and a public key
func encryptRSA(data []byte, publicKey []byte) (string, error) {
	// Parse the public key
	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to parse public key")
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Extract the public key
	rsaPub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not a valid RSA public key")
	}

	// PKCS #1 v1.5 encryption is deprecated in the standard library, but it's the
	// scheme mandated by the AT web service specification.
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, data) //nolint:staticcheck
	if err != nil {
		return "", err
	}

	// Encode to Base64
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// encryptAES encrypts the password using AES, ECB mode with PKCS5 padding
func encryptAES(password string, key []byte) (string, error) {
	// Validate key length (AES supports 16, 24, or 32 bytes)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Apply PKCS5 padding
	paddedPassword := pkcs5Padding([]byte(password), block.BlockSize())

	// Encrypt using ECB mode
	ciphertext := make([]byte, len(paddedPassword))
	for start := 0; start < len(paddedPassword); start += block.BlockSize() {
		block.Encrypt(ciphertext[start:start+block.BlockSize()], paddedPassword[start:start+block.BlockSize()])
	}

	// Encode the result to Base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// currentTimestamp returns the current UTC timestamp in ISO 8601 format
func currentTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// pkcs5Padding adds PKCS5 padding to the input
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
