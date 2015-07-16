package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/deckarep/golang-set"
)

func addToMap(frequencyMap map[string]int, value string) map[string]int {
	if len(strings.TrimSpace(value)) > 0 {
		if _, ok := frequencyMap[value]; !ok {
			frequencyMap[value] = 0
		}
		frequencyMap[value] = frequencyMap[value] + 1
	}
	return frequencyMap
}

// Scan scans for proxy targets.
func Scan() {
	var suggestedConfiguration Configuration

	suggestedProxyHosts := make(map[string]int)
	suggestedProxyPorts := make(map[string]int)
	suggestedSOCKSProxyHosts := make(map[string]int)
	suggestedSOCKSProxyPorts := make(map[string]int)
	suggestedNonProxyHosts := mapset.NewSet()

	targetsField := reflect.Indirect(reflect.ValueOf(&TargetsConfiguration{}))
	for i := 0; i < targetsField.NumField(); i++ {

		configurationField := reflect.New(reflect.TypeOf(targetsField.Field(i).Interface()).Elem()).Interface()
		withConfig, hasConfig := configurationField.(WithConfig)
		if hasConfig {
			fieldName := targetsField.Type().Field(i).Name

			if suggestedItemConfiguration := withConfig.suggestConfiguration(); suggestedItemConfiguration != nil {
				if suggestedConfiguration.Targets == nil {
					suggestedConfiguration.Targets = &TargetsConfiguration{}
				}

				noProtocol := strings.TrimPrefix(strings.TrimPrefix(suggestedItemConfiguration.ProxyHost, "http://"), "https://")
				addToMap(suggestedProxyHosts, noProtocol)
				addToMap(suggestedProxyPorts, suggestedItemConfiguration.ProxyPort)
				addToMap(suggestedSOCKSProxyHosts, suggestedItemConfiguration.SOCKSProxyHost)
				addToMap(suggestedSOCKSProxyPorts, suggestedItemConfiguration.SOCKSProxyPort)

				for _, nonProxyHost := range suggestedItemConfiguration.NonProxyHosts {
					suggestedNonProxyHosts.Add(nonProxyHost)
				}

				targetsField := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
				targetsFieldSuggested := reflect.Indirect(reflect.ValueOf(suggestedItemConfiguration.Targets))

				targetsField.FieldByName(fieldName).Set(targetsFieldSuggested.FieldByName(fieldName))
			}
		}
	}

	marshalled, err := json.MarshalIndent(suggestedConfiguration, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Proxy host candidates: %v\n", suggestedProxyHosts)
	fmt.Printf("Proxy port candidates: %v\n", suggestedProxyPorts)
	fmt.Printf("SOCKS host candidates: %v\n", suggestedSOCKSProxyHosts)
	fmt.Printf("SOCKS port candidates: %v\n", suggestedSOCKSProxyPorts)
	fmt.Printf("Non Proxy Host candidates: %s\n", suggestedNonProxyHosts)
	fmt.Println(string(marshalled))
}
