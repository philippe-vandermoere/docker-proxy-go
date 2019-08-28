package server

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
)

type Server struct {
	Name string `validate:"required"`
	Ip   string `validate:"ip4_addr"`
	Port int    `validate:"gt=0,lt=65536"`
}

func New(name string, ip string, port int) (Server, error) {
	server := Server{
		Name: name,
		Ip:   ip,
		Port: port,
	}

	err := server.validate()
	if err != nil {
		return Server{}, err
	}

	return server, nil
}

func (server Server) validate() error {
	validate := validator.New()
	err := validate.Struct(server)
	if err != nil {
		errorMessage := "Validate errors:\n"
		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "Name":
				errorMessage += "Name is required.\n"
			case "Ip":
				errorMessage += "Ip '" + server.Ip + "' is not valid.\n"
			case "Port":
				errorMessage += "Port '" + strconv.Itoa(server.Port) + "' must be between 1 and 65535.\n"
			}
		}

		return errors.New(errorMessage)
	}

	return nil
}
