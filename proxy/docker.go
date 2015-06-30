package proxy

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/andystanton/proxybastard/util"

	"golang.org/x/crypto/ssh"
)

// ExecuteBoot2DockerSSHCommand executes a bootdocker command
func ExecuteBoot2DockerSSHCommand() {
	sshKeyFile := "~/.ssh/id_boot2docker"

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

	conn, err := ssh.Dial("tcp", "192.168.59.103:22", config)
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
	err = session.Run("whoami")
	if err != nil {
		log.Fatalf("Run failed:%v", err)
	}
	log.Printf("%s", stdoutBuf.String())
}
