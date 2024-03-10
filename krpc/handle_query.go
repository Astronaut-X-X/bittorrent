package krpc

import (
	"fmt"
	"net"

	"bittorrent/utils"
)

func handleQuery(c *Client, m *Message, addr *net.UDPAddr) {
	if err := c.RoutingTable.Add(m.A.Id, addr.IP.String(), addr.Port); err != nil {
		return
	}

	switch m.Q {
	case ping:
		msg := &Message{
			T: m.T,
			Y: r,
			R: &R{
				Id: c.LocalId,
			},
		}
		c.sendMessage(msg, addr)

	case find_node:
		peers := c.RoutingTable.GetPeers(m.A.Target)

		nodes := make([]byte, 0)
		for _, peer := range peers {
			nodes = append(nodes, utils.ParseIdToByte(peer.Id)...)
			nodes = append(nodes, utils.ParseIpPortToByte(peer.Ip, peer.Port)...)
		}

		msg := &Message{
			T: m.T,
			Y: r,
			R: &R{
				Id:    c.LocalId,
				Nodes: string(nodes),
			},
		}
		c.sendMessage(msg, addr)

	case get_peers:
		peers := c.RoutingTable.GetPeers(m.A.Target)

		nodes := make([]byte, 0)
		for i, peer := range peers {
			if i == 7 {
				break
			}
			nodes = append(nodes, utils.ParseIdToByte(peer.Id)...)
			nodes = append(nodes, utils.ParseIpPortToByte(peer.Ip, peer.Port)...)
		}

		msg := &Message{
			T: m.T,
			Y: r,
			R: &R{
				Id:    c.LocalId,
				Nodes: string(nodes),
				Token: utils.RandomToken(),
			},
		}
		c.sendMessage(msg, addr)

	case announce_peer:
		fmt.Println("info_hash", m.A.InfoHash)
		fmt.Println("port", m.A.Port)
		fmt.Println("token", m.A.Token)
		fmt.Println("implied_port", m.A.ImpliedPort)
	}
}
