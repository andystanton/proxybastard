package util

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"

	"github.com/clbanning/mxj"
)

func SanitisePath(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return regexp.MustCompile("~").ReplaceAllString(path, usr.HomeDir)
}

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

func WriteXML(filename string, xmlMap mxj.Map) {
	output, err := xmlMap.XmlIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filename, output, os.ModeExclusive)
}

func LoadFileIntoSlice(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read contents of %s\n", filename)
	}
	return strings.Split(string(data), "\n")
}

func WriteSliceToFile(filename string, contents []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(contents, "\n")), os.ModeExclusive)
	if err != nil {
		log.Fatalf("Unable to write %s\n%q", filename, err)
	}
}

// ShellOut executes a command.
func ShellOut(command string, args []string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}
