package util

import (
	"net"
)

func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrSlice {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if nil != ipnet.IP.To4() {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}