package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/docker-machine-router/host"
	"github.com/fntlnz/docker-machine-router/machine"
	"github.com/fntlnz/docker-machine-router/network"
	docker "github.com/fsouza/go-dockerclient"
)

var (
	cidr  string
	debug bool
)

const (
	BANNER = `
DOCKER MACHINE ROUTER
This tool allows you to reach the container's internal ip addresses from the host by routing the OS X host traffic trough the Docker Machine VM.
(c) 2016 Lorenzo Fontana
version: %s
`
	version = "0.2.0"
)

func usage() {
	fmt.Fprintf(os.Stderr, BANNER, version)
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&cidr, "cidr", "10.18.0.0/16", "Used to identify the IP address range that can be allocated by docker-machine-router")
	flag.BoolVar(&debug, "debug", false, "Start in debug mode, provides a lot more information")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	logger := logrus.New()

	if os.Getuid() != 0 {
		logger.Fatal("docker-machine-router must be run as root")
	}

	client, err := docker.NewClientFromEnv()

	if err != nil {
		logger.Fatal("An error occurred contacting the Docker daemon: ", err.Error())
	}

	err = client.Ping()

	if err != nil {
		logger.Fatal("Can't reach the docker daemon")
	}

	_, err = network.CreateNetwork(client, cidr)

	if err != nil {
		logger.Fatal("An error occurred while creating the Docker Network: ", err.Error())
	}

	machineIP, err := machine.MachineIPFromClient(client)

	if err != nil {
		logger.Fatal("An error occurred while trying to obtain the Docker Machine ip address: ", err.Error())
	}

	host.DeallocateHostIp(machineIP.String(), cidr)
	host.AllocateHostIp(machineIP.String(), cidr)
}
