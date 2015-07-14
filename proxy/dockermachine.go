package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

func (dockerMachineConfiguration DockerMachineConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	if dockerMachineConfiguration.Enabled {
		for _, machine := range listDockerMachines() {
			sshRunConfiguration := inspectMachine(machine)
			removeFromBoot2DockerProfile(sshRunConfiguration)
			addToBoot2DockerProfile(sshRunConfiguration, proxyHost, proxyPort)
			rebootBoot2docker(sshRunConfiguration)
		}
	}
}

func (dockerMachineConfiguration DockerMachineConfiguration) removeProxySettings() {
	if dockerMachineConfiguration.Enabled {
		for _, machine := range listDockerMachines() {
			sshRunConfiguration := inspectMachine(machine)
			removeFromBoot2DockerProfile(sshRunConfiguration)
			rebootBoot2docker(sshRunConfiguration)
		}
	}
}

func inspectMachine(machine string) util.RunSSHConfiguration {
	var data map[string]interface{}

	inspection, err := util.ShellOut("docker-machine", []string{"inspect", machine})
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(inspection), &data); err != nil {
		log.Fatal(err)
	}

	driverInfo := data["Driver"].(map[string]interface{})

	return util.RunSSHConfiguration{
		SSHHost: driverInfo["IPAddress"].(string),
		SSHPort: "22", // a value is specified in driverInfo["SSHPort"] but 22 appears to be used instead
		SSHUser: driverInfo["SSHUser"].(string),
		SSHKey:  fmt.Sprintf("%s/%s", data["StorePath"].(string), "id_rsa"),
	}
}

func listDockerMachines() []string {
	machineOut, err := util.ShellOut("docker-machine", []string{"ls", "-q", "--filter", "driver=virtualbox"})
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(machineOut), "\n")
}
