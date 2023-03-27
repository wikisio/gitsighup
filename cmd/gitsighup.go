package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wikis.io/action"

	config2 "wikis.io/config"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.PUT("/api/v1/services/:name", func(c *gin.Context) {
		var serviceName = c.Param("name")
		var tag = c.Query("tag")

		path, ok := getPath(c, config2.GlobalConfig, serviceName)
		if !ok {
			return
		}

		err = action.Action(serviceName, path, tag)
		if err == nil {
			c.Status(http.StatusOK)
			return
		}

		err = action.Restart(serviceName, path)

		c.JSON(http.StatusBadRequest, map[string]string{
			"code":    "3005",
			"message": fmt.Sprintf("failed to update configuration: %v", err),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func getPath(c *gin.Context, config *config2.Config, serviceName string) (string, bool) {
	var path string
	for _, i := range config.Services {
		if i.Name == serviceName {
			path = i.ConfigPath
			break
		}
	}

	if path == "" {
		c.JSON(http.StatusNotFound, map[string]string{
			"code":    "3004",
			"message": fmt.Sprintf("Unknown sevice '%s'", serviceName),
		})
		return "", false
	}
	return path, true
}
