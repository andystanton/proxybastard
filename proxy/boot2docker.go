package proxy

import (
	"fmt"
	"log"

	"github.com/andystanton/proxybastard/util"
)

func (boot2DockerConfiguration Boot2DockerConfiguration) validate() error {
	return nil
}

func (boot2DockerConfiguration Boot2DockerConfiguration) isEnabled() bool {
	return boot2DockerConfiguration.Enabled
}

func (boot2DockerConfiguration Boot2DockerConfiguration) suggestConfiguration() *Configuration {
	return nil
}

func (boot2DockerConfiguration Boot2DockerConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	runSSHConfiguration := boot2DockerConfiguration.extractRunSSHConfiguration()

	removeFromBoot2DockerProfile(runSSHConfiguration)
	addToBoot2DockerProfile(runSSHConfiguration, proxyHost, proxyPort)
	rebootBoot2docker(runSSHConfiguration)
}

func (boot2DockerConfiguration Boot2DockerConfiguration) removeProxySettings() {
	runSSHConfiguration := boot2DockerConfiguration.extractRunSSHConfiguration()

	removeFromBoot2DockerProfile(runSSHConfiguration)
	rebootBoot2docker(runSSHConfiguration)
}

func (boot2DockerConfiguration Boot2DockerConfiguration) extractRunSSHConfiguration() util.RunSSHConfiguration {
	return util.RunSSHConfiguration{
		SSHHost: boot2DockerConfiguration.SSHHost,
		SSHPort: boot2DockerConfiguration.SSHPort,
		SSHUser: "docker",
		SSHKey:  util.SanitisePath(boot2DockerConfiguration.SSHKey),
	}
}

func rebootBoot2docker(runSSHConfiguration util.RunSSHConfiguration) {
	util.RunSSHCommand(runSSHConfiguration, "sudo reboot now")
}

func checkBoot2Docker(runSSHConfiguration util.RunSSHConfiguration) {
	boot2dockerProfile := "/var/lib/boot2docker/profile"
	log.Println(util.RunSSHCommand(runSSHConfiguration, fmt.Sprintf("cat %s", boot2dockerProfile)))
}

func addToBoot2DockerProfile(runSSHConfiguration util.RunSSHConfiguration, proxyHost string, proxyPort string) {
	boot2dockerProfile := "/var/lib/boot2docker/profile"
	addScript := `
b2d_profile=%s
b2d_proxy=%s
if [ ! -f "${b2d_profile}" ]; then
	touch "${b2d_profile}"
fi
sudo sh -c "echo -e \"export http_proxy=${b2d_proxy}\" >>${b2d_profile}"
sudo sh -c "echo -e \"export https_proxy=${b2d_proxy}\" >>${b2d_profile}"
`

	util.RunSSHCommand(runSSHConfiguration, fmt.Sprintf(addScript, boot2dockerProfile, fmt.Sprintf("%s:%s", proxyHost, proxyPort)))
}

func removeFromBoot2DockerProfile(runSSHConfiguration util.RunSSHConfiguration) {
	boot2dockerProfile := "/var/lib/boot2docker/profile"
	removeScript := `
b2d_profile=%s
if [ -f "${b2d_profile}" ]; then
	sudo sed -i '/http\(s\)\{0,1\}_proxy=/d' ${b2d_profile}
fi
`

	util.RunSSHCommand(runSSHConfiguration, fmt.Sprintf(removeScript, boot2dockerProfile))
}
