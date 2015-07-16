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

	if lookupConfiguration.Targets != nil {

		targetsField := reflect.Indirect(reflect.ValueOf(lookupConfiguration.Targets))

		for i := 0; i < targetsField.NumField(); i++ {

			configurationField := reflect.New(reflect.TypeOf(targetsField.Field(i).Interface()).Elem()).Interface()
			withConfig, hasConfig := configurationField.(WithConfig)
			if hasConfig {
				fieldName := targetsField.Type().Field(i).Name

				if suggestedConfiguration.Targets == nil {
					suggestedConfiguration.Targets = &TargetsConfiguration{}
				}

				targetsField := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
				if withConfig.suggestConfiguration() != nil {
					targetsField.FieldByName(fieldName).Set(reflect.ValueOf(withConfig.suggestConfiguration()))
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
