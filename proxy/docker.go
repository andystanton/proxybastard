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
func ExecuteBoot2DockerSSHCommand() {

	checkBoot2Docker()
}

func rebootBoot2docker() {
	doSSHCommand("sudo reboot now")
}

func checkBoot2Docker() {
	boot2dockerProfile := "/var/lib/boot2docker/profile"
	log.Println(doSSHCommand(fmt.Sprintf("cat %s", boot2dockerProfile)))
}

func addToBoot2Docker(proxyHost string, proxyPort string) {
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

	doSSHCommand(fmt.Sprintf(addScript, boot2dockerProfile, fmt.Sprintf("%s:%s", proxyHost, proxyPort)))
}

func removeFromBoot2Docker() {
	boot2dockerProfile := "/var/lib/boot2docker/profile"

	removeScript := `
b2d_profile=%s
if [ -f "${b2d_profile}" ]; then
	sudo sed -i '/http\(s\)\{0,1\}_proxy=/d' ${b2d_profile}
fi
`

	doSSHCommand(fmt.Sprintf(removeScript, boot2dockerProfile))
}

func doSSHCommand(command string) string {
	sshKeyFile := "~/.ssh/id_boot2docker"
	boot2dockerIP := "192.168.59.103"
	boot2dockerSSHPort := 22

	keyBytes, err := ioutil.ReadFile(util.SanitisePath(sshKeyFile))
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

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", boot2dockerIP, boot2dockerSSHPort), config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("session failed:%v", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	if err != nil {
		log.Fatalf("Run failed:%v", err)
	}
	return stdoutBuf.String()
}
