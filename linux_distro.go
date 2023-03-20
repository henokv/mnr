package main

import (
	"fmt"
	"strings"
)

var OsNotSupported = fmt.Errorf("OS is not supported")

type LinuxDistro interface {
	GetNicPrefix() string
	WriteFiles() error
	getTempSuffix() string
}

func NewLinuxDistro() (LinuxDistro, error) {
	name, version := GetOs()
	switch name {
	case "rhel":
		if strings.HasPrefix(version, "8.") {
			return NewRhel("8")
		} else if strings.HasPrefix(version, "7.") {
			return NewRhel("7")
		}
	case "ubuntu":
		if strings.HasPrefix(version, "20.") || strings.HasPrefix(version, "18.") {
			return nil, OsNotSupported // TODO: implement for Ubuntu 18/20
		} else if strings.HasPrefix(version, "22.") {
			return nil, OsNotSupported // TODO: implement for Ubuntu 22 no docs?
		}
	case "debian":
		if strings.HasPrefix(version, "10.") {
			return nil, OsNotSupported // TODO: implement for Debian 10
		}
	case "suse":
		if strings.HasPrefix(version, "15.") || strings.HasPrefix(version, "12.") || strings.HasPrefix(version, "11.") {
			return nil, OsNotSupported // TODO: implement for SUSE 11/12/15
		}
	}
	return nil, OsNotSupported
}
