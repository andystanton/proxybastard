package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

type dockerMachineSSHConfiguration struct {
	sshHost string
	sshPort string
	sshUser string
	sshKey  string
}

func (dockerMachineConfiguration DockerMachineConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	if dockerMachineConfiguration.Enabled {
		for _, machine := range listDockerMachines() {
			sshConfig := inspectMachine(machine)
			removeFromBoot2DockerProfile(sshConfig.sshHost, sshConfig.sshPort, sshConfig.sshUser, sshConfig.sshKey)
			addToBoot2DockerProfile(sshConfig.sshHost, sshConfig.sshPort, sshConfig.sshUser, sshConfig.sshKey, proxyHost, proxyPort)
			rebootBoot2docker(sshConfig.sshHost, sshConfig.sshPort, sshConfig.sshUser, sshConfig.sshKey)
		}
	}
}

func (dockerMachineConfiguration DockerMachineConfiguration) removeProxySettings() {
	if dockerMachineConfiguration.Enabled {
		for _, machine := range listDockerMachines() {
			sshConfig := inspectMachine(machine)
			removeFromBoot2DockerProfile(sshConfig.sshHost, sshConfig.sshPort, sshConfig.sshUser, sshConfig.sshKey)
			rebootBoot2docker(sshConfig.sshHost, sshConfig.sshPort, sshConfig.sshUser, sshConfig.sshKey)
		}
	}
}

func inspectMachine(machine string) dockerMachineSSHConfiguration {
	var data map[string]interface{}

	inspection, err := util.ShellOut("docker-machine", []string{"inspect", machine})
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(inspection), &data); err != nil {
		log.Fatal(err)
	}
	var driverInfo map[string]interface{}
	driverInfo = data["Driver"].(map[string]interface{})

	sshHost := driverInfo["IPAddress"].(string)
	storePath := data["StorePath"].(string)
	sshUser := driverInfo["SSHUser"].(string)
	// a value is specified here but 22 appears to be used instead
	//sshPort := driverInfo["SSHPort"]
	sshPort := "22"
	sshKey := fmt.Sprintf("%s/%s", storePath, "id_rsa")

	fmt.Printf("SSH Host: %s\n", sshHost)
	fmt.Printf("SSH User: %s\n", sshUser)
	fmt.Printf("SSH Port: %s\n", sshPort)
	fmt.Printf("Store Path: %s\n", storePath)

	fmt.Printf("SSH Key: %s/%s\n", storePath, "id_rsa")

	return dockerMachineSSHConfiguration{
		sshHost: sshHost,
		sshPort: sshPort,
		sshUser: sshUser,
		sshKey:  sshKey,
	}
}

func listDockerMachines() []string {
	machineOut, err := util.ShellOut("docker-machine", []string{"ls", "-q", "--filter", "driver=virtualbox"})
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(machineOut), "\n")
}
