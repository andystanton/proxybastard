package proxy

func Bastardise(config Configuration, enableProxies bool) {
	for _, shellFile := range config.Targets.Shell {
		sanitisedPath := TildeToUserHome(shellFile)
		if enableProxies {
			writeSliceToFile(sanitisedPath, RemoveProxyVars(loadFileIntoSlice(sanitisedPath)))
			writeSliceToFile(sanitisedPath, AddProxyVars(loadFileIntoSlice(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		} else {
			writeSliceToFile(sanitisedPath, RemoveProxyVars(loadFileIntoSlice(sanitisedPath)))
		}
	}

	for _, mavenFile := range config.Targets.Maven {
		sanitisedPath := TildeToUserHome(mavenFile)
		if enableProxies {
			writeXML(sanitisedPath, RemoveProxyVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
			writeXML(sanitisedPath, AddProxyVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		} else {
			writeXML(sanitisedPath, RemoveProxyVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		}
	}
}
