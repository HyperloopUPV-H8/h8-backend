package sniffer

import (
	"fmt"
	"strings"
)

func getFilters(srcAddrs []string, dstAddrs []string) string {
	result := "udp && ("
	result += getFilterHosts("src", srcAddrs) + ") || ("
	return "udp && (" + getFilterHosts("src", srcAddrs) + ") && (" + getFilterHosts("dst", dstAddrs) + ")"
}

func getFilterHosts(dir string, addrs []string) (result string) {
	for _, addr := range addrs {
		result += fmt.Sprintf("%s host %s || ", dir, addr)
	}
	return strings.TrimSuffix(result, " || ")
}
