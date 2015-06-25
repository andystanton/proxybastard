package proxy

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/clbanning/mxj"
)

// AddToMaven adds to maven.
func AddToMaven(config Configuration) {
	for _, mavenFile := range config.Targets.Maven.Files {
		sanitisedPath := TildeToUserHome(mavenFile)
		writeXML(sanitisedPath, AddEnvVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
	}
}

// RemoveFromMaven removes from Maven.
func RemoveFromMaven(config Configuration) {
	for _, mavenFile := range config.Targets.Maven.Files {
		sanitisedPath := TildeToUserHome(mavenFile)
		writeXML(sanitisedPath, RemoveEnvVarsMaven(loadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
	}
}

// AddEnvVarsMaven adds proxy vars to Maven.
func AddEnvVarsMaven(settingsXML mxj.Map, proxyHost string, proxyPort string, nonProxyHosts []string) mxj.Map {
	proxies, err := buildProxyVars(proxyHost, proxyPort, nonProxyHosts, true).ValuesForPath("proxies")
	if err != nil {
		log.Fatal("Unable to find proxies data in generated xml", err)
	}
	settingsXML.SetValueForPath(proxies[0], "settings.proxies")
	return settingsXML
}

// RemoveEnvVarsMaven adds proxy vars to Maven.
func RemoveEnvVarsMaven(settingsXML mxj.Map, proxyHost string, proxyPort string, nonProxyHosts []string) mxj.Map {
	proxies, err := buildProxyVars(proxyHost, proxyPort, nonProxyHosts, false).ValuesForPath("proxies")
	if err != nil {
		log.Fatal("Unable to find proxies data in generated xml", err)
	}
	settingsXML.SetValueForPath(proxies[0], "settings.proxies")
	return settingsXML
}

func buildProxyVars(proxyHost string, proxyPort string, nonProxyHosts []string, active bool) mxj.Map {
	shortHost := regexp.MustCompile("^http(s?)://").ReplaceAllString(proxyHost, "")
	nonProxyHostString := strings.Join(nonProxyHosts, ",")

	template := `
<proxies>
	<proxy>
		<protocol>http</protocol>
		<host>%s</host>
		<port>%s</port>
		<nonProxyHosts>%s</nonProxyHosts>
		<active>%t</active>
	</proxy>
	<proxy>
		<protocol>https</protocol>
		<host>%s</host>
		<port>%s</port>
		<nonProxyHosts>%s</nonProxyHosts>
		<active>%t</active>
	</proxy>
</proxies>`

	updated := fmt.Sprintf(template,
		shortHost, proxyPort, nonProxyHostString, active,
		shortHost, proxyPort, nonProxyHostString, active)

	xml, err := mxj.NewMapXml([]byte(updated))
	if err != nil {
		log.Fatal("Unable to generate required xml", err)
	}

	return xml
}
