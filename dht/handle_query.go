package dht

import "net"

func handleQuery(d *DHT, m *Message, addr *net.UDPAddr) {
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

	case get_peers:

	case announce_peer:

	}
}
