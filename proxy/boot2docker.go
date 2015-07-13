package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/andystanton/proxybastard/util"

	"golang.org/x/crypto/ssh"
)

// ExecuteBoot2DockerSSHCommand executes a bootdocker command
func (boot2DockerConfiguration Boot2DockerConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	if boot2DockerConfiguration.Enabled {
		removeFromBoot2DockerProfile(
			boot2DockerConfiguration.SSHHost,
			boot2DockerConfiguration.SSHPort,
			util.SanitisePath(boot2DockerConfiguration.SSHKey))
		addToBoot2DockerProfile(
			boot2DockerConfiguration.SSHHost,
			boot2DockerConfiguration.SSHPort,
			util.SanitisePath(boot2DockerConfiguration.SSHKey),
			proxyHost,
			proxyPort)
		rebootBoot2docker(
			boot2DockerConfiguration.SSHHost,
			boot2DockerConfiguration.SSHPort,
			util.SanitisePath(boot2DockerConfiguration.SSHKey))
	}
}

func (boot2DockerConfiguration Boot2DockerConfiguration) removeProxySettings() {
	if boot2DockerConfiguration.Enabled {
		removeFromBoot2DockerProfile(
			boot2DockerConfiguration.SSHHost,
			boot2DockerConfiguration.SSHPort,
			util.SanitisePath(boot2DockerConfiguration.SSHKey))
		rebootBoot2docker(
			boot2DockerConfiguration.SSHHost,
			boot2DockerConfiguration.SSHPort,
			util.SanitisePath(boot2DockerConfiguration.SSHKey))
	}
}

func rebootBoot2docker(sshHost string, sshPort string, sshKey string) {
	runSSHCommand(sshHost, sshPort, sshKey, "sudo reboot now")
}

func checkBoot2Docker(sshHost string, sshPort string, sshKey string) {
	boot2dockerProfile := "/var/lib/boot2docker/profile"
	log.Println(runSSHCommand(sshHost, sshPort, sshKey, fmt.Sprintf("cat %s", boot2dockerProfile)))
}

func addToBoot2DockerProfile(sshHost string, sshPort string, sshKey string, proxyHost string, proxyPort string) {
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

	runSSHCommand(sshHost, sshPort, sshKey, fmt.Sprintf(addScript, boot2dockerProfile, fmt.Sprintf("%s:%s", proxyHost, proxyPort)))
}

func removeFromBoot2DockerProfile(sshHost string, sshPort string, sshKey string) {
	boot2dockerProfile := "/var/lib/boot2docker/profile"
	removeScript := `
b2d_profile=%s
if [ -f "${b2d_profile}" ]; then
	sudo sed -i '/http\(s\)\{0,1\}_proxy=/d' ${b2d_profile}
fi
`

	runSSHCommand(sshHost, sshPort, sshKey, fmt.Sprintf(removeScript, boot2dockerProfile))
}

func runSSHCommand(sshHost string, sshPort string, sshKey string, command string) string {
	keyBytes, err := ioutil.ReadFile(sshKey)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), config)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %s", fmt.Sprintf("%s:%s", sshHost, sshPort), err)
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Session failed: %s", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	if err != nil {
		log.Fatalf("Run failed: %s", err)
	}
	return stdoutBuf.String()
}
