package selfSigned

import (
	"errors"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
)

type SelfSigned struct {
}

func New() SelfSigned {
	return SelfSigned{}
}

func (selfSigned SelfSigned) CreateCertificate(certificate typeCertificate.Certificate) error {
	return errors.New("selfSigned not implemented")
}

func (selfSigned SelfSigned) GetName() string {
	return "SelfSigned"
}
