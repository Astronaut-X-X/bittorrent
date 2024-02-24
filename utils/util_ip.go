package utils

import (
	"fmt"
	"net"
)

func ParseIp(addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(udpAddr.IP, udpAddr.Port)
}

func ParseIpPortToByte(ip_str string, port int) []byte {
	ip := net.ParseIP(ip_str)

	port_1 := byte(port)
	port_2 := byte(port << 8)
	return append(ip, port_2, port_1)
}
