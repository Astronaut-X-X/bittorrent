package routingTable

import (
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	PeerIdError = errors.New("error peer id")
)

// Peer is a tracker in DHT
type Peer struct {
	Id         string       // length of 20 byte
	Ip         string       // Peer ip
	Port       int          // Peer port
	Addr       *net.UDPAddr // Peer Address
	CreateTime time.Time    // Create time
}

// NewPeer Create a Peer
func NewPeer(id string, ip string, port int) (*Peer, error) {

	if len(id) != 20 {
		return nil, PeerIdError
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	return &Peer{
		Id:         id,
		Ip:         ip,
		Port:       port,
		Addr:       addr,
		CreateTime: time.Now(),
	}, nil
}
