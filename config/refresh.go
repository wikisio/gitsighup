package config

import (
	"fmt"
	"os"
	"syscall"
)

func Refresh(c <-chan os.Signal) {
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Ignore panic in refresh: %v\n", err)
				}
			}()

			var s = <-c
			switch s {
			case syscall.SIGHUP:
				var err = LoadConfig()
				if err != nil {
					fmt.Printf("failed to refresh config: %v\n", err)
				}
			}
		}()
	}
}
