package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vishvananda/netlink"
	"net"
	"os"
)

func GetOs() (name string, version string) {
	id := viper.GetString("ID")
	versionId := viper.GetString("VERSION_ID")
	fmt.Println(fmt.Sprintf("%s :: %s", id, versionId))
	return id, versionId
}

func GetIPs(nic net.Interface) (ipv4Addrs []net.Addr, err error) {
	addrs, err := nic.Addrs()
	if err != nil {
		return ipv4Addrs, fmt.Errorf("error getting nic addrs: %s", err)
	}
	for _, addr := range addrs {
		if ipv4Addr, ok := addr.(*net.IPNet); ok && !ipv4Addr.IP.IsLoopback() {
			if ipv4Addr.IP.To4() != nil {
				ipv4Addrs = append(ipv4Addrs, ipv4Addr)
			}
		}
	}
	fmt.Printf("%v", ipv4Addrs)
	return ipv4Addrs, nil
}

func GetDefaultIPv4Route() (*netlink.Route, error) {
	defaultRoutes, err := netlink.RouteListFiltered(netlink.FAMILY_V4, &netlink.Route{Dst: nil}, netlink.RT_FILTER_DST)
	if err != nil {
		return nil, fmt.Errorf("error getting defaultRoutes: %s", err)
	}
	if len(defaultRoutes) == 0 {
		return nil, fmt.Errorf("no default routes found")
	}
	defaultRouteLowestPrio := &defaultRoutes[0]
	for _, route := range defaultRoutes {
		if defaultRouteLowestPrio.Priority > route.Priority {
			defaultRouteLowestPrio = &route
		}
	}
	return defaultRouteLowestPrio, nil
}

func WriteFile(content, path, backupId string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error getting file info: %s", err)
		}
		// File does not exist  yet, ok to create
	} else {
		if fileInfo.IsDir() {
			return fmt.Errorf("%s is a directory", path)
		}
		// TODO Overwrite file or move file to backup
		newPath := fmt.Sprintf("%s.%s", path, backupId)
		err := os.Rename(path, fmt.Sprintf("%s.bak", newPath))
		if err != nil {
			return fmt.Errorf("error renaming file: %s", err)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}
	return nil
}
