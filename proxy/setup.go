package proxy

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/andystanton/proxybastard/util"
	"github.com/deckarep/golang-set"
)

func addToMap(frequencyMap map[string]int, value string) map[string]int {
	if len(strings.TrimSpace(value)) > 0 {
		if _, ok := frequencyMap[value]; !ok {
			frequencyMap[value] = 0
		}
		frequencyMap[value] = frequencyMap[value] + 1
	}
	return frequencyMap
}

func awaitInput(prompt string, pattern string) string {
	var matched string
	var found bool
	fmt.Printf("%s> ", prompt)
	for i := 0; i < 3; i++ {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		rPattern := regexp.MustCompile(pattern)
		result := strings.TrimSpace(text)
		if rPattern.MatchString(result) {
			matched = result
			found = true
			break
		} else if i < 2 {
			log.Println("That doesn't look right - try again...")
			log.Println()
			fmt.Printf("%s\n> ", prompt)
		} else {
			log.Println()
			log.Print("Three failed attempts - aborting!")
		}
	}
	if !found {
		os.Exit(1)
	}
	return matched
}

func awaitFileInput(prompt string) string {
	var matched string
	var found bool
	fmt.Printf("%s> ", prompt)
	for i := 0; i < 3; i++ {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		result := util.SanitisePath(strings.TrimSpace(text))
		if util.FileExists(result) {
			matched = result
			found = true
			break
		} else if i < 2 {
			log.Println("That doesn't look right - try again...")
			log.Println()
			fmt.Printf("%s\n> ", prompt)
		} else {
			log.Println()
			log.Print("Three failed attempts - aborting!")
		}
	}
	if !found {
		os.Exit(1)
	}
	return matched
}

