package main

import (
	"github.com/spf13/viper"
	"log"
	"os/user"
	"runtime"
)

func main() {
	linuxDistro, err := NewLinuxDistro()
	if err != nil {
		log.Fatalf("error creating new linux: %s", err)
	}
	linuxDistro.WriteFiles()
}

func init() {
	goos := runtime.GOOS
	if goos != "linux" {
		log.Fatalf("Unsupported OS: %s", goos)
	}
	viper.SetConfigType("env")
	viper.SetConfigFile("/etc/os-release")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %s", err)
	}
	if user.Uid != "0" {
		log.Fatalf("Must be run as root")
	}

}
