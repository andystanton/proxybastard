package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

var lookupConfiguration = Configuration{
	Targets: TargetsConfiguration{
		APM: &APMConfiguration{
			Files: []string{"~/.atom/.apmrc"},
		},
		NPM: &NPMConfiguration{
			Files: []string{"~/.npmrc"},
		},
	},
}

// Scan scans for proxy targets.
func Scan() {
	reflected := reflect.ValueOf(lookupConfiguration.Targets)
	for i := 0; i < reflected.NumField(); i++ {

		fmt.Println(reflect.Indirect(reflect.ValueOf(reflected.Field(i))).Interface())
		// reflectedField := reflect.ValueOf(reflected.Field(i).Interface())
		// for j := 0; j < reflectedField.NumField(); j++ {
		//
		// 	fieldName := reflectedField.Type().Field(j).Name
		// 	if fieldName == "Files" {
		//
		// 		files := reflectedField.Field(j).Interface().([]string)
		// 		for _, file := range files {
		//
		// 			sanitisedFile := util.SanitisePath(file)
		// 			fmt.Println(sanitisedFile)
		// 		}
		// 	}
		// }
	}
	marshalled, err := json.Marshal(lookupConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}
