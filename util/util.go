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

	"code.google.com/p/go-uuid/uuid"
	"github.com/clbanning/mxj"
	"golang.org/x/crypto/ssh"
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

func writeXML(filename string, xmlMap mxj.Map) {
	output, err := xmlMap.XmlIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filename, output, 0644)
}

// SafeWriteXML writes xml to a random filename in the same
// directory as the target filename and then renames it to the target filename.
func SafeWriteXML(filename string, xmlMap mxj.Map) {
	unique := false
	var safeFilename string
	for !unique {
		safeFilename = generateRandomFilename(filename)
		_, err := os.Stat(safeFilename)
		unique = os.IsNotExist(err)
	}
	writeXML(safeFilename, xmlMap)
	err := os.Rename(safeFilename, filename)
	if err != nil {
		log.Fatal(err)
	}
}

// LoadFileIntoSlice loads a file into a string slice.
func LoadFileIntoSlice(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(data), "\n"), nil
}

func writeSliceToFile(filename string, contents []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(contents, "\n")), 0644)
	if err != nil {
		log.Fatalf("Unable to write %s\n%q", filename, err)
	}
}

// SafeWriteSliceToFile writes a slice to a random filename in the same
// directory as the target filename and then renames it to the target filename.
func SafeWriteSliceToFile(filename string, contents []string) {
	unique := false
	var safeFilename string
	for !unique {
		safeFilename = generateRandomFilename(filename)
		_, err := os.Stat(safeFilename)
		unique = os.IsNotExist(err)
	}
	writeSliceToFile(safeFilename, contents)
	err := os.Rename(safeFilename, filename)
	if err != nil {
		log.Fatal(err)
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

func generateRandomFilename(filepath string) string {
	filepathRegex := regexp.MustCompile("^(.+)/.+$")
	matches := filepathRegex.FindStringSubmatch(filepath)

	newFilename := uuid.New()

	if len(matches) > 0 {
		newFilename = fmt.Sprintf("%s/%s", matches[1], newFilename)
	}

	return newFilename
}
