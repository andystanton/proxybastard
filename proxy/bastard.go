package proxy

func Bastardise(config Configuration, enableProxies bool) {
	for _, shellFile := range config.ShellFiles {
		sanitisedPath := tildeToUserHome(shellFile)
		if enableProxies {
			writeSliceToFile(sanitisedPath, AddProxyVars(loadFileIntoSlice(sanitisedPath), config.ProxyHost, config.ProxyPort))
		} else {
			writeSliceToFile(sanitisedPath, RemoveProxyVars(loadFileIntoSlice(sanitisedPath)))
		}
	}
}
