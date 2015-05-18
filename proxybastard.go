package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("bastard.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Unable to open log file")
	}
	defer f.Close()

	log.SetOutput(f)
}
