package server

import (
	"github.com/Pallinder/go-randomdata"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	name := randomdata.Alphanumeric(10)
	ip := randomdata.IpV4Address()
	port := randomdata.Number(1, 65536)
	server, err := New(name, ip, port)
	if err != nil {
		t.Error(err)
	}

	if server.Name != name && server.Ip != ip && server.Port != port {
		t.Fail()
	}
}

func TestNewBadName(t *testing.T) {
	name := randomdata.Alphanumeric(0)
	ip := randomdata.IpV4Address()
	port := randomdata.Number(1, 65536)
	server, err := New(name, ip, port)
	if err == nil {
		t.Error(err)
	}

	if server != (Server{}) {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nName is required.\n" {
		t.Fail()
	}
}

func TestNewBadIp(t *testing.T) {
	name := randomdata.Alphanumeric(10)
	ip := "1.1.1.1.1"
	port := randomdata.Number(1, 65536)
	server, err := New(name, ip, port)
	if err == nil {
		t.Error(err)
	}

	if server != (Server{}) {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nIp '"+ip+"' is not valid.\n" {
		t.Fail()
	}
}

func TestNewBadPort(t *testing.T) {
	name := randomdata.Alphanumeric(10)
	ip := randomdata.IpV4Address()
	port := 0
	server, err := New(name, ip, port)
	if err == nil {
		t.Error(err)
	}

	if server != (Server{}) {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nPort '"+strconv.Itoa(port)+"' must be between 1 and 65535.\n" {
		t.Fail()
	}
}

func TestNewBadNameIpPort(t *testing.T) {
	name := randomdata.Alphanumeric(0)
	ip := "1.1.1.1.1"
	port := 65536
	server, err := New(name, ip, port)
	if err == nil {
		t.Error(err)
	}

	if server != (Server{}) {
		t.Fail()
	}

	errorMessage := "Validate errors:\n"
	errorMessage += "Name is required.\n"
	errorMessage += "Ip '" + ip + "' is not valid.\n"
	errorMessage += "Port '" + strconv.Itoa(port) + "' must be between 1 and 65535.\n"

	if err.Error() != errorMessage {
		t.Fail()
	}
}
