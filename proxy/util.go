package proxy

import (
	"io/ioutil"
	"log"
	"strings"
)

func loadFileIntoSlice(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read contents of %s\n", filename)
	}
	return strings.Split(string(data), "\n")
}
