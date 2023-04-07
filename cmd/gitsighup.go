package main

import (
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

func main() {
	var configFile = flag.String("c", "empty", "the yaml config file")
	flag.Parse()
	config2.GlobalConfigFile = *configFile

	var err = config2.LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	var c = make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP)
	go config2.Refresh(c)

	var namespace string
	var service string
	var filename string
	var dst string
	for _, i := range config2.GlobalConfig.Services {
		for _, j := range i.ConfigPath {
			filename = j.Src
			dst = j.Dst
		}
		namespace = i.NameSpace
		service = i.Name
	}
	url := ""
	SendRequest(url, namespace, service, filename, dst)

}

func SendRequest(url string, namespace string, service string, filename string, dst string) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("check login table error:%v", err)
		}
	}()
	resp, err := http.Get(url)
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
			action.ADD(namespace, service, filename, dst)
		case "EDIT":
			action.EDIT(namespace, service, filename, dst)
		case "KEEP":
		}
	}

}

type CfgResult struct {
	OperType string `json:"operType"` // 传输配置信息的类型，ADD，EDIT，KEEP
	CfgMsg   string `json:"cfgMsg"`   // 传输配置的信息
}
