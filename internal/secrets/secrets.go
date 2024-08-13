package secrets

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func LoadPrivateKeyFromPEM[K crypto.Signer](privateKeyPEM []byte, password string) (K, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return *new(K), errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return *new(K), err
	}

	key, ok := privateKey.(K)
	if !ok {
		return *new(K), fmt.Errorf("private key is not %T", new(K))
	}

	return key, nil
}
