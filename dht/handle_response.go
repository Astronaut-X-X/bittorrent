package dht

import (
	"fmt"
	"net"
)

func handleResponse(d *DHT, m *Message, addr *net.UDPAddr) {
	d.routingTable.Add(m.R.Id, addr.String(), addr.IP.String(), addr.Port)

	if m.R == nil {
		return
	}
	if m.R.Nodes != "" {
		length := len(m.R.Nodes)
		for i := 0; i < length; i += 26 {
			id := m.R.Nodes[i : i+20]
			ip := net.IPv4(m.R.Nodes[i+20], m.R.Nodes[i+21], m.R.Nodes[1+22], m.R.Nodes[1+23])
			port := int(m.R.Nodes[i+24])*256 + int(m.R.Nodes[i+25])

			d.routingTable.Add(id, ip.String(), ip.String(), port)

			fmt.Println(string(id), ip.String(), port)
			d.log.Println(string(id), ip.String(), port)
		}
	}
	if m.R.Values != nil {
		// TODO send get mateinfo
	}
}
