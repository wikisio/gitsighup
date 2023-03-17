package config

import (
	"fmt"
	"os"
	"syscall"
)

func Refresh(c <-chan os.Signal) {
	for {
		var s = <-c
		switch s {
		case syscall.SIGHUP:
			var err = LoadConfig()
			if err != nil {
				fmt.Printf("failed to refresh config: %v", err)
			}
		}
	}
}
