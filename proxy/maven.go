package proxy

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
	"github.com/clbanning/mxj"
)

func addToMaven(config Configuration) {
	if config.Targets.Maven.Enabled {
		for _, mavenFile := range config.Targets.Maven.Files {
			sanitisedPath := util.SanitisePath(mavenFile)
			util.WriteXML(sanitisedPath, addToMavenXML(util.LoadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		}
	}
}

func removeFromMaven(config Configuration) {
	if config.Targets.Maven.Enabled {
		for _, mavenFile := range config.Targets.Maven.Files {
			sanitisedPath := util.SanitisePath(mavenFile)
			util.WriteXML(sanitisedPath, removeFromMavenXML(util.LoadXML(sanitisedPath), config.ProxyHost, config.ProxyPort, config.NonProxyHosts))
		}
	}
}

func addToMavenXML(settingsXML mxj.Map, proxyHost string, proxyPort string, nonProxyHosts []string) mxj.Map {
	proxies, err := buildProxyVars(proxyHost, proxyPort, nonProxyHosts, true).ValuesForPath("proxies")
	if err != nil {
		log.Fatal("Unable to find proxies data in generated xml", err)
	}
	settingsXML.SetValueForPath(proxies[0], "settings.proxies")
	return settingsXML
}

func removeFromMavenXML(settingsXML mxj.Map, proxyHost string, proxyPort string, nonProxyHosts []string) mxj.Map {
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
		<id>http://%s:%s</id>
		<protocol>http</protocol>
		<host>%s</host>
		<port>%s</port>
		<nonProxyHosts>%s</nonProxyHosts>
		<active>%t</active>
	</proxy>
	<proxy>
		<id>https://%s:%s</id>
		<protocol>https</protocol>
		<host>%s</host>
		<port>%s</port>
		<nonProxyHosts>%s</nonProxyHosts>
		<active>%t</active>
	</proxy>
</proxies>`

	updated := fmt.Sprintf(template,
		shortHost, proxyPort, shortHost, proxyPort, nonProxyHostString, active,
		shortHost, proxyPort, shortHost, proxyPort, nonProxyHostString, active)

	xml, err := mxj.NewMapXml([]byte(updated))
	if err != nil {
		log.Fatal("Unable to generate required xml", err)
	}

	return xml
}
