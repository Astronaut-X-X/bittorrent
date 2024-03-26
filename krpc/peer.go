package krpc

import (
	"net"
)

type Peer struct {
	Ip       net.IP
	Port     int
	InfoHash string
}

func NewPeer(ip net.IP, port int, infoHash string) *Peer {
	return &Peer{
		Ip:       ip,
		Port:     port,
		InfoHash: infoHash,
	}
}
