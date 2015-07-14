package proxy

import (
	"fmt"
	"os"
	"reflect"

	"github.com/andystanton/proxybastard/util"
)

type lookupLocations struct {
	Configs     []string
	Executables []string
}

type lookup struct {
	APM lookupLocations
}

var lookupThing = lookup{
	APM: lookupLocations{
		Configs:     []string{"~/.atom/.apmrc"},
		Executables: []string{"apm"},
	},
}

// Scan scans for proxy targets.
func Scan() {
	reflected := reflect.ValueOf(lookupThing)
	for i := 0; i < reflected.NumField(); i++ {
		fmt.Println(reflected.Type().Field(i).Name)
		lookupLocation, _ := reflected.Field(i).Interface().(lookupLocations)

		for _, config := range lookupLocation.Configs {
			fmt.Printf(" c %s\n", util.SanitisePath(config))
			if _, err := os.Stat(util.SanitisePath(config)); err == nil {
				fmt.Println("FOUND!")
			}
		}
		for _, config := range lookupLocation.Executables {
			fmt.Printf(" x %s\n", config)
		}
	}
}
