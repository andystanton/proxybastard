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

// BackupMode for the backup/restore functionality.
type BackupMode int

const (
	// Backup indicates a backup should be performed.
	Backup BackupMode = iota
	// Restore indicates a restore should be performed.
	Restore
)

// DirtyBackupOperation performs a dirty backup or restore.
func DirtyBackupOperation(config Configuration, mode BackupMode) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	userHome := usr.HomeDir
	userHomeRegex := regexp.MustCompile(strings.Replace(fmt.Sprintf("^%s/.+$", userHome), "/", "\\/", -1))
	backupPath := fmt.Sprintf("%s/%s", userHome, ".proxybastard/backup")

	reflected := reflect.Indirect(reflect.ValueOf(config.Targets))
	if reflect.ValueOf(config.Targets) != reflect.Zero(reflect.TypeOf(config.Targets)) {
		for i := 0; i < reflected.NumField(); i++ {

			fieldInterface := reflected.Field(i).Interface()

			if reflect.ValueOf(fieldInterface) != reflect.Zero(reflect.TypeOf(fieldInterface)) {
				reflectedField := reflect.Indirect(reflect.ValueOf(fieldInterface))

				if reflectedField.FieldByName("Files").Kind() != reflect.Invalid {
					files := reflectedField.FieldByName("Files").Interface().([]string)
					for _, file := range files {

						sanitisedFile := util.SanitisePath(file)
						if userHomeRegex.MatchString(sanitisedFile) {
							fileBits := regexp.MustCompile(fmt.Sprintf("%s/(.+/)?(.+)", userHome)).FindStringSubmatch(sanitisedFile)

							pathToFile := fileBits[1]
							fileName := fileBits[2]

							fileBackupPath := strings.TrimSuffix(fmt.Sprintf("%s/%s", backupPath, pathToFile), "/")

							if mode == Backup {
								util.ShellOut("mkdir", []string{"-p", fileBackupPath})

								fmt.Printf("Backing up %s to %s\n", sanitisedFile, fmt.Sprintf("%s/%s", fileBackupPath, fileName))
								util.ShellOut("cp", []string{"-rf", sanitisedFile, fmt.Sprintf("%s/%s", fileBackupPath, fileName)})
							} else {
								fmt.Printf("Restoring %s from %s\n", sanitisedFile, fmt.Sprintf("%s/%s", fileBackupPath, fileName))
								util.ShellOut("cp", []string{"-rf", fmt.Sprintf("%s/%s", fileBackupPath, fileName), sanitisedFile})
							}

						} else {
							fmt.Printf("Unable to backup %s - not in user home\n", sanitisedFile)
						}
					}
				}
			}
		}
	}
}
