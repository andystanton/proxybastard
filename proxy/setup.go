package proxy

import (
	"bufio"
	"encoding/json"
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

func awaitInput(prompt string, pattern string, prefix string) string {
	var matched string
	var found bool
	fmt.Printf("%s%s\n", prefix, prompt)
	fmt.Printf("%s> ", prefix)
	for i := 0; i < 3; i++ {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		rPattern := regexp.MustCompile(pattern)

		if result := strings.TrimSpace(text); rPattern.MatchString(result) {
			matched = result
			found = true
			break
		} else if i < 2 {
			log.Printf("%sThat doesn't look right - try again...\n", prefix)
			log.Println()
			fmt.Printf("%s%s\n", prefix, prompt)
			fmt.Printf("%s> ", prefix)
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

func awaitFileInput(prompt string, prefix string) string {
	var matched string
	var found bool
	fmt.Printf("%s%s\n", prefix, prompt)
	fmt.Printf("%s> ", prefix)
	for i := 0; i < 3; i++ {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if file := util.SanitisePath(strings.TrimSpace(text)); util.FileExists(file) {
			matched = file
			found = true
			break
		} else if i < 2 {
			log.Printf("%sThat doesn't look right - try again...\n", prefix)
			log.Println()
			fmt.Printf("%s%s\n", prefix, prompt)
			fmt.Printf("%s> ", prefix)
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
			message := fmt.Sprintf("Use suggested http proxy %s:%s? [Yn]", suggestedConfiguration.ProxyHost, suggestedConfiguration.ProxyPort)
			input := awaitInput(message, "(y|n|^$)", "")
			httpProxySet = strings.EqualFold(input, "y") || strings.EqualFold(input, "")
			actualConfiguration.ProxyHost = suggestedConfiguration.ProxyHost
			actualConfiguration.ProxyPort = suggestedConfiguration.ProxyPort
			fmt.Println()
		}

		if !httpProxySet {
			proxyHostPattern := "(?:https?://)?(.+):(\\d+)"
			proxyHostRegexp := regexp.MustCompile(proxyHostPattern)
			matches := proxyHostRegexp.FindStringSubmatch(awaitInput("Please enter an http proxy e.g. http://proxybastard:1234", proxyHostPattern, ""))
			actualConfiguration.ProxyHost = fmt.Sprintf("http://%s", matches[1])
			actualConfiguration.ProxyPort = matches[2]
			httpProxySet = true
			fmt.Println()
		}

		socksProxySet := false
		if len(suggestedConfiguration.SOCKSProxyHost) > 0 {
			message := fmt.Sprintf("Use suggested SOCKS proxy %s:%s? [Yn]", suggestedConfiguration.SOCKSProxyHost, suggestedConfiguration.SOCKSProxyPort)
			input := awaitInput(message, "(y|n|^$)", "")
			socksProxySet = strings.EqualFold(input, "y") || strings.EqualFold(input, "")
			actualConfiguration.SOCKSProxyHost = suggestedConfiguration.SOCKSProxyHost
			actualConfiguration.SOCKSProxyPort = suggestedConfiguration.SOCKSProxyPort
			fmt.Println()
		}

		if !socksProxySet {
			socksHostPattern := "(?:(.+):(\\d+)|^$)"
			sockHostRegexp := regexp.MustCompile(socksHostPattern)
			matches := sockHostRegexp.FindStringSubmatch(awaitInput("Please enter a SOCKS proxy or press return for none e.g. socks.proxybastard:1234", socksHostPattern, ""))
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
					message := fmt.Sprintf("Found %s! Use suggested configuration? [Yn]", fieldName)
					input := awaitInput(message, "(y|n|^$)", "")
					configurationSet := strings.EqualFold(input, "y") || strings.EqualFold(input, "")

					if configurationSet {
						// great! put the suggested config in the actual config
						if actualConfiguration.Targets == nil {
							actualConfiguration.Targets = &TargetsConfiguration{}
						}
						actualField := reflect.Indirect(reflect.ValueOf(actualConfiguration.Targets)).FieldByName(fieldName)
						actualField.Set(reflect.ValueOf(targetsField.Field(i).Interface()))
					} else {
						fmt.Println()

						message := fmt.Sprintf("Enable %s? [Yn]", fieldName)
						input := awaitInput(message, "(y|n|^$)", "  ")
						if strings.EqualFold(input, "y") || strings.EqualFold(input, "") {
							if actualConfiguration.Targets == nil {
								actualConfiguration.Targets = &TargetsConfiguration{}
							}
							actualField := reflect.New(reflect.TypeOf(targetField.Interface())).Interface()
							reflect.Indirect(reflect.ValueOf(actualField)).FieldByName("Enabled").Set(reflect.ValueOf(true))
							fmt.Println()

							if util.ValueHasField(targetField, "Files") {
								fieldFiles := targetField.FieldByName("Files").Interface().([]string)
								message := fmt.Sprintf("Use suggested file %s? [Yn]", strings.Join(fieldFiles, ", "))
								input := awaitInput(message, "(y|n|^$)", "  ")
								if strings.EqualFold(input, "y") || strings.EqualFold(input, "") {
									reflect.Indirect(reflect.ValueOf(actualField)).FieldByName("Files").Set(reflect.ValueOf(fieldFiles))
								} else {
									fmt.Println()
									message := fmt.Sprintf("Enter the path to a valid %s configuration file", fieldName)
									input := awaitFileInput(message, "  ")
									reflect.Indirect(reflect.ValueOf(actualField)).FieldByName("Files").Set(reflect.ValueOf([]string{input}))
								}
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

		w := tabwriter.NewWriter(os.Stdout, 20, 0, 1, ' ', 0)

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
						newFiles := []string{}
						for _, file := range fieldFiles {
							newFiles = append(newFiles, util.UnsanitisePath(file))
						}
						fmt.Fprintf(w, " - Files\t : %s\n", strings.Join(newFiles, ","))
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

		input := awaitInput("Write these settings to ~/.proxybastard/config.json? [Yn]", "(y|n|^$)", "")
		fmt.Println()
		readyToWrite = strings.EqualFold(input, "y") || strings.EqualFold(input, "")
	}

	marshalled, err := json.MarshalIndent(actualConfiguration, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n\n", string(marshalled))
	if readyToWrite {
		util.WriteSliceToFile(util.SanitisePath("~/.proxybastard/config.json"), strings.Split(string(marshalled), "\n"))
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
