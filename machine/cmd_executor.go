package machine

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"golang.org/x/crypto/ssh"
)

type MachineCmdExecutor struct {
	IpAddress net.IP
	Port      int
	CertPath  string
}

func NewMachineCmdExecutor(client *docker.Client) (*MachineCmdExecutor, error) {

	certPath, ok := os.LookupEnv("DOCKER_CERT_PATH")

	if ok == false {
		return nil, fmt.Errorf("DOCKER_CERT_PATH environment variable is not set")
	}

	ip, err := MachineIPFromClient(client)

	if err != nil {
		return nil, err
	}

	return &MachineCmdExecutor{
		IpAddress: ip,
		Port:      22,
		CertPath:  certPath,
	}, nil
}

func (m *MachineCmdExecutor) ExecuteCommand(command string) error {
	sshConfig := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			publicKeyFile(fmt.Sprintf("%s/id_rsa", m.CertPath)),
		},
	}

	hostPort := fmt.Sprintf("%s:%d", m.IpAddress.String(), m.Port)

	connection, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		return err
	}

	session, err := connection.NewSession()
	if err != nil {
		return err
	}

	err = session.Run(command)
	if err != nil {
		return err
	}
	return nil
}

func (m *MachineCmdExecutor) PreAllocateIPClass(class string) error {
	var c string
	var i int
	for i = 1; i < 256; i++ {
		c += fmt.Sprintf("sudo ip addr add %s.%d/8 dev eth0;", class, i)
	}
	return m.ExecuteCommand(c)
}

func (m *MachineCmdExecutor) DeAllocateIPClass(class string) error {
	var c string
	var i int
	for i = 1; i < 256; i++ {
		c += fmt.Sprintf("sudo ip addr del %s.%d/8 dev eth0;", class, i)
	}
	return m.ExecuteCommand(c)
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
