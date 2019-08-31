package proxy

import (
	"errors"
	typeDocker "github.com/docker/docker/api/types"
	"github.com/philippe-vandermoere/docker-proxy-go/certificate"
	"github.com/philippe-vandermoere/docker-proxy-go/client/docker"
	"github.com/philippe-vandermoere/docker-proxy-go/nginx"
	"github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	typeProxy "github.com/philippe-vandermoere/docker-proxy-go/types/proxy"
	typeServer "github.com/philippe-vandermoere/docker-proxy-go/types/server"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

const dockerComposeLabelProject = "com.docker.compose.project"
const dockerComposeLabelService = "com.docker.compose.service"
const dockerProxyLabelProject = "docker-proxy"
const defaultPort = 80
const defaultPath = "/"
const dockerProxyLabelDomain = "com.docker-proxy.domain"
const dockerProxyLabelPath = "com.docker-proxy.path"
const dockerProxyLabelPort = "com.docker-proxy.port"
const dockerProxyLabelSsl = "com.docker-proxy.ssl"
const dockerCertificateProviderPrefix = "com.docker-proxy.certificate-provider"

func Run() error {
	proxyCollection, err := list()
	if err != nil {
		return err
	}

	err = nginx.Proxy(proxyCollection)
	if err != nil {
		return err
	}

	err = applyConfig()
	if err != nil {
		return err
	}

	return nil
}

func getNginxContainer() (typeDocker.Container, error) {
	var container typeDocker.Container
	containers, err := clientDocker.ContainerList()
	if err != nil {
		return container, err
	}

	for _, container := range containers {
		project := container.Labels[dockerComposeLabelProject]
		service := container.Labels[dockerComposeLabelService]
		if project == dockerProxyLabelProject && service == "nginx" {
			return container, nil
		}
	}

	return container, errors.New("unable to find docker proxy nginx container")
}

func getNetwork() (typeDocker.NetworkResource, error) {
	var network typeDocker.NetworkResource
	networks, err := clientDocker.NetworkList()
	if err != nil {
		return network, err
	}

	for _, network := range networks {
		if network.Labels[dockerComposeLabelProject] == dockerProxyLabelProject {
			return network, nil
		}
	}

	return network, errors.New("unable to find docker-proxy Network")
}

func list() (typeProxy.Collection, error) {
	proxyCollection := make(typeProxy.Collection)
	containers, err := clientDocker.ContainerList()
	if err != nil {
		return proxyCollection, err
	}

	for _, container := range containers {
		domain, ok := container.Labels[dockerProxyLabelDomain]
		if ok {
			if _, ok := proxyCollection[domain]; !ok {
				proxyCollection[domain], err = typeProxy.New(domain, getCertificate(domain, container))

				if err != nil {
					return proxyCollection, err
				}
			}

			server, err := getServer(container)
			if err != nil {
				return proxyCollection, err
			}

			proxyCollection[domain].AddServer(getPath(container), server)
		}
	}

	return proxyCollection, nil
}

func getPath(container typeDocker.Container) string {
	path, ok := container.Labels[dockerProxyLabelPath]
	if ok {
		return path
	} else {
		return defaultPath
	}
}

func getServer(container typeDocker.Container) (*typeServer.Server, error) {
	var server *typeServer.Server
	var ip string
	port := defaultPort
	labelPort, ok := container.Labels[dockerProxyLabelPort]
	if ok {
		portInt64, err := strconv.ParseInt(labelPort, 10, 0)
		if err == nil {
			port = int(portInt64)
		}
	}

	proxyNetwork, err := getNetwork()
	if err != nil {
		return server, err
	}

	container, err = clientDocker.NetworkConnect(proxyNetwork, container)
	if err != nil {
		return server, err
	}

	for _, network := range container.NetworkSettings.Networks {
		if proxyNetwork.ID == network.NetworkID {
			ip = network.IPAddress
		}
	}

	server, err = typeServer.New(
		strings.Trim(container.Names[0], "/"),
		ip,
		port,
	)

	if err != nil {
		return server, err
	}

	return server, nil
}

func applyConfig() error {
	nginxContainer, err := getNginxContainer()
	if err != nil {
		return err
	}

	result, err := clientDocker.ContainerExec(nginxContainer, []string{"nginx", "-s", "reload"})
	if err != nil {
		return err
	}

	if result.ExitCode != 0 {
		return errors.New(result.StdError)
	}

	log.Info("Reload nginx configuration.")
	return nil
}

func getCertificate(domain string, container typeDocker.Container) *typeCertificate.Certificate {
	https, ok := container.Labels[dockerProxyLabelSsl]
	if !ok || https != "True" {
		return nil
	}

	options := make(map[string]string)
	regex := regexp.MustCompile(dockerCertificateProviderPrefix)
	for label, value := range container.Labels {
		if regex.MatchString(label) {
			options[strings.TrimPrefix(label, dockerCertificateProviderPrefix+".")] = value
		}
	}

	cert, err := certificate.GetCertificate(domain, options)
	if err != nil {
		return nil
	}

	return cert
}
