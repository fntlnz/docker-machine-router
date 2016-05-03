package machine

import docker "github.com/fsouza/go-dockerclient"

const NETWORK_NAME = "dmr"

func CreateNetwork(client *docker.Client, cidr string) (*docker.Network, error) {
	ipamCfg := []docker.IPAMConfig{
		docker.IPAMConfig{
			Subnet: cidr,
		},
	}
	ipamOpts := docker.IPAMOptions{
		"default",
		ipamCfg,
	}

	netOpts := docker.CreateNetworkOptions{
		NETWORK_NAME,
		false,
		"bridge",
		ipamOpts,
		nil,
		false,
		false,
	}

	n, err := client.NetworkInfo(NETWORK_NAME)

	if err != nil {
		return nil, err
	}

	if n != nil {
		client.RemoveNetwork(n.ID)
	}

	return client.CreateNetwork(netOpts)

}
