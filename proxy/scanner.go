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

func addFilesToSuggestConfiguration(suggestedConfiguration *Configuration, fieldName string, files []string) {
	for _, file := range files {

		sanitisedFile := util.SanitisePath(file)
		if util.FileExists(sanitisedFile) {

			fmt.Printf("Found %s - adding to suggested config\n", file)
			if suggestedConfiguration.Targets == nil {
				suggestedConfiguration.Targets = &TargetsConfiguration{}
			}

			targetsField := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
			configurationFieldPtr := targetsField.FieldByName(fieldName).Interface()

			if util.InterfaceIsZero(configurationFieldPtr) {
				configurationField := reflect.New(reflect.TypeOf(targetsField.FieldByName(fieldName).Interface()).Elem())
				reflect.Indirect(configurationField).FieldByName("Files").Set(reflect.ValueOf([]string{}))
				targetsField.FieldByName(fieldName).Set(configurationField)
			}

			suggestedFiles := reflect.Indirect(targetsField.FieldByName(fieldName)).FieldByName("Files").Interface().([]string)
			reflect.Indirect(targetsField.FieldByName(fieldName)).FieldByName("Files").Set(reflect.ValueOf(append(suggestedFiles, file)))
		}
	}
}

// Scan scans for proxy targets.
func Scan() {
	var suggestedConfiguration Configuration

	if lookupConfiguration.Targets != nil {

		targetsField := reflect.Indirect(reflect.ValueOf(lookupConfiguration.Targets))
		for i := 0; i < targetsField.NumField(); i++ {

			configurationFieldPtr := targetsField.Field(i).Interface()

			if !util.InterfaceIsZero(configurationFieldPtr) {
				configurationField := reflect.Indirect(reflect.ValueOf(configurationFieldPtr))

				if util.ValueHasField(configurationField, "Files") {
					addFilesToSuggestConfiguration(
						&suggestedConfiguration,
						targetsField.Type().Field(i).Name,
						configurationField.FieldByName("Files").Interface().([]string))
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
