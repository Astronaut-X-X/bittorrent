package utils

import (
	"errors"
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

func ParseAddrToByte(addr *net.UDPAddr) []byte {
	if addr == nil {
		return nil
	}

	return append(addr.IP.To4(), byte(addr.Port/256), byte(addr.Port%256))
}

func ParseByteToIpPort(data []byte) (net.IP, int, error) {
	if len(data) != 6 {
		return nil, 0, errors.New("error data")
	}

	ip := net.IPv4(data[0], data[1], data[2], data[3])
	port := int(data[4])*256 + int(data[5])

	return ip, port, nil
}

func ParseByteToAddr(data []byte) (*net.UDPAddr, error) {
	ip, port, err := ParseByteToIpPort(data)
	if err != nil {
		return nil, err
	}

	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	return addr, nil
}
