package util

import (
	"net"
	"time"
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

func GetCurrentTimestamp() int64 {
	return int64(float64(time.Now().UnixNano()) / 1e6)
}