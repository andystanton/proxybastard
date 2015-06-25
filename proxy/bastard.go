package proxy

// Bastardise bastardises.
func Bastardise(config Configuration, enableProxies bool) {
	if enableProxies {
		AddToShell(config)
		AddToMaven(config)
	} else {
		RemoveFromShell(config)
		RemoveFromMaven(config)
	}
}
