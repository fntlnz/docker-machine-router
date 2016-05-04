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
		Name:           NETWORK_NAME,
		CheckDuplicate: false,
		Driver:         "bridge",
		IPAM:           ipamOpts,
		Options:        nil,
		Internal:       false,
		EnableIPv6:     false,
	}

	n, _ := client.NetworkInfo(NETWORK_NAME)

	if n != nil {
		client.RemoveNetwork(n.ID)
	}

	return client.CreateNetwork(netOpts)

}
