package action

import (
	"fmt"
	"os"
	"os/exec"
)

func Action(servcieName string, configPath string, gitTag string) error {
	var err = os.Chdir(configPath)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v", servcieName, err)
	}

	var gitCmd = fmt.Sprintf("git pull -f origin %s", gitTag)
	if err := runCmd(gitCmd); err != nil {
		return err
	}

	var systemctlCmd = fmt.Sprintf("systemctl kill --signal=HUP %s", servcieName)
	if err := runCmd(systemctlCmd); err != nil {
		return err
	}

	return nil
}

func runCmd(commandline string) error {
	var cmd = exec.Command(commandline)
	var err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run  '%s', error: %v", commandline, err)
	}

	return nil
}
