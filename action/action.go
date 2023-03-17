package action

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Action(serviceName string, configPath string, gitTag string) error {
	var err = os.Chdir(configPath)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v\n", serviceName, err)
	}

	if err := runCmd("git", "fetch", "--all"); err != nil {
		return err
	}

	if err := runCmd("git", "reset", "--hard", "origin/"+gitTag); err != nil {
		return err
	}

	if err := runCmd("git", "rebase", "origin/"+gitTag); err != nil {
		return err
	}

	if err := runCmd("systemctl", "kill", "--signal=HUP", serviceName); err != nil {
		return err
	}

	return nil
}

func runCmd(commandline string, args ...string) error {
	var cmd = exec.Command(commandline, args...)
	fmt.Println(cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	var err = cmd.Run()

	if err != nil {
		if strings.Contains(err.Error(), "signal: hangup") {
			// ignore
		} else {
			return fmt.Errorf("failed to run  '%s', error: %v\n", commandline, err)
		}
	}

	return nil
}
