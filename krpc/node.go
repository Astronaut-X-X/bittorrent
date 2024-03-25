package krpc

import (
	"bittorrent/utils"
	"net"
)

type Node struct {
	Id   string
	Addr *net.UDPAddr
}

func NewNode(id string, addr *net.UDPAddr) *Node {
	return &Node{
		Id:   id,
		Addr: addr,
	}
}

func (n *Node) toByte() []byte {
	return utils.ParseAddrToByte(n.Addr)
}
