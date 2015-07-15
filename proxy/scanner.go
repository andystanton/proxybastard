package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

var lookupConfiguration = Configuration{
	Targets: &TargetsConfiguration{
		NPM: &NPMConfiguration{
			Enabled: true,
		},
	},
}

// Scan scans for proxy targets.
func Scan() {
	reflected := reflect.Indirect(reflect.ValueOf(lookupConfiguration.Targets))
	var suggestedConfiguration Configuration

	if reflect.ValueOf(lookupConfiguration.Targets) != reflect.Zero(reflect.TypeOf(lookupConfiguration.Targets)) {
		for i := 0; i < reflected.NumField(); i++ {
			fieldInterface := reflected.Field(i).Interface()
			if reflect.ValueOf(fieldInterface) != reflect.Zero(reflect.TypeOf(fieldInterface)) {
				fmt.Println(reflect.Indirect(reflect.ValueOf(fieldInterface)))
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

		}
	}

	marshalled, err := json.Marshal(suggestedConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}
