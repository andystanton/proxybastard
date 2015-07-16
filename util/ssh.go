package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

// RunSSHConfiguration defines the configuration necessary to run an ssh command.
type RunSSHConfiguration struct {
	SSHHost string
	SSHPort string
	SSHUser string
	SSHKey  string
}

// RunSSHCommand runs an ssh command.
func RunSSHCommand(runSSHConfiguration RunSSHConfiguration, command string) string {
	keyBytes, err := ioutil.ReadFile(runSSHConfiguration.SSHKey)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: runSSHConfiguration.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", runSSHConfiguration.SSHHost, runSSHConfiguration.SSHPort), config)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %s", fmt.Sprintf("%s:%s", runSSHConfiguration.SSHHost, runSSHConfiguration.SSHPort), err)
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Session failed: %s", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	if err = session.Run(command); err != nil {
		log.Fatalf("Run failed: %s", err)
	}
	return stdoutBuf.String()
}
