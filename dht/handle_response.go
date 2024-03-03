package dht

import (
	routingTable "bittorrent/routingTabel"
	"fmt"
	"net"
)

func handleResponse(d *DHT, m *Message, addr *net.UDPAddr) {
	if err := d.routingTable.Add(m.R.Id, addr.IP.String(), addr.Port); err != nil {
		return
	}

	peers := make([]*routingTable.Peer, 0)

	if m.R.Nodes != "" {
		length := len(m.R.Nodes)
		for i := 0; i < length; i += 26 {
			id := m.R.Nodes[i : i+20]
			ip := net.IPv4(m.R.Nodes[i+20], m.R.Nodes[i+21], m.R.Nodes[1+22], m.R.Nodes[1+23])
			port := int(m.R.Nodes[i+24])*256 + int(m.R.Nodes[i+25])

			peer, err := routingTable.NewPeer(id, ip.String(), port)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			peers = append(peers, peer)
		}

		d.routingTable.AddPeers(peers)
	}

	values := make([]string, 0)
	if m.R.Values != nil {
		d.log.Println("[values]", m.R.Values)
		// TODO send get mateinfo

	}

	if t, ok := d.tm.Load(m.T); ok {
		fmt.Println("[transaction]", m.T)
		if t.Query.Q == ping {
			t.Response <- true
			t.timer.Stop()
			d.tm.Delete(t)
		}

		if t.Query.Q == find_node {
			t.Response <- true
			t.timer.Stop()
			d.tm.Delete(t)
		}

		if t.Query.Q == get_peers {
			if len(peers) != 0 {
				t.Peers = append(t.Peers, peers...)
				peer := t.Peers[0]
				t.Peers = t.Peers[1:]

				GetPeers(t.DHT, peer.Addr, t.Query.A.InfoHash, nil)
				t.timer.Reset(Timeout)
			}

			if len(values) != 0 {
				t.Response <- true
				t.timer.Stop()
				d.tm.Delete(t)
			}
		}
	}
}
