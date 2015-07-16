package proxy

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
	"github.com/clbanning/mxj"
)

func (mavenConfiguration MavenConfiguration) validate() error {
	return nil
}

func (mavenConfiguration MavenConfiguration) isEnabled() bool {
	return mavenConfiguration.Enabled
}

func (mavenConfiguration MavenConfiguration) suggestConfiguration() interface{} {
	return nil
}

func (mavenConfiguration MavenConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	for _, mavenFile := range mavenConfiguration.Files {
		sanitisedPath := util.SanitisePath(mavenFile)
		util.WriteXML(sanitisedPath, addToMavenXML(util.LoadXML(sanitisedPath), proxyHost, proxyPort, nonProxyHosts))
	}
}

func (mavenConfiguration MavenConfiguration) removeProxySettings() {
	for _, mavenFile := range mavenConfiguration.Files {
		sanitisedPath := util.SanitisePath(mavenFile)
		util.WriteXML(sanitisedPath, removeFromMavenXML(util.LoadXML(sanitisedPath)))
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

func removeFromMavenXML(settingsXML mxj.Map) mxj.Map {
	settingsXML.Remove("settings.proxies")
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
