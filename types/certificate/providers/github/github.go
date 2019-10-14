package github

import (
	"errors"
	"github.com/philippe-vandermoere/docker-proxy-go/client/github"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	"gopkg.in/go-playground/validator.v9"
)

type Github struct {
	Repository           string `validate:"required"`
	CertificatePath      string `validate:"required"`
	PrivateKeyPath       string `validate:"required"`
	CertificateChainPath string
	Token                string
	Reference            string
}

func New(repository string, certificatePath string, privateKeyPath string, certificateChainPath string, token string, reference string) *Github {
	return &Github{
		Repository:           repository,
		CertificatePath:      certificatePath,
		PrivateKeyPath:       privateKeyPath,
		CertificateChainPath: certificateChainPath,
		Reference:            reference,
		Token:                token,
	}
}

func (github *Github) validate() error {
	validate := validator.New()
	err := validate.Struct(github)
	if err != nil {
		errorMessage := "Validate errors:\n"
		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "Repository":
				errorMessage += "Repository is required.\n"
			case "CertificatePath":
				errorMessage += "CertificatePath is required.\n"
			case "PrivateKeyPath":
				errorMessage += "PrivateKeyPath is required.\n"
			}
		}

		return errors.New(errorMessage)
	}

	return nil
}

func (github *Github) CreateCertificate(certificate *typeCertificate.Certificate) error {
	err := github.validate()
	if err != nil {
		return err
	}

	certificateContent, err := clientGithub.GetFileContent(
		github.Repository,
		github.CertificatePath,
		github.Reference,
		github.Token,
	)
	if err != nil {
		return err
	}

	privateKeyContent, err := clientGithub.GetFileContent(
		github.Repository,
		github.PrivateKeyPath,
		github.Reference,
		github.Token,
	)
	if err != nil {
		return err
	}

	certificateChainContent := ""
	if github.CertificateChainPath != "" {
		certificateChainContent, err = clientGithub.GetFileContent(
			github.Repository,
			github.CertificateChainPath,
			github.Reference,
			github.Token,
		)
		if err != nil {
			return err
		}
	}

	err = certificate.Write(certificateContent, privateKeyContent, certificateChainContent)
	if err != nil {
		return err
	}

	return nil
}

func (github *Github) GetName() string {
	return "Github"
}