// Setup presents the user with setup options.
func Setup(version string, acceptDefaults bool) {
	suggestedConfiguration := suggestConfiguration()
	actualConfiguration := Configuration{}
	actualConfiguration.Version = version
	readyToWrite := false

	if acceptDefaults {
		actualConfiguration = suggestedConfiguration
		readyToWrite = true
	} else {
		httpProxySet := false
		if len(suggestedConfiguration.ProxyHost) > 0 {
			message := fmt.Sprintf("Use suggested http proxy %s:%s? [Yn]\n", suggestedConfiguration.ProxyHost, suggestedConfiguration.ProxyPort)
			input := awaitInput(message, "(y|n|^$)")
			httpProxySet = strings.EqualFold(input, "y") || strings.EqualFold(input, "")
			actualConfiguration.ProxyHost = suggestedConfiguration.ProxyHost
			actualConfiguration.ProxyPort = suggestedConfiguration.ProxyPort
			fmt.Println()
		}

		if !httpProxySet {
			proxyHostPattern := "(?:https?://)?(.+):(\\d+)"
			proxyHostRegexp := regexp.MustCompile(proxyHostPattern)
			matches := proxyHostRegexp.FindStringSubmatch(awaitInput("Please enter an http proxy e.g. http://proxybastard:1234\n", proxyHostPattern))
			actualConfiguration.ProxyHost = fmt.Sprintf("http://%s", matches[1])
			actualConfiguration.ProxyPort = matches[2]
			httpProxySet = true
			fmt.Println()
		}

		socksProxySet := false
		if len(suggestedConfiguration.SOCKSProxyHost) > 0 {
			message := fmt.Sprintf("Use suggested SOCKS proxy %s:%s? [Yn]\n", suggestedConfiguration.SOCKSProxyHost, suggestedConfiguration.SOCKSProxyPort)
			input := awaitInput(message, "(y|n|^$)")
			socksProxySet = strings.EqualFold(input, "y") || strings.EqualFold(input, "")
			actualConfiguration.SOCKSProxyHost = suggestedConfiguration.SOCKSProxyHost
			actualConfiguration.SOCKSProxyPort = suggestedConfiguration.SOCKSProxyPort
			fmt.Println()
		}

		if !socksProxySet {
			socksHostPattern := "(?:(.+):(\\d+)|^$)"
			sockHostRegexp := regexp.MustCompile(socksHostPattern)
			matches := sockHostRegexp.FindStringSubmatch(awaitInput("Please enter a SOCKS proxy or press return for none e.g. socks.proxybastard:1234\n", socksHostPattern))
			if len(matches) > 0 {
				socksProxySet = true
				actualConfiguration.SOCKSProxyHost = matches[1]
				actualConfiguration.SOCKSProxyPort = matches[2]
			}
			fmt.Println()
		}

		if suggestedConfiguration.Targets != nil {
			targetsField := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
			for i := 0; i < targetsField.NumField(); i++ {
				fieldName := targetsField.Type().Field(i).Name

				// valueForFieldRequired := false
				if !util.InterfaceIsZero(targetsField.Field(i).Interface()) {
					targetField := reflect.Indirect(reflect.ValueOf(targetsField.Field(i).Interface()))

					// step 1 - go simple and ask if user wants suggested config
					message := fmt.Sprintf("Found %s! Use suggested configuration? [Yn]\n", fieldName)
					input := awaitInput(message, "(y|n|^$)")
					configurationSet := strings.EqualFold(input, "y") || strings.EqualFold(input, "")

					if configurationSet {
						// great! put the suggested config in the actual config
						if actualConfiguration.Targets == nil {
							actualConfiguration.Targets = &TargetsConfiguration{}
						}
						actualField := reflect.Indirect(reflect.ValueOf(actualConfiguration.Targets)).FieldByName(fieldName)
						actualField.Set(reflect.ValueOf(targetsField.Field(i).Interface()))
					} else {
						// hmmm, the user wants to do something else. well that's ok - there are three types of
						// config. From the simplest to the most complex:
						// 1. config only contains the "enabled" element - allow user to enable or disable
						// 2. config contains "enabled" and "files" - allow user to enable or disable and manually enter files.
						// 3. config is more complex and has non-standard elements in it. examples include shell which allows
						//	  java opts to be set, and stunnel which includes a flag to kill the stunnel process.
						fmt.Println()

						message := fmt.Sprintf("\tEnable %s? [Yn]\n\t", fieldName)
						input := awaitInput(message, "(y|n|^$)")
						if strings.EqualFold(input, "y") || strings.EqualFold(input, "") {
							if actualConfiguration.Targets == nil {
								actualConfiguration.Targets = &TargetsConfiguration{}
							}
							actualField := reflect.New(reflect.TypeOf(targetField.Interface())).Interface()
							reflect.Indirect(reflect.ValueOf(actualField)).FieldByName("Enabled").Set(reflect.ValueOf(true))
							fmt.Println()

							if util.ValueHasField(targetField, "Files") {
								fieldFiles := targetField.FieldByName("Files").Interface().([]string)
								message := fmt.Sprintf("\tUse suggested file %s? [Yn]\n\t", strings.Join(fieldFiles, ", "))
								input := awaitInput(message, "(y|n|^$)")
								if strings.EqualFold(input, "y") || strings.EqualFold(input, "") {

								} else {
									message := fmt.Sprintf("\tEnter the path to a valid %s configuration file\n\t", fieldName)
									input := awaitFileInput(message)
									reflect.Indirect(reflect.ValueOf(actualField)).FieldByName("Files").Set(reflect.ValueOf([]string{input}))
								}
								fmt.Println()
							}

							if util.ValueHasMethod(targetField, "CustomPrompt") {
								customMethod := targetField.MethodByName("CustomPrompt")
								actualField = customMethod.Call([]reflect.Value{})[0].Interface()
							}

							reflect.Indirect(reflect.ValueOf(actualConfiguration.Targets)).FieldByName(fieldName).Set(reflect.ValueOf(actualField))
						}
					}
					fmt.Println()
				} else {
					fmt.Println("Suggested config does not contain " + fieldName)
				}
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 12, 0, 1, ' ', 0)

		fmt.Fprintln(w, "Settings")
		fmt.Fprintln(w, "================================================================")
		fmt.Fprintln(w, "HTTP Proxy")
		fmt.Fprintf(w, " - Host\t : %s\n", actualConfiguration.ProxyHost)
		fmt.Fprintf(w, " - Port\t : %s\n", actualConfiguration.ProxyPort)

		if len(actualConfiguration.SOCKSProxyHost) > 0 {
			fmt.Fprintln(w)
			fmt.Fprintln(w, "SOCKS Proxy")
			fmt.Fprintf(w, " - Host\t : %s\n", actualConfiguration.SOCKSProxyHost)
			fmt.Fprintf(w, " - Port\t : %s\n", actualConfiguration.SOCKSProxyPort)
		}
		fmt.Fprintln(w, "================================================================")
		fmt.Fprintln(w)

		if actualConfiguration.Targets != nil {
			targetsField := reflect.Indirect(reflect.ValueOf(actualConfiguration.Targets))
			for i := 0; i < targetsField.NumField(); i++ {
				fieldName := targetsField.Type().Field(i).Name

				if !util.InterfaceIsZero(targetsField.Field(i).Interface()) {
					targetField := reflect.Indirect(reflect.ValueOf(targetsField.Field(i).Interface()))
					withConfig, _ := targetField.Interface().(WithConfig)
					fmt.Fprintln(w, fieldName)
					fmt.Fprintf(w, " - Enabled\t : %v\n", withConfig.isEnabled())

					if util.ValueHasField(targetField, "Files") {
						fieldFiles := targetField.FieldByName("Files").Interface().([]string)
						fmt.Fprintf(w, " - Files\t : %s\n", strings.Join(fieldFiles, ","))
					}

					if util.ValueHasMethod(targetField, "CustomFields") {
						customMethod := targetField.MethodByName("CustomFields")
						extraFields := customMethod.Call([]reflect.Value{})[0].Interface().(map[string]string)
						for k, v := range extraFields {
							fmt.Fprintf(w, " - %s\t : %s\n", k, v)
						}
					}
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintln(w, "================================================================")
		fmt.Fprintln(w)
		w.Flush()

		input := awaitInput("Write these settings to ~/.proxybastard/config.json? [Yn]\n", "(y|n|^$)")
		fmt.Println()
		readyToWrite = strings.EqualFold(input, "y") || strings.EqualFold(input, "")
	}

	// marshalled, err := json.MarshalIndent(actualConfiguration, "", "    ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n\n", string(marshalled))
	if readyToWrite {
		fmt.Println("Done")
	} else {
		fmt.Println("kthx")
	}
}

func getHighestFrequency(frequencyMap map[string]int) string {
	var highestValue int
	var mostFrequentKey string

	keys := []string{}
	for key := range frequencyMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if value := frequencyMap[key]; value > highestValue {
			highestValue = value
			mostFrequentKey = key
		}
	}
	return mostFrequentKey
}

func suggestConfiguration() Configuration {
	var suggestedConfiguration Configuration

	suggestedProxyHosts := make(map[string]int)
	suggestedProxyPorts := make(map[string]int)
	suggestedSOCKSProxyHosts := make(map[string]int)
	suggestedSOCKSProxyPorts := make(map[string]int)
	suggestedNonProxyHosts := mapset.NewSet()

	targetsField := reflect.Indirect(reflect.ValueOf(&TargetsConfiguration{}))
	for i := 0; i < targetsField.NumField(); i++ {

		configurationField := reflect.New(reflect.TypeOf(targetsField.Field(i).Interface()).Elem()).Interface()
		withConfig, hasConfig := configurationField.(WithConfig)
		if hasConfig {
			fieldName := targetsField.Type().Field(i).Name

			if suggestedItemConfiguration := withConfig.suggestConfiguration(); suggestedItemConfiguration != nil {
				if suggestedConfiguration.Targets == nil {
					suggestedConfiguration.Targets = &TargetsConfiguration{}
				}

				addToMap(suggestedProxyHosts, util.SanitiseHTTPProxyURL(suggestedItemConfiguration.ProxyHost))
				addToMap(suggestedProxyPorts, suggestedItemConfiguration.ProxyPort)
				addToMap(suggestedSOCKSProxyHosts, suggestedItemConfiguration.SOCKSProxyHost)
				addToMap(suggestedSOCKSProxyPorts, suggestedItemConfiguration.SOCKSProxyPort)

				for _, nonProxyHost := range suggestedItemConfiguration.NonProxyHosts {
					suggestedNonProxyHosts.Add(nonProxyHost)
				}

				targetsField := reflect.Indirect(reflect.ValueOf(suggestedConfiguration.Targets))
				targetsFieldSuggested := reflect.Indirect(reflect.ValueOf(suggestedItemConfiguration.Targets))

				targetsField.FieldByName(fieldName).Set(targetsFieldSuggested.FieldByName(fieldName))
			}
		}
	}

	suggestedConfiguration.ProxyHost = getHighestFrequency(suggestedProxyHosts)
	suggestedConfiguration.ProxyPort = getHighestFrequency(suggestedProxyPorts)
	suggestedConfiguration.SOCKSProxyHost = getHighestFrequency(suggestedSOCKSProxyHosts)
	suggestedConfiguration.SOCKSProxyPort = getHighestFrequency(suggestedSOCKSProxyPorts)
	for suggestedNonProxyHost := range suggestedNonProxyHosts.Iter() {
		suggestedConfiguration.NonProxyHosts = append(suggestedConfiguration.NonProxyHosts, suggestedNonProxyHost.(string))
	}

	return suggestedConfiguration
}
