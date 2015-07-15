package proxy

import (
	"fmt"
	"log"
	"os/user"
	"reflect"
	"regexp"
	"strings"

	"github.com/andystanton/proxybastard/util"
)

// DirtyBackup does a backup by using reflection over the struct. Filthy.
func DirtyBackup(config Configuration) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	userHome := usr.HomeDir
	userHomeRegex := regexp.MustCompile(strings.Replace(fmt.Sprintf("^%s/.+$", userHome), "/", "\\/", -1))
	backupPath := fmt.Sprintf("%s/%s", userHome, ".proxybastard/backup")

	util.ShellOut("mkdir", []string{"-p", backupPath})

	reflected := reflect.ValueOf(config.Targets)
	for i := 0; i < reflected.NumField(); i++ {

		reflectedField := reflect.ValueOf(reflected.Field(i).Interface())
		for j := 0; j < reflectedField.NumField(); j++ {

			fieldName := reflectedField.Type().Field(j).Name
			if fieldName == "Files" {

				fmt.Println(reflected.Type().Field(i).Name)
				files := reflectedField.Field(j).Interface().([]string)
				for _, file := range files {

					sanitisedFile := util.SanitisePath(file)
					if userHomeRegex.MatchString(sanitisedFile) {
						fileBits := regexp.MustCompile(fmt.Sprintf("%s/(.+/)?(.+)", userHome)).FindStringSubmatch(sanitisedFile)
						pathToFile := strings.TrimSuffix(fileBits[1], "/")
						fileName := fileBits[2]
						fileBackupPath := fmt.Sprintf("%s/%s", backupPath, pathToFile)
						fmt.Println(fileBackupPath)
						util.ShellOut("mkdir", []string{"-p", fileBackupPath})
						fmt.Printf("copying %s to %s\n", sanitisedFile, fmt.Sprintf("%s/%s", fileBackupPath, fileName))
						util.ShellOut("cp", []string{"-rf", sanitisedFile, fmt.Sprintf("%s/%s", fileBackupPath, fileName)})
					} else {
						fmt.Println("Unable to backup paths not in user home!")
					}
				}
			}
		}
	}
}
