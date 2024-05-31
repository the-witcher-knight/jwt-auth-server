package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"
)

type SignatureService interface {
	GenerateToken(payload []byte) (string, error)

	GetJWKs() (JWKS, error)
}

type JWK struct {
	Kty string   `json:"kty"`
	X5c []string `json:"x5c"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	Use string   `json:"use"`
	Alg string   `json:"alg"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type rsaSignatureService struct {
	privateKey *rsa.PrivateKey
}

func NewRSASignatureService(privateKeyPath string) (SignatureService, error) {
	privateKey, err := loadPrivateKeyFromPEMFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return rsaSignatureService{
		privateKey: privateKey,
	}, nil
}

func (hdl rsaSignatureService) Algo() string {
	return "RS256"
}

func (hdl rsaSignatureService) GenerateToken(payload []byte) (string, error) {
	header := fmt.Sprintf(`{"alg":"%s","typ":"JWT"}`, hdl.Algo())

	message := base64URLEncode([]byte(header)) + "." + base64URLEncode(payload)

	hashed := sha256.Sum256([]byte(message))

	sig, err := rsa.SignPKCS1v15(rand.Reader, hdl.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return message + "." + base64URLEncode(sig), nil
}

func (hdl rsaSignatureService) GetJWKs() (JWKS, error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "example.com",
			Organization: []string{"Example Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateKey := hdl.privateKey
	pubKey := privateKey.PublicKey

	// Self-sign the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &pubKey, privateKey)
	if err != nil {
		return JWKS{}, fmt.Errorf("failed to create certificate: %v", err)
	}

	return JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				X5c: []string{
					base64.StdEncoding.EncodeToString(certDER),
				},
				N:   base64URLEncode(pubKey.N.Bytes()),
				E:   base64URLEncode(big.NewInt(int64(pubKey.E)).Bytes()),
				Use: "sig",
				Alg: hdl.Algo(),
			},
		},
	}, nil
}

func base64URLEncode(str []byte) string {
	encoded := base64.URLEncoding.EncodeToString(str)
	return strings.TrimRight(encoded, "=")
}

func base64URLDecode(encoded string) ([]byte, error) {
	if l := len(encoded) % 4; l > 0 {
		encoded += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(encoded)
}

func loadPrivateKeyFromPEMFile(path string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w, failed to read private key from %s", err, path)
	}

	block, _ := pem.Decode(bytes) // Just have 1 block
	if block == nil {
		return nil, fmt.Errorf("%w, failed to parse PEM block", err)
	}

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w, failed to parse private", err)
	}

	rsaKey, ok := pk.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	return rsaKey, nil
}
