package proxy

// EnableProxies enable proxies.
func EnableProxies(config Configuration) {
	AddToShell(config)
	AddToMaven(config)
	AddToGit(config)
}

// DisableProxies disables proxies
func DisableProxies(config Configuration) {
	RemoveFromShell(config)
	RemoveFromMaven(config)
	RemoveFromGit(config)
}
