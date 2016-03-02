package host

import "os/exec"

// At the moment I'm doing this using a plain command.
// This could be achieved using the Route apis
// see: https://developer.apple.com/library/mac/documentation/Darwin/Reference/ManPages/man4/route.4.html#//apple_ref/doc/man/4/route

func AllocateHostIp(machineIP string, hostIP string) error {
	err := exec.Command("route", "-n", "add", hostIP, machineIP).Run()

	if err != nil {
		return err
	}
	return nil
}

func DeallocateHostIp(machineIP string, hostIP string) error {
	err := exec.Command("route", "-n", "delete", hostIP, machineIP).Run()

	if err != nil {
		return err
	}
	return nil
}
