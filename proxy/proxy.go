package proxy

import (
	"fmt"
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

	reflected := reflect.ValueOf(config.Targets)
	for i := 0; i < reflected.NumField(); i++ {
		configWithProxy, hasProxySettings := reflected.Field(i).Interface().(WithProxy)

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
				if mode == Enable {
					fmt.Printf("Enabling http proxy settings for %s\n", configName)

					configWithProxy.addProxySettings(proxyHost, proxyPort, nonProxyHosts)
				} else {
					fmt.Printf("Disabling http proxy settings for %s\n", configName)

					configWithProxy.removeProxySettings()
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

				if mode == Enable {
					fmt.Printf("Enabling SOCKS proxy settings for %s\n", configName)

					configWithSOCKSProxy.addSOCKSProxySettings(socksProxyHost, socksProxyPort)
				} else {
					fmt.Printf("Disabling SOCKS proxy settings for %s\n", configName)

					configWithSOCKSProxy.removeSOCKSProxySettings()
				}
			}(
				reflected.Type().Field(i).Name,
				configWithSOCKSProxy,
				mode,
				config.SOCKSProxyHost,
				config.SOCKSProxyPort)
		}
	}
	wg.Wait()
}
