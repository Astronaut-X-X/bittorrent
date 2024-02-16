package routingTable

import (
	"time"
)

type Peer struct {
	Id string

	Address string
	Port    int

	AddTime time.Time
}

func NewPeer(id string, address string, port int) *Peer {
	return &Peer{
		Id:      id,
		Address: address,
		Port:    port,
		AddTime: time.Now(),
	}
}
