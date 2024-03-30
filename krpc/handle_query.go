package krpc

import (
	"net"

	"bittorrent/utils"
)

func handleQuery(c *Client, m *Message, addr *net.UDPAddr) {
	if m.A == nil || m.A.Id == "" {
		return
	}
	node := NewNode(m.A.Id, addr)
	c.HandleNode(node, NoNeedAppendQueue)

	msg := &Message{
		T: m.T,
		Y: r,
		R: &R{
			Id: c.LocalId,
		},
	}

	switch m.Q {
	case ping:
		c.sendMessage(msg, addr)

	case find_node:
		nodes := c.SearchNode(m.A.Target)
		byteData := make([]byte, 0)
		for _, node_ := range nodes {
			byteData = append(byteData, node_.toByte()...)
		}
		msg.R.Nodes = string(byteData)
		c.sendMessage(msg, addr)

	case get_peers:
		nodes := c.SearchNode(m.A.Target)
		byteData := make([]byte, 0)
		for _, node_ := range nodes {
			byteData = append(byteData, node_.toByte()...)
		}
		msg.R.Nodes = string(byteData)
		msg.R.Token = utils.RandomToken()
		c.sendMessage(msg, addr)
		// on get_peers
		c.OnGetPeers(node, m)

	case announce_peer:
		//fmt.Println("info_hash", m.A.InfoHash)
		//fmt.Println("port", m.A.Port)
		//fmt.Println("token", m.A.Token)
		//fmt.Println("implied_port", m.A.ImpliedPort)
		c.OnAnnouncePeer(node, m)
	}
}
