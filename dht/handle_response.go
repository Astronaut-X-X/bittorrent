package dht

import "net"

func handleResponse(d *DHT, m *Message, addr *net.UDPAddr) {
	if m.R == nil {
		return
	}
	if m.R.Nodes != "" {
		num := len(m.R.Nodes) / (20 + 4 + 2)
		for i := 0; i < num; i++ {
			s := i * 26
			eid := s + 20
			id := m.R.Nodes[s:eid]
			ip := net.IPv4(m.R.Nodes[s+21], m.R.Nodes[s+22], m.R.Nodes[s+23], m.R.Nodes[s+24])
			port := int(m.R.Nodes[s+25])*256 + int(m.R.Nodes[s+26])
			d.routingTable.Add(id, ip.String(), ip.String(), port)

			d.log.Println(string(id), ip.String(), port)
		}
	}
	if m.R.Values != nil {
		// TODO send get mateinfo
	}
}
