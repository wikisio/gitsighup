package action

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	config2 "wikis.io/config"
)

func Action(serviceName string, configPath string, gitTag string) error {
	var err = os.Chdir(configPath)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v", serviceName, err)
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

func Restart(serviceName string, configPath string) error {

	if err := runCmd("systemctl", "restart", serviceName); err != nil {
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
			return fmt.Errorf("failed to run  '%s', error: %v", commandline, err)
		}
	}

	return nil
}

func ADD(CfgMsg string, namespace string, service string, filename string, dst string) error {
	var err = os.Chdir(dst)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v", service, err)
	}

	data, err := base64.StdEncoding.DecodeString(CfgMsg)
	if err != nil {
		return err
	}
	file, err := os.Create(filename + ".yml")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)

	if err := runCmd("systemctl", "kill", "--signal=HUP", service); err != nil {
		return err
	}
	return Restart(service, dst)
}

func EDIT(CfgMsg string, namespace string, service string, filename string, dst string) error {
	var err = os.Chdir(dst)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v", service, err)
	}
	data, err := base64.StdEncoding.DecodeString(CfgMsg)
	if err != nil {
		return err
	}

	file, err := os.Open(filename + ".yml")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)

	if err := runCmd("systemctl", "kill", "--signal=HUP", service); err != nil {
		return err
	}
	if service == "gitsighup" {
		err = config2.LoadConfig()
		if err != nil {
			os.Exit(1)
		}
	}

	return Restart(service, dst)
}
