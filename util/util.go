package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/clbanning/mxj"
)

// RunSSHConfiguration defines the configuration necessary to run an ssh command.
type RunSSHConfiguration struct {
	SSHHost string
	SSHPort string
	SSHUser string
	SSHKey  string
}

// SanitisePath translates ~ into $HOME.
func SanitisePath(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return regexp.MustCompile("~").ReplaceAllString(path, usr.HomeDir)
}

// LoadXML loads a file into an xml map.
func LoadXML(filename string) mxj.Map {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	v, err := mxj.NewMapXml(data)
	if err != nil {
		log.Fatal(err)
	}
	return mxj.Map(v)
}

// WriteXML writes xml to a file.
func WriteXML(filename string, xmlMap mxj.Map) {
	output, err := xmlMap.XmlIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filename, output, os.ModeExclusive)
}

// LoadFileIntoSlice loads a file into a string slice.
func LoadFileIntoSlice(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read contents of %s\n", filename)
	}
	return strings.Split(string(data), "\n")
}

// WriteSliceToFile writes a slice to file.
func WriteSliceToFile(filename string, contents []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(contents, "\n")), os.ModeExclusive)
	if err != nil {
		log.Fatalf("Unable to write %s\n%q", filename, err)
	}
}

// ShellOut executes a command.
func ShellOut(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	return string(output), err
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
	err = session.Run(command)
	if err != nil {
		log.Fatalf("Run failed: %s", err)
	}
	return stdoutBuf.String()
}
