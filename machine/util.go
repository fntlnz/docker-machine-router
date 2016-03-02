package machine

import (
	"net"
	"net/url"

	"github.com/fsouza/go-dockerclient"
)

func MachineIPFromClient(client *docker.Client) (net.IP, error) {
	parsedIP, err := url.Parse(client.Endpoint())

	if err != nil {
		return nil, err
	}

	ipString, _, err := net.SplitHostPort(parsedIP.Host)

	if err != nil {
		return nil, err
	}

	return net.ParseIP(ipString), nil
}
