package execute

import (
	"github.com/Pallinder/go-randomdata"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	stdOutput := randomdata.Alphanumeric(randomdata.Number(0, 4096))
	stdError := randomdata.Alphanumeric(randomdata.Number(0, 4096))
	exitCode := randomdata.Number(0, 255)
	result, err := New(stdOutput, stdError, exitCode)
	if err != nil {
		t.Error(err)
	}

	if result.StdOutput != stdOutput && result.StdError != stdError && result.ExitCode != exitCode {
		t.Fail()
	}
}

func TestNewBadExitCode(t *testing.T) {
	stdOutput := randomdata.Alphanumeric(randomdata.Number(0, 4096))
	stdError := randomdata.Alphanumeric(randomdata.Number(0, 4096))
	var exitCode int
	if randomdata.Boolean() {
		exitCode = randomdata.Number(-4294967295, -1)
	} else {
		exitCode = randomdata.Number(256, 4294967295)
	}

	result, err := New(stdOutput, stdError, exitCode)
	if err == nil {
		t.FailNow()
	}

	if result != nil {
		t.Fail()
	}

	if err.Error() != "Validate errors:\nExitCode '"+strconv.Itoa(exitCode)+"' must be between 0 and 255.\n" {
		t.Fail()
	}
}
