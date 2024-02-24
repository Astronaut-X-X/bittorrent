package dht

import (
	"bittorrent/utils"
	"encoding/hex"
	"net"
)

func handleQuery(d *DHT, m *Message, addr *net.UDPAddr) {
	d.routingTable.Add(m.A.Id, addr.String(), addr.IP.String(), addr.Port)

	switch m.Q {
	case ping:
		msg := &Message{
			T: m.T,
			Y: r,
			R: &R{
				Id: d.routingTable.LocalId,
			},
		}
		sendMessage(d, msg, addr)

	case find_node:
		peers := d.routingTable.GetPeers(m.A.Target)

		nodes := make([]byte, 0)
		for _, peer := range peers {
			nodes = append(nodes, utils.ParseIdToByte(peer.Id)...)
			nodes = append(nodes, utils.ParseIpPortToByte(peer.Ip, peer.Port)...)
		}

		msg := &Message{
			T: m.T,
			Y: r,
			R: &R{
				Id:    d.routingTable.LocalId,
				Nodes: hex.EncodeToString(nodes),
			},
		}
		sendMessage(d, msg, addr)

	case get_peers:
		peers := d.routingTable.GetPeers(m.A.Target)

		nodes := make([]byte, 0)
		for _, peer := range peers {
			nodes = append(nodes, utils.ParseIdToByte(peer.Id)...)
			nodes = append(nodes, utils.ParseIpPortToByte(peer.Ip, peer.Port)...)
		}

		msg := &Message{
			T: m.T,
			Y: r,
			R: &R{
				Id:    d.routingTable.LocalId,
				Nodes: hex.EncodeToString(nodes),
				Token: utils.RandomToken(),
			},
		}
		sendMessage(d, msg, addr)

	case announce_peer:
		// TODO Get mateinfo
	}
}
