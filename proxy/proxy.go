package proxy

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/andystanton/proxybastard/util"
)

// Mode is the application mode: Enable|Disable
type Mode int

const (
	// Enable mode.
	Enable Mode = iota
	// Disable mode.
	Disable
)

// ToggleProxies toggles proxy settings.
func ToggleProxies(config Configuration, mode Mode) {

	if config.Targets != nil {
		var wg sync.WaitGroup

		targetsField := reflect.Indirect(reflect.ValueOf(config.Targets))
		for i := 0; i < targetsField.NumField(); i++ {
			configurationFieldPtr := targetsField.Field(i).Interface()

			if !util.InterfaceIsZero(configurationFieldPtr) {

				configItem, hasConfig := configurationFieldPtr.(WithConfig)
				if hasConfig && configItem.isEnabled() {
					if err := configItem.validate(); err != nil {
						log.Fatal(err)
					}

					if configWithProxy, hasProxySettings := configItem.(WithProxy); hasProxySettings {
						wg.Add(1)

						go func(
							configName string,
							configWithProxy WithProxy,
							mode Mode,
							proxyHost string,
							proxyPort string,
							nonProxyHosts []string) {

							defer wg.Done()

							if mode == Enable {
								configWithProxy.addProxySettings(proxyHost, proxyPort, nonProxyHosts)
								fmt.Printf("Enabled http proxy settings for %s\n", configName)
							} else {
								configWithProxy.removeProxySettings()
								fmt.Printf("Disabled http proxy settings for %s\n", configName)
							}

						}(
							targetsField.Type().Field(i).Name,
							configWithProxy,
							mode,
							config.ProxyHost,
							config.ProxyPort,
							config.NonProxyHosts)
					}

					if configWithSOCKSProxy, hasSOCKSProxySettings := configItem.(WithSOCKSProxy); hasSOCKSProxySettings {
						wg.Add(1)

						go func(
							configName string,
							configWithSOCKSProxy WithSOCKSProxy,
							mode Mode,
							socksProxyHost string,
							socksProxyPort string) {

							defer wg.Done()

							if mode == Enable {
								configWithSOCKSProxy.addSOCKSProxySettings(socksProxyHost, socksProxyPort)
								fmt.Printf("Enabled SOCKS proxy settings for %s\n", configName)
							} else {
								configWithSOCKSProxy.removeSOCKSProxySettings()
								fmt.Printf("Disabled SOCKS proxy settings for %s\n", configName)
							}
						}(
							targetsField.Type().Field(i).Name,
							configWithSOCKSProxy,
							mode,
							config.SOCKSProxyHost,
							config.SOCKSProxyPort)
					}
				}
			}
		}
		wg.Wait()
	}
}
