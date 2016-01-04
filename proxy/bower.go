package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/andystanton/proxybastard/util"
)

func (bowerConfiguration BowerConfiguration) validate() error {
	return nil
}

func (bowerConfiguration BowerConfiguration) isEnabled() bool {
	return bowerConfiguration.Enabled
}

func loadBowerRc(filename string) (map[string]interface{}, error) {
	configBytes, err := ioutil.ReadFile(util.SanitisePath(filename))
	if err != nil {
		return make(map[string]interface{}), err
	}
	return parseBowerRCJSON(configBytes), nil
}

func parseBowerRCJSON(jsoncontent []byte) map[string]interface{} {
	bowerRC := make(map[string]interface{})
	err := json.Unmarshal(jsoncontent, &bowerRC)
	if err != nil {
		log.Fatal(err)
	}
	if val, ok := bowerRC["proxy"]; ok {
		bowerRC["proxy"] = util.SanitiseHTTPProxyURL(val.(string))
	}
	if val, ok := bowerRC["https-proxy"]; ok {
		bowerRC["https-proxy"] = util.SanitiseHTTPProxyURL(val.(string))
	}
	return bowerRC
}

func extractProxyFromBowerContents(bowerRC map[string]interface{}) (string, string) {
	found := false
	var fullProxy string
	suggestedProxy, suggestedPort := "", ""
	if val, ok := bowerRC["proxy"]; ok {
		found = true
		fullProxy = val.(string)
	} else if val, ok := bowerRC["https-proxy"]; ok {
		found = true
		fullProxy = val.(string)
	}
	if found {
		hostRegexp := regexp.MustCompile("(.+):(.+)")
		hostMatches := hostRegexp.FindStringSubmatch(fullProxy)
		if len(hostMatches) > 0 {
			suggestedProxy = hostMatches[1]
			suggestedPort = hostMatches[2]
		} else {
			suggestedProxy = fullProxy
		}
	}
	return suggestedProxy, suggestedPort
}

func (bowerConfiguration BowerConfiguration) suggestConfiguration() (configuration *Configuration) {
	bowerExecutable := "bower"
	bowerFile := "~/.bowerrc"
	bowerFileSanitised := util.SanitisePath(bowerFile)

	_, err := util.ShellOut("which", []string{bowerExecutable})
	hasBower := err == nil
	hasBowerRC := util.FileExists(bowerFileSanitised)

	var suggestedProxy, suggestedPort string

	if hasBower {

		if hasBowerRC {
			bowerRC, _ := loadBowerRc(bowerFileSanitised)
			suggestedProxy, suggestedPort = extractProxyFromBowerContents(bowerRC)
		}

		return &Configuration{
			ProxyHost: suggestedProxy,
			ProxyPort: suggestedPort,
			Targets: &TargetsConfiguration{
				Bower: &BowerConfiguration{
					Enabled: true,
					Files:   []string{bowerFile},
				},
			},
		}
	}
	return nil
}

func addBowerProxySettings(contents map[string]interface{}, proxyHost string, proxyPort string) map[string]interface{} {
	out := contents
	out["proxy"] = fmt.Sprintf("%s:%s", proxyHost, proxyPort)
	out["https-proxy"] = fmt.Sprintf("%s:%s", proxyHost, proxyPort)
	return out
}

func removeBowerProxySettings(contents map[string]interface{}) map[string]interface{} {
	out := contents
	delete(out, "proxy")
	delete(out, "https-proxy")
	return out
}

func (bowerConfiguration BowerConfiguration) addProxySettings(proxyHost string, proxyPort string, nonProxyHosts []string) {
	bowerConfiguration.removeProxySettings()
	for _, file := range bowerConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := loadBowerRc(sanitisedPath)
		bowerRC, _ := json.MarshalIndent(addBowerProxySettings(contents, proxyHost, proxyPort), "", "    ")
		util.WriteToFile(sanitisedPath, bowerRC)
	}
}

func (bowerConfiguration BowerConfiguration) removeProxySettings() {
	for _, file := range bowerConfiguration.Files {
		sanitisedPath := util.SanitisePath(file)
		contents, _ := loadBowerRc(sanitisedPath)
		bowerRC, _ := json.Marshal(removeBowerProxySettings(contents))
		util.WriteToFile(sanitisedPath, bowerRC)
	}
}
