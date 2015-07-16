package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/clbanning/mxj"
)

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

// WriteXML writes an xml map to a file.
func WriteXML(filename string, xmlMap mxj.Map) {
	output, err := xmlMap.XmlIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}
	safeWrite(filename, output)
}

// LoadFileIntoSlice loads a file into a string slice.
func LoadFileIntoSlice(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(data), "\n"), nil
}

// WriteSliceToFile writes a string slice to a file.
func WriteSliceToFile(filename string, contents []string) {
	safeWrite(filename, []byte(strings.Join(contents, "\n")))
}

// FileExists returns whether or not a file exists.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil && !os.IsNotExist(err)
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

func safeWrite(filename string, content []byte) {
	unique := false
	var safeFilename string
	for !unique {
		safeFilename = generateRandomFilename(filename)
		_, err := os.Stat(safeFilename)
		unique = os.IsNotExist(err)
	}

	if err := ioutil.WriteFile(safeFilename, content, 0644); err != nil {
		log.Fatalf("Unable to write %s\n%q", safeFilename, err)
	}

	if err := os.Rename(safeFilename, filename); err != nil {
		log.Fatal(err)
	}
}
