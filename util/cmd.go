package util

import "os/exec"

// ShellOut executes a command.
func ShellOut(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	return string(output), err
}
