package proxy

import "reflect"

// WithProxy is a thing that can have proxy settings added or removed.
type WithProxy interface {
	addProxySettings(string, string, []string)
	removeProxySettings()
}

// WithSOCKSProxy is a thing that can have SOCKS proxy settings added or removed.
type WithSOCKSProxy interface {
	addSocksProxySettings(string, string)
	removeSocksProxySettings()
}

// EnableProxies enable proxies.
func EnableProxies(config Configuration) {
	v := reflect.ValueOf(config.Targets)

	for i := 0; i < v.NumField(); i++ {
		configWithProxy, hasProxySettings := v.Field(i).Interface().(WithProxy)
		if hasProxySettings {
			configWithProxy.addProxySettings(config.ProxyHost, config.ProxyPort, config.NonProxyHosts)
		}

		configWithSOCKSProxy, hasSocksProxySettings := v.Field(i).Interface().(WithSOCKSProxy)
		if hasSocksProxySettings {
			configWithSOCKSProxy.addSocksProxySettings(config.SOCKSProxyHost, config.SOCKSProxyPort)
		}
	}
}

// DisableProxies disables proxies
func DisableProxies(config Configuration) {
	v := reflect.ValueOf(config.Targets)

	for i := 0; i < v.NumField(); i++ {
		configWithProxy, hasProxySettings := v.Field(i).Interface().(WithProxy)
		if hasProxySettings {
			configWithProxy.removeProxySettings()
		}

		configWithSOCKSProxy, hasSocksProxySettings := v.Field(i).Interface().(WithSOCKSProxy)
		if hasSocksProxySettings {
			configWithSOCKSProxy.removeSocksProxySettings()
		}
	}
}
