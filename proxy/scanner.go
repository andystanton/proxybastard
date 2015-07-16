package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/deckarep/golang-set"
)

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

				if len(strings.TrimSpace(suggestedItemConfiguration.ProxyHost)) > 0 {
					noProtocol := strings.TrimPrefix(strings.TrimPrefix(suggestedItemConfiguration.ProxyHost, "http://"), "https://")
					if _, ok := suggestedProxyHosts[noProtocol]; !ok {
						suggestedProxyHosts[noProtocol] = 0
					}
					suggestedProxyHosts[noProtocol] = suggestedProxyHosts[noProtocol] + 1
				}

				if len(strings.TrimSpace(suggestedItemConfiguration.ProxyPort)) > 0 {
					if _, ok := suggestedProxyPorts[suggestedItemConfiguration.ProxyPort]; !ok {
						suggestedProxyPorts[suggestedItemConfiguration.ProxyPort] = 0
					}
					suggestedProxyPorts[suggestedItemConfiguration.ProxyPort] = suggestedProxyPorts[suggestedItemConfiguration.ProxyPort] + 1
				}

				if len(strings.TrimSpace(suggestedItemConfiguration.SOCKSProxyHost)) > 0 {
					if _, ok := suggestedSOCKSProxyHosts[suggestedItemConfiguration.SOCKSProxyHost]; !ok {
						suggestedSOCKSProxyHosts[suggestedItemConfiguration.SOCKSProxyHost] = 0
					}
					suggestedSOCKSProxyHosts[suggestedItemConfiguration.SOCKSProxyHost] = suggestedSOCKSProxyHosts[suggestedItemConfiguration.SOCKSProxyHost] + 1
				}

				if len(strings.TrimSpace(suggestedItemConfiguration.SOCKSProxyPort)) > 0 {
					if _, ok := suggestedSOCKSProxyPorts[suggestedItemConfiguration.SOCKSProxyPort]; !ok {
						suggestedSOCKSProxyPorts[suggestedItemConfiguration.SOCKSProxyPort] = 0
					}
					suggestedSOCKSProxyPorts[suggestedItemConfiguration.SOCKSProxyPort] = suggestedSOCKSProxyPorts[suggestedItemConfiguration.SOCKSProxyPort] + 1
				}

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
