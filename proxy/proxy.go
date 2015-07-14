package proxy

import (
	"reflect"
	"sync"
)

// EnableProxies enables proxies.
func EnableProxies(config Configuration) {
	toggleProxies(config, true)
}

// DisableProxies disables proxies
func DisableProxies(config Configuration) {
	toggleProxies(config, false)
}

func toggleProxies(config Configuration, enable bool) {
	var wg sync.WaitGroup
	v := reflect.ValueOf(config.Targets)

	for i := 0; i < v.NumField(); i++ {
		configWithProxy, hasProxySettings := v.Field(i).Interface().(WithProxy)

		if hasProxySettings && configWithProxy.isEnabled() {
			wg.Add(1)

			go func(configWithProxy WithProxy, enable bool, proxyHost string, proxyPort string, nonProxyHosts []string) {
				defer wg.Done()
				if enable {
					configWithProxy.addProxySettings(proxyHost, proxyPort, nonProxyHosts)
				} else {
					configWithProxy.removeProxySettings()
				}
			}(configWithProxy, enable, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
		}

		configWithSOCKSProxy, hasSOCKSProxySettings := v.Field(i).Interface().(WithSOCKSProxy)

		if hasSOCKSProxySettings && configWithSOCKSProxy.isEnabled() {
			wg.Add(1)

			go func(configWithSOCKSProxy WithSOCKSProxy, enable bool, socksProxyHost string, socksProxyPort string) {
				defer wg.Done()

				if enable {
					configWithSOCKSProxy.addSOCKSProxySettings(socksProxyHost, socksProxyPort)
				} else {
					configWithSOCKSProxy.removeSOCKSProxySettings()
				}

			}(configWithSOCKSProxy, enable, config.SOCKSProxyHost, config.SOCKSProxyPort)
		}
	}
	wg.Wait()
}
