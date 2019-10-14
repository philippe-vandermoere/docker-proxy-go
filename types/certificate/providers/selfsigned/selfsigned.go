package selfSigned

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	"math/big"
	"time"
)

type SelfSigned struct {
}

func New() *SelfSigned {
	return &SelfSigned{}
}

func (selfSigned *SelfSigned) CreateCertificate(certificate *typeCertificate.Certificate) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1337),
		Issuer: pkix.Name{
			CommonName: certificate.Domain,
		},
		DNSNames:              []string{certificate.Domain},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 365),
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, ca, ca, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	certificateContent := &bytes.Buffer{}
	err = pem.Encode(certificateContent, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		return err
	}

	privateKeyContent := &bytes.Buffer{}
	err = pem.Encode(privateKeyContent, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return err
	}

	err = certificate.Write(certificateContent.String(), privateKeyContent.String(), "")
	if err != nil {
		return err
	}

	return nil
}

func (selfSigned *SelfSigned) GetName() string {
	return "self-signed"
}
