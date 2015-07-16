package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/deckarep/golang-set"
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

func updateSuggestedConfiguration(suggestedConfiguration *Configuration, configurationToAdd Configuration) {

}

// Scan scans for proxy targets.
func Scan() {
	var suggestedConfiguration Configuration

	suggestedProxyHosts := mapset.NewSet()
	suggestedProxyPorts := mapset.NewSet()
	suggestedSOCKSProxyHosts := mapset.NewSet()
	suggestedSOCKSProxyPorts := mapset.NewSet()
	suggestedNonProxyHosts := mapset.NewSet()

	if lookupConfiguration.Targets != nil {

		targetsField := reflect.Indirect(reflect.ValueOf(lookupConfiguration.Targets))

		for i := 0; i < targetsField.NumField(); i++ {

			configurationField := reflect.New(reflect.TypeOf(targetsField.Field(i).Interface()).Elem()).Interface()
			withConfig, hasConfig := configurationField.(WithConfig)
			if hasConfig {
				fieldName := targetsField.Type().Field(i).Name

				if suggestedItemConfiguration := withConfig.suggestConfiguration(); suggestedItemConfiguration != nil {
					if suggestedConfiguration.Targets == nil {
						suggestedConfiguration.Targets = &TargetsConfiguration{}
					}

					if len(strings.TrimSpace(suggestedItemConfiguration.ProxyHost)) > 0 {
						suggestedProxyHosts.Add(suggestedItemConfiguration.ProxyHost)
					}

					if len(strings.TrimSpace(suggestedItemConfiguration.ProxyPort)) > 0 {
						suggestedProxyPorts.Add(suggestedItemConfiguration.ProxyPort)
					}

					if len(strings.TrimSpace(suggestedItemConfiguration.SOCKSProxyHost)) > 0 {
						suggestedSOCKSProxyHosts.Add(suggestedItemConfiguration.SOCKSProxyHost)
					}

					if len(strings.TrimSpace(suggestedItemConfiguration.SOCKSProxyPort)) > 0 {
						suggestedSOCKSProxyPorts.Add(suggestedItemConfiguration.SOCKSProxyPort)
					}

					for _, nph := range suggestedItemConfiguration.NonProxyHosts {
						suggestedNonProxyHosts.Add(nph)
					}

					targetsField := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
					targetsFieldSuggested := reflect.Indirect(reflect.ValueOf(suggestedItemConfiguration.Targets))

					targetsField.FieldByName(fieldName).Set(targetsFieldSuggested.FieldByName(fieldName))
				}
			}
		}
	}
	marshalled, err := json.MarshalIndent(suggestedConfiguration, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Proxy host candidates: %s\n", suggestedProxyHosts)
	fmt.Printf("Proxy port candidates: %s\n", suggestedProxyPorts)
	fmt.Printf("SOCKS host candidates: %s\n", suggestedSOCKSProxyHosts)
	fmt.Printf("SOCKS port candidates: %s\n", suggestedSOCKSProxyPorts)
	fmt.Printf("Non Proxy Host candidates: %s\n", suggestedNonProxyHosts)
	fmt.Println(string(marshalled))
}
