package proxy

// Configuration represents the Proxybastard configuration.
type Configuration struct {
	ProxyHost  string
	ProxyPort  *int
	ShellFiles []string
	MavenFiles []string
}
