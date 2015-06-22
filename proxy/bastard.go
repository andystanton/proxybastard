package proxy

// Bastardise bastardises.
func Bastardise(config Configuration, enableProxies bool) {
	for _, shellFile := range config.Targets.Shell.Files {
		sanitisedPath := TildeToUserHome(shellFile)
		if enableProxies {
			writeSliceToFile(sanitisedPath, RemoveEnvVars(loadFileIntoSlice(sanitisedPath)))
			writeSliceToFile(sanitisedPath, AddEnvVars(loadFileIntoSlice(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		} else {
			writeSliceToFile(sanitisedPath, RemoveEnvVars(loadFileIntoSlice(sanitisedPath)))
		}
	}

	for _, mavenFile := range config.Targets.Maven.Files {
		sanitisedPath := TildeToUserHome(mavenFile)
		if enableProxies {
			writeXML(sanitisedPath, RemoveEnvVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
			writeXML(sanitisedPath, AddEnvVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		} else {
			writeXML(sanitisedPath, RemoveEnvVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		}
	}
}
