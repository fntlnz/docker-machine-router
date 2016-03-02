package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/docker-machine-router/events"
	"github.com/fntlnz/docker-machine-router/machine"
	docker "github.com/fsouza/go-dockerclient"
)

const IP_CLASS = "10.0.0"

func main() {
	logger := logrus.New()

	if os.Getuid() != 0 {
		logger.Fatal("docker-machine-router must be run as root")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("An error occurred: ", r)
		}
	}()

	client, err := docker.NewClientFromEnv()

	if err != nil {
		logger.Fatal("An error occurred connecting to the Docker daemon")
	}

	err = client.Ping()

	if err != nil {
		logger.Fatalf("Can't reach the docker daemon")
	}

	mc, err := machine.NewMachineCmdExecutor(client)

	if err != nil {
		logger.Fatal("An error occurred creating the docker-machine command executor: %s", err.Error())
	}

	logger.Info("GC IPs on docker machine")
	err = mc.DeAllocateIPClass(IP_CLASS)

	if err != nil {
		logger.Warn("Problem occurred GC IPs on docker machine")
	}

	logger.Info("Preallocating IPs on docker machine")
	err = mc.PreAllocateIPClass(IP_CLASS)

	if err != nil {
		logger.Warnf("Problem occurred preallocating IPs on docker machine")
	}

	listener := events.NewDockerListener(client, logger)
	err = listener.Listen()

	if err != nil {
		logger.Fatal("Error occurred attaching to the Docker daemon event stream: ", err.Error())
	}
}
