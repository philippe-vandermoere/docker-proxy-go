package proxy

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/icrowley/fake"
	typeCertificate "github.com/philippe-vandermoere/docker-proxy-go/types/certificate"
	typeServer "github.com/philippe-vandermoere/docker-proxy-go/types/server"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	domain := fake.DomainName()
	proxy, err := New(domain, nil)
	if err != nil {
		t.Error(err)
	}

	if proxy.Domain != domain {
		t.Fail()
	}

	if len(proxy.Servers) != 0 {
		t.Fail()
	}

	if proxy.IsHttps() {
		t.Fail()
	}
}

func TestNewWithCertificate(t *testing.T) {
	domain := fake.DomainName()
	certificate := &typeCertificate.Certificate{Domain: domain}
	proxy, err := New(domain, certificate)
	if err != nil {
		t.Error(err)
	}

	if proxy.Domain != domain {
		t.Fail()
	}

	if !proxy.IsHttps() {
		t.Fail()
	}
}

func TestNewBadCertificate(t *testing.T) {
	domain := fake.DomainName()
	certificateDomain := ""
	var certificate *typeCertificate.Certificate
	if randomdata.Boolean() {
		certificateDomain = fake.DomainName()
		certificate = &typeCertificate.Certificate{Domain: certificateDomain}
	} else {
		certificate = &typeCertificate.Certificate{}
	}

	proxy, err := New(domain, certificate)
	if err == nil {
		t.FailNow()
	}

	if proxy != nil {
		t.Fail()
	}

	if err.Error() != "The domain of certificate '"+certificateDomain+"' must be '"+domain+"'." {
		t.Fail()
	}
}

func TestNewBadDomain(t *testing.T) {
	domain := "aaéééé.gthgfjhg4hgf5h32dfh3df2h5dhf4"
	proxy, err := New(domain, nil)
	if err == nil {
		t.FailNow()
	}

	if proxy != nil {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nDomain '"+domain+"' is not valid.\n" {
		t.Fail()
	}
}

func TestServers(t *testing.T) {
	domain := fake.DomainName()
	proxy, err := New(domain, nil)
	if err != nil {
		t.Error(err)
	}

	if len(proxy.Servers) != 0 {
		t.Fail()
	}

	var paths []string
	for i := 0; i < 5; i++ {
		var path string
		if i == 0 {
			path = "/"
		} else {
			path = "/" + randomdata.Alphanumeric(randomdata.Number(1, 10))
		}

		paths = append(paths, path)

		var servers []*typeServer.Server
		for j := 0; j < 5; j++ {
			server := &typeServer.Server{
				Name: randomdata.Alphanumeric(randomdata.Number(1, 255)),
				Ip:   fake.IPv4(),
				Port: randomdata.Number(1, 65535),
			}

			servers = append(servers, server)
			proxy.AddServer(path, server)
			if proxy.Servers[path][j] != server {
				t.Fail()
			}
		}

		if len(proxy.Servers[path]) != 5 {
			t.Fail()
		}

		if !reflect.DeepEqual(proxy.GetServers(path), servers) {
			t.Log(proxy.GetServers(path), servers)
			t.Fail()
		}
	}

	if len(proxy.Servers) != 5 {
		t.Fail()
	}

	sort.Strings(paths)

	if !reflect.DeepEqual(proxy.GetPaths(), paths) {
		t.Log(proxy.GetPaths(), paths)
		t.Fail()
	}
}

func TestGetUpstream(t *testing.T) {
	domain := fake.DomainName()
	proxy, err := New(domain, nil)
	if err != nil {
		t.Error(err)
	}

	path := "/"
	upstream := strings.ReplaceAll(domain, ".", "_")
	if proxy.GetUpstream(path) != upstream {
		t.Fail()
	}

	path += randomdata.Alphanumeric(randomdata.Number(1, 10))
	upstream += strings.ReplaceAll(path, "/", "_")
	if proxy.GetUpstream(path) != upstream {
		t.Log(proxy.GetUpstream(path), upstream)
		t.Fail()
	}
}

func TestGetHrefHttp(t *testing.T) {
	domain := fake.DomainName()
	proxy, err := New(domain, nil)
	if err != nil {
		t.Error(err)
	}

	path := "/"
	if randomdata.Boolean() {
		path += randomdata.Alphanumeric(randomdata.Number(1, 10))
	}

	for i := 0; i < 10; i++ {
		port := "80"
		if i != 0 {
			port = strconv.Itoa(randomdata.Number(1, 65535))
			err = os.Setenv("HTTP_PORT", port)
			if err != nil {
				t.Error(err)
			}
		}

		href := "http://" + domain + ":" + port + path
		if proxy.GetHref(path) != href {
			t.Log(proxy.GetHref(path), href)
			t.Fail()
		}
	}
}

func TestGetHrefHttps(t *testing.T) {
	domain := fake.DomainName()
	certificate := &typeCertificate.Certificate{Domain: domain}
	proxy, err := New(domain, certificate)
	if err != nil {
		t.Error(err)
	}

	path := "/"
	if randomdata.Boolean() {
		path += randomdata.Alphanumeric(randomdata.Number(1, 10))
	}

	for i := 0; i < 10; i++ {
		port := "443"
		if i != 0 {
			port = strconv.Itoa(randomdata.Number(1, 65535))
			err = os.Setenv("HTTPS_PORT", port)
			if err != nil {
				t.Error(err)
			}
		}

		href := "https://" + domain + ":" + port + path
		if proxy.GetHref(path) != href {
			t.Log(proxy.GetHref(path), href)
			t.Fail()
		}
	}
}
