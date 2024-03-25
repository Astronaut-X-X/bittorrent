package routing

import (
	"fmt"
	"net"
	"time"
)

type Node struct {
	NodeId string
	Ip     string
	Port   int
	Addr   *net.UDPAddr
	Create time.Time
}

func NewNode(nodeId string, ip string, port int) *Node {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil
	}

	return &Node{
		NodeId: nodeId,
		Ip:     ip,
		Port:   port,
		Addr:   addr,
		Create: time.Now(),
	}
}
