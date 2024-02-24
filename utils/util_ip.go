package utils

import (
	"fmt"
	"net"
)

func parseIp(addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(udpAddr.IP, udpAddr.Port)
}
