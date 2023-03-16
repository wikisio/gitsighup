package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"wikis.io/action"

	config2 "wikis.io/config"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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

	r := gin.Default()
	r.GET("/api/v1/services/:name", func(c *gin.Context) {
		var serviceName = c.Param("name")
		var tag = c.Query("tag")

		var path string
		for _, i := range config.Services {
			if i.Name == serviceName {
				path = i.Name
				break
			}
		}

		if path == "" {
			c.Status(http.StatusNotFound)
			return
		}

		err = action.Action(serviceName, tag, path)
		if err == nil {
			c.Status(http.StatusOK)
			return
		}

		c.Status(http.StatusBadRequest)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
