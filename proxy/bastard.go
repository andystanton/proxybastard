package proxy

// Bastardise bastardises.
func Bastardise(config Configuration, enableProxies bool) {
	if enableProxies {
		AddToShell(config)
		AddToMaven(config)
		AddToGit(config)
	} else {
		RemoveFromShell(config)
		RemoveFromMaven(config)
		RemoveFromGit(config)
	}
}
