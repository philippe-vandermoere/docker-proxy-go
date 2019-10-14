package typeCertificate

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"os"
	"time"
)

type Certificate struct {
	Domain string `validate:"hostname_rfc1123"`
}

func New(domain string) (*Certificate, error) {
	certificate := &Certificate{
		Domain: domain,
	}

	err := certificate.validate()
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func (certificate *Certificate) validate() error {
	validate := validator.New()
	err := validate.Struct(certificate)
	if err != nil {
		errorMessage := "Validate errors:\n"
		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "Domain":
				errorMessage += "Domain '" + certificate.Domain + "' is not valid.\n"
			}
		}

		return errors.New(errorMessage)
	}

	return nil
}

func (certificate *Certificate) GetFileName() string {
	return os.Getenv("CERTIFICATE_DIRECTORY") + "/" + certificate.Domain + "/" + "certificate.pem"
}

func (certificate *Certificate) GetPrivateKeyFileName() string {
	return os.Getenv("CERTIFICATE_DIRECTORY") + "/" + certificate.Domain + "/" + "privateKey.pem"
}

func (certificate *Certificate) GetChainFileName() string {
	return os.Getenv("CERTIFICATE_DIRECTORY") + "/" + certificate.Domain + "/" + "chain.pem"
}

func (certificate *Certificate) HasChain() bool {
	info, err := os.Stat(certificate.GetChainFileName())
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func (certificate *Certificate) Write(certificateContent string, privateKeyContent string, certificateChainContent string) error {
	err := os.MkdirAll(os.Getenv("CERTIFICATE_DIRECTORY")+"/"+certificate.Domain, 755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(certificate.GetFileName(), []byte(certificateContent), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(certificate.GetPrivateKeyFileName(), []byte(privateKeyContent), 0644)
	if err != nil {
		return err
	}

	if certificateChainContent != "" {
		err = ioutil.WriteFile(certificate.GetChainFileName(), []byte(certificateChainContent), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (certificate *Certificate) IsValid() error {
	if _, err := os.Stat(certificate.GetFileName()); os.IsNotExist(err) {
		return err
	}

	if _, err := os.Stat(certificate.GetPrivateKeyFileName()); os.IsNotExist(err) {
		return err
	}

	certPEM, err := ioutil.ReadFile(certificate.GetFileName())
	if err != nil {
		return err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return err
	}

	roots := x509.NewCertPool()
	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return err
	}

	for _, tmp := range certs {
		roots.AddCert(tmp)
	}

	cert := certs[0]

	opts := x509.VerifyOptions{
		DNSName:     certificate.Domain,
		Roots:       roots,
		CurrentTime: time.Now().AddDate(0, 0, 1),
	}

	if _, err := cert.Verify(opts); err != nil {
		return err
	}

	return nil
}
