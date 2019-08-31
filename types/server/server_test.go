package server

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/icrowley/fake"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	name := randomdata.Alphanumeric(randomdata.Number(1, 255))
	ip := fake.IPv4()
	port := randomdata.Number(1, 65535)
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
	ip := fake.IPv4()
	port := randomdata.Number(1, 65535)
	server, err := New(name, ip, port)
	if err == nil {
		t.FailNow()
	}

	if server != nil {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nName is required.\n" {
		t.Fail()
	}
}

func TestNewBadIp(t *testing.T) {
	name := randomdata.Alphanumeric(randomdata.Number(1, 255))
	ip := "1.1.1.1.1"
	port := randomdata.Number(1, 65535)
	server, err := New(name, ip, port)
	if err == nil {
		t.FailNow()
	}

	if server != nil {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nIp '"+ip+"' is not valid.\n" {
		t.Fail()
	}
}

func TestNewBadPort(t *testing.T) {
	name := randomdata.Alphanumeric(randomdata.Number(1, 255))
	ip := fake.IPv4()
	var port int
	if randomdata.Boolean() {
		port = randomdata.Number(65536, 4294967295)
	} else {
		port = randomdata.Number(-4294967295, 0)
	}

	server, err := New(name, ip, port)
	if err == nil {
		t.FailNow()
	}

	if server != nil {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nPort '"+strconv.Itoa(port)+"' must be between 1 and 65535.\n" {
		t.Fail()
	}
}

func TestNewBadNameIpPort(t *testing.T) {
	name := randomdata.Alphanumeric(0)
	ip := "1.1.1.1.1"
	var port int
	if randomdata.Boolean() {
		port = randomdata.Number(65536, 4294967295)
	} else {
		port = randomdata.Number(-4294967295, 0)
	}
	server, err := New(name, ip, port)
	if err == nil {
		t.FailNow()
	}

	if server != nil {
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
