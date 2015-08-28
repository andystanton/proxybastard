package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

func (dockerMachineConfiguration DockerMachineConfiguration) validate() error {
	return nil
}

func (dockerMachineConfiguration DockerMachineConfiguration) isEnabled() bool {
	return dockerMachineConfiguration.Enabled
}

func (dockerMachineConfiguration DockerMachineConfiguration) suggestConfiguration() (configuration *Configuration) {
	configuration = nil
	dockerMachineExecutable := "docker-machine"
	_, err := util.ShellOut("which", []string{dockerMachineExecutable})
	if err == nil {
		hosts := listDockerMachines()
		if len(hosts) > 0 {
			configuration = &Configuration{
				Targets: &TargetsConfiguration{
					DockerMachine: &DockerMachineConfiguration{
						Enabled: true,
						Hosts:   hosts,
					},
				},
			}
		}
	}
	return
}

func (dockerMachineConfiguration DockerMachineConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	for _, machine := range listDockerMachines() {
		sshRunConfiguration := dockerMachineConfiguration.extractRunSSHConfiguration(machine)

		removeFromBoot2DockerProfile(sshRunConfiguration)
		addToBoot2DockerProfile(sshRunConfiguration, proxyHost, proxyPort)
		rebootBoot2docker(sshRunConfiguration)
	}
}

func (dockerMachineConfiguration DockerMachineConfiguration) removeProxySettings() {
	for _, machine := range listDockerMachines() {
		sshRunConfiguration := dockerMachineConfiguration.extractRunSSHConfiguration(machine)

		removeFromBoot2DockerProfile(sshRunConfiguration)
		rebootBoot2docker(sshRunConfiguration)
	}
}

func (dockerMachineConfiguration DockerMachineConfiguration) extractRunSSHConfiguration(machine string) util.RunSSHConfiguration {
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
	machines := []string{}
	machineOut, err := util.ShellOut("docker-machine", []string{"ls", "-q", "--filter", "driver=virtualbox"})
	if err != nil {
		log.Fatal(err)
	}
	trimmed := strings.TrimSpace(machineOut)
	if len(trimmed) > 0 {
		machines = strings.Split(strings.TrimSpace(machineOut), "\n")
	}
	return machines
}
