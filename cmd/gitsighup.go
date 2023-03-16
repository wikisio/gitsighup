package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	config2 "wikis.io/config"
)

func main() {
	var configFile = flag.String("c", "", "the yaml config file")
	flag.Parse()

	var f, err = os.Open(*configFile)
	if err != nil {
		fmt.Printf("Failed to open config file, %s, error: %v", configFile, err)
		os.Exit(1)
	}

	var config config2.Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		fmt.Printf("Failed to read yaml content, error: %v", err)
		os.Exit(2)
	}

	fmt.Printf("%v", config)

}
