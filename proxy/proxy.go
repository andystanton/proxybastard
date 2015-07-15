package proxy

import (
	"fmt"
	"log"
	"reflect"
	"sync"
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
	var wg sync.WaitGroup

	reflected := reflect.Indirect(reflect.ValueOf(config.Targets))
	if reflect.ValueOf(config.Targets) != reflect.Zero(reflect.TypeOf(config.Targets)) {
		for i := 0; i < reflected.NumField(); i++ {
			fieldInterface := reflected.Field(i).Interface()

			if reflect.ValueOf(fieldInterface) != reflect.Zero(reflect.TypeOf(fieldInterface)) {
				configWithProxy, hasProxySettings := fieldInterface.(WithProxy)

				if hasProxySettings && configWithProxy.isEnabled() {
					wg.Add(1)

					go func(
						configName string,
						configWithProxy WithProxy,
						mode Mode,
						proxyHost string,
						proxyPort string,
						nonProxyHosts []string) {

						defer wg.Done()
						if err := configWithProxy.validate(); err != nil {
							log.Fatal(err)
						}

						if mode == Enable {
							configWithProxy.addProxySettings(proxyHost, proxyPort, nonProxyHosts)
							fmt.Printf("Enabled http proxy settings for %s\n", configName)
						} else {
							configWithProxy.removeProxySettings()
							fmt.Printf("Disabled http proxy settings for %s\n", configName)
						}

					}(
						reflected.Type().Field(i).Name,
						configWithProxy,
						mode,
						config.ProxyHost,
						config.ProxyPort,
						config.NonProxyHosts)
				}

				configWithSOCKSProxy, hasSOCKSProxySettings := reflected.Field(i).Interface().(WithSOCKSProxy)

				if hasSOCKSProxySettings && configWithSOCKSProxy.isEnabled() {
					wg.Add(1)

					go func(
						configName string,
						configWithSOCKSProxy WithSOCKSProxy,
						mode Mode,
						socksProxyHost string,
						socksProxyPort string) {

						defer wg.Done()
						if err := configWithSOCKSProxy.validate(); err != nil {
							log.Fatal(err)
						}

						if mode == Enable {
							configWithSOCKSProxy.addSOCKSProxySettings(socksProxyHost, socksProxyPort)
							fmt.Printf("Enabled SOCKS proxy settings for %s\n", configName)
						} else {
							configWithSOCKSProxy.removeSOCKSProxySettings()
							fmt.Printf("Disabled SOCKS proxy settings for %s\n", configName)
						}
					}(
						reflected.Type().Field(i).Name,
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
