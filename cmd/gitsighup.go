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
				url := "http://127.0.0.1:3000/api/v1/configsrv/"
				SendRequest(url, namespace, service, filename, dst, jwt)
			}

		}
	}

}

func SendRequest(url string, namespace string, service string, filename string, dst string, jwt string) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("connect failed:%v", err)
		}
	}()
	url = url + namespace + "/" + service + "/" + filename
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//登录
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var result CfgResult
		json.Unmarshal(body, &result)
		switch result.OperType {
		case "ADD":
			action.ADD(result.CfgMsg, namespace, service, filename, dst)
		case "EDIT":
			action.EDIT(result.CfgMsg, namespace, service, filename, dst)
		case "KEEP":
		}
	}

}

type CfgResult struct {
	OperType string `json:"operType"` // 传输配置信息的类型，ADD，EDIT，KEEP
	CfgMsg   string `json:"cfgMsg"`   // 传输配置的信息
}

func Login() (string, error) {
	client := &http.Client{}
	buf := bytes.NewBuffer([]byte(`{ "username": "testuser",
"password": "123456"}`))
	req, err := http.NewRequest("POST", "http://127.0.0.1:3000/api/v1/session/login", buf)
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
