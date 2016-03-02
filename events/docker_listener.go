package events

import (
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/docker-machine-router/host"
	"github.com/fntlnz/docker-machine-router/machine"
	docker "github.com/fsouza/go-dockerclient"
)

type DockerListener struct {
	client *docker.Client
	logger *logrus.Logger
}

func NewDockerListener(c *docker.Client, l *logrus.Logger) *DockerListener {
	return &DockerListener{
		client: c,
		logger: l,
	}
}

func (l *DockerListener) Listen() error {
	machineIP, err := machine.MachineIPFromClient(l.client)
	if err != nil {
		return fmt.Errorf("Can't get docker machine IP")
	}

	eventsChannel := make(chan *docker.APIEvents)

	l.logger.Debugf("Listening for docker daemon events")
	err = l.client.AddEventListener(eventsChannel)

	if err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-eventsChannel:

			l.logger.Infof("[EVENT] id: %s\tStatus:%s\tDate: %s", event.ID, event.Status, time.Now().Format("15:04:05"))

			if event.Status == "create" || event.Status == "start" {
				info, err := l.client.InspectContainer(event.ID)
				if err != nil {
					l.logger.Warnf("Error inspecting container: %s", event.ID)
					continue
				}
				// Deallocation should not be there.
				l.deAllocateHostIPsByPortBindings(machineIP, info.HostConfig.PortBindings)
				l.allocateHostIPsByPortBindings(machineIP, info.HostConfig.PortBindings)
			}

			if !ok {
				eventsChannel = nil
			}

			if eventsChannel == nil {
				l.logger.Warn("eventsChannel nil")
				break
			}
		}
	}
}

func (l *DockerListener) allocateHostIPsByPortBindings(machineIP net.IP, bindings map[docker.Port][]docker.PortBinding) {
	for _, binding := range bindings {
		for _, setting := range binding {
			l.logger.Infof("Allocating ip: %s", setting.HostIP)
			err := host.AllocateHostIp(machineIP.String(), setting.HostIP)
			if err != nil {
				l.logger.Warnf("Error allocating ip: %s", setting.HostIP)
			}
		}
	}
}

func (l *DockerListener) deAllocateHostIPsByPortBindings(machineIP net.IP, bindings map[docker.Port][]docker.PortBinding) {
	for _, binding := range bindings {
		for _, setting := range binding {
			l.logger.Infof("Deallocating ip: %s", setting.HostIP)
			err := host.DeallocateHostIp(machineIP.String(), setting.HostIP)
			if err != nil {
				l.logger.Warnf("Error deallocating ip: %s", setting.HostIP)
			}
		}
	}
}
