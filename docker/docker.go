package docker

import (
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/philippe-vandermoere/docker-proxy-go/types/execute"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Client() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return cli, err
	}

	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
}

func Container(id string) (types.Container, error) {
	var container types.Container
	containers, err := ContainerList()

	if err != nil {
		return container, err
	}

	for _, container := range containers {
		if container.ID == id {
			return container, nil
		}
	}

	return container, errors.New("unable to find container " + id)
}

func ContainerList() ([]types.Container, error) {
	client, err := Client()
	if err != nil {
		return []types.Container{}, err
	}

	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil
}

func ContainerExec(container types.Container, command []string) (typeExecute.Result, error) {
	var executeResult typeExecute.Result
	client, err := Client()
	if err != nil {
		return executeResult, err
	}

	responseCreate, err := client.ContainerExecCreate(
		context.Background(),
		container.ID,
		types.ExecConfig{
			AttachStdin:  false,
			AttachStderr: true,
			AttachStdout: true,
			Tty:          false,
			Cmd:          command,
		},
	)

	if err != nil {
		return executeResult, err
	}

	responseAttach, err := client.ContainerExecAttach(
		context.Background(),
		responseCreate.ID,
		types.ExecStartCheck{},
	)

	if err != nil {
		return executeResult, err
	}

	var stdOutput, stdError bytes.Buffer
	_, err = stdcopy.StdCopy(&stdOutput, &stdError, responseAttach.Reader)
	if err != nil {
		return executeResult, err
	}

	responseInspect, err := client.ContainerExecInspect(context.Background(), responseCreate.ID)
	if err != nil {
		return executeResult, err
	}

	executeResult.StdOutput = stdOutput.String()
	executeResult.StdError = stdError.String()
	executeResult.ExitCode = responseInspect.ExitCode

	return executeResult, nil
}

func NetworkList() ([]types.NetworkResource, error) {
	client, err := Client()
	if err != nil {
		return []types.NetworkResource{}, err
	}

	networks, err := client.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return []types.NetworkResource{}, err
	}

	return networks, nil
}

func NetworkConnect(network types.NetworkResource, container types.Container) (types.Container, error) {
	for _, networkContainer := range container.NetworkSettings.Networks {
		if networkContainer.NetworkID == network.ID {
			return container, nil
		}
	}

	client, err := Client()
	if err != nil {
		return container, err
	}

	err = client.NetworkConnect(
		context.Background(),
		network.ID, container.ID,
		nil,
	)

	if err != nil {
		return container, err
	}

	log.Debug("Connect container '" + strings.Trim(container.Names[0], "/") + "' to network '" + network.Name + "'.")

	return Container(container.ID)
}
