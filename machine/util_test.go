package machine

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func TestGetIpFromClient(t *testing.T) {
	endpoint := "tcp://104.236.41.100:2376"
	m, _ := docker.NewClient(endpoint)
	d, _ := MachineIPFromClient(m)
	if d.String() != "104.236.41.100" {
		t.Fail()
	}
}
