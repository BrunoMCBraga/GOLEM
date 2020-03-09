package util

import (
	"errors"
	"fmt"
	"net"
)

func convertIPCIDRStringsToInternalRepresentation(ipsSlice []string) ([]*net.IPNet, error) {

	netBlocks := make([]*net.IPNet, 0)

	for _, cidrString := range ipsSlice {
		_, netBlock, parseCIDRError := net.ParseCIDR(cidrString)
		if parseCIDRError != nil {
			return nil, errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->util->network_util->ConvertIPCIDRStringsToInternalRepresentation->net.ParseCIDR:", parseCIDRError.Error()))
		}

		netBlocks = append(netBlocks, netBlock)
	}
	return netBlocks, nil
}

func IsIPReserved(ipString string, privateIPBlocks []*net.IPNet) bool {

	parsedIP := net.ParseIP(ipString)
	if parsedIP.IsLoopback() || parsedIP.IsLinkLocalUnicast() || parsedIP.IsLinkLocalMulticast() {
		return true
	} else {
		for _, privateIPBlock := range privateIPBlocks {
			if privateIPBlock.Contains(parsedIP) {
				return true
			}
		}
	}
	return false
}

func GetListOfReservedNetCIDRs() ([]*net.IPNet, error) {
	reservedIPsSlice := []string{
		"0.0.0.0/8",
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"100.64.0.0/10",
		"169.254.0.0/16",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"255.255.255.255/32",
		"::1/128",
		"fe80::/10",
		"fc00::/7",
	}
	return convertIPCIDRStringsToInternalRepresentation(reservedIPsSlice)
}
