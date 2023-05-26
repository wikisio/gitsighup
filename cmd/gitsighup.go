package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wikis.io/action"
	config2 "wikis.io/config"
)

var jwt string = "jwt"

func main() {
	var configFile = flag.String("c", "C:/Users/82742/fork/configsrv.yml", "the yaml config file")
	flag.Parse()
	config2.GlobalConfigFile = *configFile

	var err = config2.LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	var c = make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP)
	go config2.Refresh(c)
	jwt, err = Login()
	if err != nil {
		os.Exit(401)
	}

	var namespace string
	var service string
	var filename string
	var dst string
	for {
		for _, i := range config2.GlobalConfig.Services {
			namespace = i.NameSpace
			service = i.Name
			for _, j := range i.ConfigPath {
				filename = j.Src
				dst = j.Dst
				url := config2.ApiUrl
				fmt.Print(SendRequest(url, namespace, service, filename, dst, jwt))
			}

		}
	}
}

func SendRequest(url string, namespace string, service string, filename string, dst string, jwt string) error {
	url = url + namespace + "/" + service + "/" + filename
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new request err: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request err: %v", err)
	}
	defer resp.Body.Close()
	//登录
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("io read err: %v", err)
		}
		var result CfgResult
		if err = json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("unmarshal err: %v", err)
		}
		switch result.OperType {
		case "ADD":
			err = action.Add(result.CfgMsg, namespace, service, filename, dst)
		case "EDIT":
			err = action.Edit(result.CfgMsg, namespace, service, filename, dst)
		case "KEEP":
			err = nil
		}
		return err
	}
	return nil
}

type CfgResult struct {
	OperType string `json:"operType"` // 传输配置信息的类型，ADD，EDIT，KEEP
	CfgMsg   string `json:"cfgMsg"`   // 传输配置的信息
}

func Login() (string, error) {
	client := &http.Client{}
	loginMsg := fmt.Sprintf(`{ "username": "%s",
"password": "%s"}`, config2.RequestMsg.Username, config2.RequestMsg.Password)
	buf := bytes.NewBuffer([]byte(loginMsg))
	req, err := http.NewRequest(config2.RequestMsg.Type, config2.RequestMsg.Url, buf)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	returnBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result string
	json.Unmarshal(returnBody, &result)
	return result, nil
}
