package execute

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
)

type Result struct {
	StdOutput string
	StdError  string
	ExitCode  int `validate:"gte=0,lte=255"`
}

func New(stdOutput string, stdError string, exitCode int) (*Result, error) {
	result := &Result{
		StdOutput: stdOutput,
		StdError:  stdError,
		ExitCode:  exitCode,
	}

	err := result.validate()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (result *Result) validate() error {
	validate := validator.New()
	err := validate.Struct(result)
	if err != nil {
		errorMessage := "Validate errors:\n"
		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "ExitCode":
				errorMessage += "ExitCode '" + strconv.Itoa(result.ExitCode) + "' must be between 0 and 255.\n"
			}
		}

		return errors.New(errorMessage)
	}

	return nil
}
