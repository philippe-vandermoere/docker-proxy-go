package github

import (
	"errors"
	"github.com/philippe-vandermoere/docker-proxy-go/github-client"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	"gopkg.in/go-playground/validator.v9"
)

type Github struct {
	Repository      string `validate:"required"`
	CertificatePath string `validate:"required"`
	PrivateKeyPath  string `validate:"required"`
	Token           string
	Reference       string
}

func New(repository string, certificatePath string, privateKeyPath string, token string, reference string) Github {
	return Github{
		Repository:      repository,
		CertificatePath: certificatePath,
		PrivateKeyPath:  privateKeyPath,
		Reference:       reference,
		Token:           token,
	}
}

func (github Github) Validate() error {
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
			default:
				errorMessage += err.StructField() + "\n"
			}
		}

		return errors.New(errorMessage)
	}

	return nil
}

func (github Github) CreateCertificate(certificate typeCertificate.Certificate) error {
	err := github.Validate()
	if err != nil {
		return err
	}

	certificateContent, err := githubClient.GetFileContent(
		github.Repository,
		github.CertificatePath,
		github.Reference,
		github.Token,
	)
	if err != nil {
		return err
	}

	privateKeyContent, err := githubClient.GetFileContent(
		github.Repository,
		github.PrivateKeyPath,
		github.Reference,
		github.Token,
	)
	if err != nil {
		return err
	}

	err = certificate.Write(certificateContent, privateKeyContent)
	if err != nil {
		return err
	}

	return nil
}

func (github Github) GetName() string {
	return "Github"
}
