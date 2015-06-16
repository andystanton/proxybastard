package proxy

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/clbanning/mxj"
)

func tildeToUserHome(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return regexp.MustCompile("~").ReplaceAllString(path, usr.HomeDir)
}

func loadXML(filename string) mxj.Map {
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

func loadFileIntoSlice(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read contents of %s\n", filename)
	}
	return strings.Split(string(data), "\n")
}

func writeSliceToFile(filename string, contents []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(contents, "\n")), os.ModeExclusive)
	if err != nil {
		log.Fatalf("Unable to write %s\n%q", filename, err)
	}
}
