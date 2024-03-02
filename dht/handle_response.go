package dht

import (
	"fmt"
	"net"
)

func handleResponse(d *DHT, m *Message, addr *net.UDPAddr) {
	if err := d.routingTable.Add(m.R.Id, addr.IP.String(), addr.Port); err != nil {
		return
	}

	if m.R.Nodes != "" {
		length := len(m.R.Nodes)
		for i := 0; i < length; i += 26 {
			id := m.R.Nodes[i : i+20]
			ip := net.IPv4(m.R.Nodes[i+20], m.R.Nodes[i+21], m.R.Nodes[1+22], m.R.Nodes[1+23])
			port := int(m.R.Nodes[i+24])*256 + int(m.R.Nodes[i+25])

			if err := d.routingTable.Add(id, ip.String(), port); err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Println("[Nodes]", ip.String(), port)
			d.log.Println("[Nodes]", ip.String(), port)
		}
	}

	if m.R.Values != nil {
		// TODO send get mateinfo
	}
}
