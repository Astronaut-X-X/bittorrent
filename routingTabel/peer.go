package routingTable

import (
	"net"
	"time"
)

type Peer struct {
	Id string

	Address net.Addr
	Port    int

	AddTime time.Time
}

func NewPeer(id string, address net.Addr, port int) *Peer {
	return &Peer{
		Id:      id,
		Address: address,
		Port:    port,
		AddTime: time.Now(),
	}
}
