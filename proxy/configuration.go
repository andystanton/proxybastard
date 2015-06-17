package proxy

// Configuration represents the Proxybastard configuration.
type Configuration struct {
	ProxyHost     string
	ProxyPort     *int
	NonProxyHosts []string
	ShellFiles    []string
	MavenFiles    []string
}
