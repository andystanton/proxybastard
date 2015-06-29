package proxy

import (
	"log"

	"golang.org/x/crypto/ssh"
)

// ExecuteBoot2DockerSSHCommand executes a bootdocker command
func ExecuteBoot2DockerSSHCommand() {
	config := &ssh.ClientConfig{
		User: "username",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()
}
