package proxy

// EnableProxies enable proxies.
func EnableProxies(config Configuration) {
	addToShell(config)
	addToMaven(config)
	addToGit(config)
	addToNPM(config)
	addToSSH(config)
	addToAPM(config)
	addToSubversion(config)
}

// DisableProxies disables proxies
func DisableProxies(config Configuration) {
	removeFromShell(config)
	removeFromMaven(config)
	removeFromGit(config)
	removeFromNPM(config)
	removeFromSSH(config)
	removeFromAPM(config)
	removeFromSubversion(config)
}
