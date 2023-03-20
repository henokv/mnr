package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"strings"
)

type Linux struct {
	nicPrefix string
	//tempSuffix string
}

func NewLinux(nicPrefix string) *Linux {
	linux := Linux{nicPrefix: nicPrefix}
	return &linux
}

//func (linux *Linux) getTempSuffix() string {
//	if len(linux.tempSuffix) == 0 {
//		linux.tempSuffix = fmt.Sprintf("%d", time.Now().Unix())
//	}
//	return linux.tempSuffix
//}

func (linux *Linux) getRelevantNics() (filteredNics []net.Interface, err error) {
	nics, err := net.Interfaces()
	if err != nil {
		return nics, err
	}
	for _, nic := range nics {
		if strings.HasPrefix(nic.Name, linux.GetNicPrefix()) {
			filteredNics = append(filteredNics, nic)
		}
	}
	return filteredNics, err
}

func (linux *Linux) IfRuleFileContent(nic net.Interface, ipv4Addrs []net.Addr) string {
	path := fmt.Sprintf("/etc/sysconfig/network-scripts/rule-%s", nic.Name)
	ipv4Addr, _, err := net.ParseCIDR(ipv4Addrs[0].String())
	if err != nil {
		log.Fatalf("error parsing CIDR: %s", err)
	}
	name := nic.Name
	content := fmt.Sprintf(`from %s/32 table %s-rt
to %s/32 table %s-rt`, ipv4Addr, name, ipv4Addr, name)
	fmt.Printf("%s", path)
	return content
}

func (linux *Linux) IfRouteFileContent(nic net.Interface, ipv4Addrs []net.Addr, defaultRoute *netlink.Route) string {
	path := fmt.Sprintf("/etc/sysconfig/network-scripts/route-%s", nic.Name)
	_, network, err := net.ParseCIDR(ipv4Addrs[0].String())
	if err != nil {
		log.Fatalf("error parsing CIDR: %s", err)
	}
	name := nic.Name
	gw := defaultRoute.Gw
	content := fmt.Sprintf(`%s dev %s table %s-rt
default via %s dev %s table %s-rt`, network.String(), name, name, gw, name, name)
	fmt.Printf("%s", path)
	return content
}

func (linux *Linux) GlobalRouteTablesFileContent(nics []net.Interface) {
	//path := "/etc/sysconfig/network-scripts/rule-eth0"
	path := "/etc/iproute2/rt_tables"
	content := ""
	for _, nic := range nics {
		id := strings.ReplaceAll(nic.Name, "wlp59s", "")
		content += fmt.Sprintf("2%02s %s-rt\n", id, nic.Name)
	}
	fmt.Printf("%s", path)
	fmt.Print(content)
}

func (linux *Linux) GetNicPrefix() string {
	return linux.nicPrefix
}
