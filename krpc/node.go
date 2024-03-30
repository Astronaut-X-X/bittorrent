package krpc

import (
	"bittorrent/utils"
	"errors"
	"fmt"
	"net"
)

const (
	NoNeedAppendQueue = iota
	NeedAppendQueue   = iota
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

func ParseNode(data []byte) (*Node, error) {
	if len(data) != 26 {
		return nil, errors.New("error data length")
	}

	id := data[:20]
	addr, err := utils.ParseByteToAddr(data[20:26])
	if err != nil {
		return nil, errors.New("error length of data")
	}
	node := NewNode(string(id), addr)

	return node, nil
}

func ParseNodes(data []byte) ([]*Node, error) {
	if len(data)%26 != 0 {
		return nil, errors.New("error length of data")
	}

	nodeMap := make(map[string]*Node)
	for i := 0; i < len(data); i += 26 {
		node, err := ParseNode(data[i : i+26])
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		nodeMap[node.Id] = node
	}

	nodes := make([]*Node, 0, len(nodeMap))
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}

	return nodes, nil
}
