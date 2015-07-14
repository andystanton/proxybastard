package proxy

import (
	"fmt"
	"reflect"
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

func Scan() {
	reflected := reflect.ValueOf(lookupThing)
	for i := 0; i < reflected.NumField(); i++ {
		fmt.Println(reflected.Type().Field(i).Name)
		lookupLocation, _ := reflected.Field(i).Interface().(lookupLocations)
		for _, config := range lookupLocation.Configs {
			fmt.Printf(" c %s\n", config)
		}
		for _, config := range lookupLocation.Executables {
			fmt.Printf(" x %s\n", config)
		}
	}
}
