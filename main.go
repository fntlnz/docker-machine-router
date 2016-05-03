package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fntlnz/docker-machine-router/host"
	"github.com/fntlnz/docker-machine-router/machine"
	docker "github.com/fsouza/go-dockerclient"
)

var (
	cidr string
)

const (
	BANNER  = "DOCKER MACHINE ROUTER v%s"
	version = "0.2.0"
)

func banner() string {
	return fmt.Sprintf(BANNER, version)
}
func usage() {
	fmt.Fprintf(os.Stderr, banner(), version)
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&cidr, "cidr", "10.18.0.0/16", "Used to identify the IP address range that can be allocated by docker-machine-router")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	if os.Getuid() != 0 {
		fmt.Println("❗  docker-machine-router must be run as root")
		os.Exit(1)
	}

	fmt.Println(banner())

	client, err := docker.NewClientFromEnv()

	if err != nil {
		fmt.Println("❗  An error occurred contacting the Docker daemon: ", err.Error())
		os.Exit(1)
	}

	err = client.Ping()

	if err != nil {
		fmt.Println("❗  Can't reach the docker daemon")
		os.Exit(1)
	}

	fmt.Println("Connected to: ", os.Getenv("DOCKER_HOST"))
	_, err = machine.CreateNetwork(client, cidr)

	if err != nil {
		fmt.Println("❗  An error occurred while creating the Docker Network: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("✅  Created network on Docker machine")

	machineIP, err := machine.MachineIPFromClient(client)

	if err != nil {
		fmt.Println("❗  An error occurred while trying to obtain the Docker Machine ip address: ", err.Error())
		os.Exit(1)
	}

	err = host.DeallocateHostIp(machineIP.String(), cidr)

	if err != nil {
		fmt.Println("❗  An error occurred while garbage collecting the route from the host: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("✅  Removed previous allocated routes")
	err = host.AllocateHostIp(machineIP.String(), cidr)

	if err != nil {
		fmt.Println("❗  An error occurred while creating the route on the host: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("✅  Created new routes")

	fmt.Println("🎉  The network has been setup 🎉")
}
