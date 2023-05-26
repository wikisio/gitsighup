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

func Add(CfgMsg string, namespace string, service string, filename string, dst string) error {
	var err = os.Chdir(dst)
	if err != nil {
		return fmt.Errorf("failed to change the dir %s, error: %v", dst, err)
	}

	data, err := base64.StdEncoding.DecodeString(CfgMsg)
	if err != nil {
		return fmt.Errorf("decode err : %v", err)
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file err : %v", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("write data err : %v", err)
	}

	if err := runCmd("systemctl", "kill", "--signal=HUP", service); err != nil {
		return fmt.Errorf("run cmd err : %v", err)
	}
	return Restart(service, dst)
}

func Edit(CfgMsg string, namespace string, service string, filename string, dst string) error {
	var err = os.Chdir(dst)
	if err != nil {
		return fmt.Errorf("failed to change the dir of %s, error: %v", service, err)
	}
	data, err := base64.StdEncoding.DecodeString(CfgMsg)
	if err != nil {
		return fmt.Errorf("decode err : %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open data err : %v", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("write data err : %v", err)
	}
	if err := runCmd("systemctl", "kill", "--signal=HUP", service); err != nil {
		return fmt.Errorf("run cmd err : %v", err)
	}
	if service == "gitsighup" {
		if err = config2.LoadConfig(); err != nil {
			os.Exit(1)
		}
	}

	return Restart(service, dst)
}
