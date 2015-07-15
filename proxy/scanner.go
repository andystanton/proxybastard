package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/andystanton/proxybastard/util"
)

var lookupConfiguration = Configuration{
	Targets: &TargetsConfiguration{
		NPM: &NPMConfiguration{
			Enabled: true,
			Files: []string{
				"~/.npmrc",
			},
		},
		APM: &APMConfiguration{
			Enabled: true,
			Files: []string{
				"~/.atom/.apmrc",
			},
		},
	},
}

// Scan scans for proxy targets.
func Scan() {
	var suggestedConfiguration Configuration

	reflected := reflect.Indirect(reflect.ValueOf(lookupConfiguration.Targets))
	if reflect.ValueOf(lookupConfiguration.Targets) != reflect.Zero(reflect.TypeOf(lookupConfiguration.Targets)) {
		for i := 0; i < reflected.NumField(); i++ {

			fieldInterface := reflected.Field(i).Interface()

			if reflect.ValueOf(fieldInterface) != reflect.Zero(reflect.TypeOf(fieldInterface)) {
				reflectedField := reflect.Indirect(reflect.ValueOf(fieldInterface))

				if reflectedField.FieldByName("Files").Kind() != reflect.Invalid {
					files := reflectedField.FieldByName("Files").Interface().([]string)
					for _, file := range files {

						sanitisedFile := util.SanitisePath(file)
						if util.FileExists(sanitisedFile) {

							fmt.Printf("Found %s - adding to suggested config\n", file)
							if suggestedConfiguration.Targets == nil {
								suggestedConfiguration.Targets = &TargetsConfiguration{}
							}
							reflectedTarget := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
							reflectedInterface := reflectedTarget.Field(i).Interface()
							if reflect.ValueOf(reflectedInterface) == reflect.Zero(reflect.TypeOf(reflectedInterface)) {
								newField := reflect.New(reflect.TypeOf(reflectedTarget.Field(i).Interface()).Elem())
								reflect.Indirect(newField).FieldByName("Files").Set(reflect.ValueOf([]string{}))
								reflectedTarget.Field(i).Set(newField)
							}
							files := reflect.Indirect(reflectedTarget.Field(i)).FieldByName("Files").Interface().([]string)
							reflect.Indirect(reflectedTarget.Field(i)).FieldByName("Files").Set(reflect.ValueOf(append(files, file)))
						}
					}
				}
			}
		}
	}

	marshalled, err := json.Marshal(suggestedConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}
