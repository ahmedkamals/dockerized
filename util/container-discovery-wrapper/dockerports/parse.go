package dockerports

import (
	"os"
	"github.com/fsouza/go-dockerclient"
	"errors"
)

func ParseAndSetenv() {
	if _, err := os.Stat("/docker.sock"); os.IsNotExist(err) {
		panic(errors.New("Couldn't find Docker socket inside the container, you need to mount HOST:/var/run/docker.sock to CONTAINER:/docker.sock"))
	}

	containerId := os.Getenv("HOSTNAME")

	client, err := docker.NewClient("unix:///docker.sock")

	if err != nil {
		panic(err)
	}

	container, err := client.InspectContainer(containerId)

	if err != nil {
		panic(err)
	}

	for containerPort, portBindings := range container.NetworkSettings.Ports {
		for _, portBinding := range portBindings {
			//todo: take hostBinding.HostIP into consideration if available
			os.Setenv("DOCKER_CONTAINER_PORT_"+containerPort.Port()+"_REAL", portBinding.HostPort)
		}
	}
}
