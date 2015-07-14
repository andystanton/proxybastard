package proxy

import (
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
	v := reflect.ValueOf(config.Targets)

	for i := 0; i < v.NumField(); i++ {
		configWithProxy, hasProxySettings := v.Field(i).Interface().(WithProxy)

		if hasProxySettings && configWithProxy.isEnabled() {
			wg.Add(1)

			go func(configWithProxy WithProxy, mode Mode, proxyHost string, proxyPort string, nonProxyHosts []string) {
				defer wg.Done()
				if mode == Enable {
					configWithProxy.addProxySettings(proxyHost, proxyPort, nonProxyHosts)
				} else {
					configWithProxy.removeProxySettings()
				}
			}(configWithProxy, mode, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
		}

		configWithSOCKSProxy, hasSOCKSProxySettings := v.Field(i).Interface().(WithSOCKSProxy)

		if hasSOCKSProxySettings && configWithSOCKSProxy.isEnabled() {
			wg.Add(1)

			go func(configWithSOCKSProxy WithSOCKSProxy, mode Mode, socksProxyHost string, socksProxyPort string) {
				defer wg.Done()

				if mode == Enable {
					configWithSOCKSProxy.addSOCKSProxySettings(socksProxyHost, socksProxyPort)
				} else {
					configWithSOCKSProxy.removeSOCKSProxySettings()
				}

			}(configWithSOCKSProxy, mode, config.SOCKSProxyHost, config.SOCKSProxyPort)
		}
	}
	wg.Wait()
}
