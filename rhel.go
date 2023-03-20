package main

import (
	"fmt"
	"net"
	"time"
)

type Rhel struct {
	*Linux
	tempSuffix string
	nmControl  string
}

func NewRhel(version string) (*Rhel, error) {
	var rh Rhel
	switch version {
	case "8":
		rh = Rhel{NewLinux("eth"), "eth", "yes"}
	case "7":
		rh = Rhel{NewLinux("eth"), "eth", "no"}
	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
	return &rh, nil
}

func (rhel *Rhel) GetNicPrefix() string {
	return rhel.nicPrefix
}

func (rhel *Rhel) getTempSuffix() string {
	if len(rhel.tempSuffix) == 0 {
		rhel.tempSuffix = fmt.Sprintf("%d", time.Now().Unix())
	}
	return rhel.tempSuffix
}

func (rhel *Rhel) WriteFiles() error {
	nics, err := rhel.getRelevantNics()
	if err != nil {
		return fmt.Errorf("error getting relevant nics: %s", err)
	}
	//linux.GlobalRouteTablesFileContent(nics)
	defaultRoute, err := GetDefaultIPv4Route()
	for _, nic := range nics {
		ipv4Addrs, err := GetIPs(nic)
		if err != nil {
			return fmt.Errorf("error getting IPs: %s", err)
		}
		err = WriteFile(rhel.IfConfigFileContent(nic, ipv4Addrs), fmt.Sprintf("/etc/sysconfig/network-scripts/ifcfg-%s", nic.Name), rhel.getTempSuffix())
		if err != nil {
			return err
		}
		err = WriteFile(rhel.IfRuleFileContent(nic, ipv4Addrs), fmt.Sprintf("/etc/sysconfig/network-scripts/rule-%s", nic.Name), rhel.getTempSuffix())
		if err != nil {
			return err
		}
		err = WriteFile(rhel.IfRouteFileContent(nic, ipv4Addrs, defaultRoute), fmt.Sprintf("/etc/sysconfig/network-scripts/route-%s", nic.Name), rhel.getTempSuffix())
		if err != nil {
			return err
		}
	}
	return nil
}

func (rhel *Rhel) IfConfigFileContent(nic net.Interface, ipv4Addrs []net.Addr) string {
	path := fmt.Sprintf("/etc/sysconfig/network-scripts/ifcfg-%s", nic.Name)
	content := fmt.Sprintf(`DEVICE=%s
ONBOOT=yes
BOOTPROTO=dhcp
HWADDR=%s
TYPE=Ethernet
USERCTL=no
PEERDNS=yes
IPV6INIT=no
PERSISTENT_DHCLIENT=%s
NM_CONTROLLED=yes`, nic.Name, nic.HardwareAddr, rhel.nmControl)
	fmt.Printf("%s", path)
	return content
}
