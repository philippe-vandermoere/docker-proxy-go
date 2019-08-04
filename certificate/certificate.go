package certificate

import (
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate/providers"
	log "github.com/sirupsen/logrus"
)

func GetCertificate(domain string, options map[string]string) (typeCertificate.Certificate, error) {
	certificate, err := typeCertificate.New(domain)
	if err != nil {
		return typeCertificate.Certificate{}, err
	}

	if err = certificate.IsValid(); err == nil {
		log.Info("The certificate is valid for domain '" + certificate.Domain + "'.")
		return certificate, nil
	}

	provider := providers.GetProvider(options)
	err = provider.CreateCertificate(certificate)
	if err != nil {
		log.Error("Unable to created certificate for domain '"+certificate.Domain+"'.\nerror: ", err)
		return typeCertificate.Certificate{}, err
	}

	if err = certificate.IsValid(); err != nil {
		log.Error("The created certificate is not valid for domain '"+certificate.Domain+"'.\nerror:", err)
		return typeCertificate.Certificate{}, nil
	}

	log.Info("Created certificate for domain '" + certificate.Domain + "' with provider '" + provider.GetName() + "'.")

	return certificate, nil
}
