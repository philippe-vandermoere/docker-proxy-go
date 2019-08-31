package proxy

import (
	"errors"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	typeServer "github.com/philippe-vandermoere/docker-proxy-go/types/server"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"sort"
	"strings"
)

type Collection map[string]*Proxy

type Proxy struct {
	Domain      string `validate:"hostname_rfc1123"`
	Servers     map[string][]*typeServer.Server
	Certificate *typeCertificate.Certificate
}

func New(domain string, certificate *typeCertificate.Certificate) (*Proxy, error) {
	proxy := &Proxy{
		Domain:      domain,
		Servers:     make(map[string][]*typeServer.Server),
		Certificate: certificate,
	}

	if certificate != nil && certificate.Domain != proxy.Domain {
		return nil, errors.New("The domain of certificate '" + certificate.Domain + "' must be '" + proxy.Domain + "'.")
	}

	err := proxy.validate()
	if err != nil {
		return nil, err
	}

	return proxy, nil
}

func (proxy *Proxy) validate() error {
	validate := validator.New()
	err := validate.Struct(proxy)
	if err != nil {
		errorMessage := "Validate errors:\n"
		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "Domain":
				errorMessage += "Domain '" + proxy.Domain + "' is not valid.\n"
			}
		}

		return errors.New(errorMessage)
	}

	return nil
}

func (proxy *Proxy) AddServer(path string, server *typeServer.Server) {
	proxy.Servers[path] = append(proxy.Servers[path], server)
}

func (proxy *Proxy) GetServers(path string) []*typeServer.Server {
	return proxy.Servers[path]
}

func (proxy *Proxy) GetPaths() []string {
	var paths []string

	for path := range proxy.Servers {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	return paths
}

func (proxy *Proxy) IsHttps() bool {
	if proxy.Certificate == nil {
		return false
	} else {
		return true
	}
}

// for template
func (proxy *Proxy) GetUpstream(path string) string {
	upstream := strings.ReplaceAll(proxy.Domain, ".", "_")
	if path != "/" {
		upstream += strings.ReplaceAll(path, "/", "_")
	}

	return upstream
}

// for template
func (proxy *Proxy) GetHref(path string) string {
	var href string
	var port string

	if proxy.IsHttps() {
		href = "https://"
		port = os.Getenv("HTTPS_PORT")
		if port == "" {
			port = "443"
		}
	} else {
		href = "http://"
		port = os.Getenv("HTTP_PORT")
		if port == "" {
			port = "80"
		}
	}

	return href + proxy.Domain + ":" + port + path
}
