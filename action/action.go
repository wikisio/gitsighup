package action

import (
	"fmt"
	"os"
	"os/exec"
)

func Action(serviceName string, configPath string, gitTag string) error {
	var err = os.Chdir(configPath)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v", serviceName, err)
	}

	if err := runCmd("git", "pull", "-f", "origin", gitTag); err != nil {
		return err
	}

	if err := runCmd("/usr/bin/systemctl", "kill", "--signal=HUP", serviceName); err != nil {
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
		return fmt.Errorf("failed to run  '%s', error: %v", commandline, err)
	}

	return nil
}
