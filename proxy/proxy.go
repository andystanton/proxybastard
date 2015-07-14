package proxy

import (
	"reflect"
	"sync"
)

// EnableProxies enables proxies.
func EnableProxies(config Configuration) {
	var wg sync.WaitGroup
	v := reflect.ValueOf(config.Targets)

	for i := 0; i < v.NumField(); i++ {
		configWithProxy, hasProxySettings := v.Field(i).Interface().(WithProxy)

		if hasProxySettings {
			wg.Add(1)

			go func(configWithProxy WithProxy, proxyHost string, proxyPort string, nonProxyHosts []string) {
				defer wg.Done()

				configWithProxy.addProxySettings(proxyHost, proxyPort, nonProxyHosts)
			}(configWithProxy, config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
		}

		configWithSOCKSProxy, hasSOCKSProxySettings := v.Field(i).Interface().(WithSOCKSProxy)

		if hasSOCKSProxySettings {
			wg.Add(1)

			go func(configWithSOCKSProxy WithSOCKSProxy, socksProxyHost string, socksProxyPort string) {
				defer wg.Done()

				configWithSOCKSProxy.addSOCKSProxySettings(socksProxyHost, socksProxyPort)
			}(configWithSOCKSProxy, config.SOCKSProxyHost, config.SOCKSProxyPort)
		}
	}
	wg.Wait()
}

// DisableProxies disables proxies
func DisableProxies(config Configuration) {
	var wg sync.WaitGroup
	v := reflect.ValueOf(config.Targets)

	for i := 0; i < v.NumField(); i++ {
		configWithProxy, hasProxySettings := v.Field(i).Interface().(WithProxy)

		if hasProxySettings {
			wg.Add(1)

			go func(configWithProxy WithProxy) {
				defer wg.Done()

				configWithProxy.removeProxySettings()
			}(configWithProxy)
		}

		configWithSOCKSProxy, hasSOCKSProxySettings := v.Field(i).Interface().(WithSOCKSProxy)

		if hasSOCKSProxySettings {
			wg.Add(1)

			go func(configWithSOCKSProxy WithSOCKSProxy) {
				defer wg.Done()

				configWithSOCKSProxy.removeSOCKSProxySettings()
			}(configWithSOCKSProxy)
		}
	}
	wg.Wait()
}
