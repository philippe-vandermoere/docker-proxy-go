package providers

import (
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate/providers/github"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate/providers/selfsigned"
)

type Provider interface {
	CreateCertificate(certificate *typeCertificate.Certificate) error
	GetName() string
}

func GetProvider(options map[string]string) Provider {
	switch options["name"] {
	case "github":
		return github.New(
			options["repository"],
			options["certificate_path"],
			options["private_key_path"],
			options["certificate_chain_path"],
			options["token"],
			options["reference"],
		)
	}

	return selfSigned.New()
}
