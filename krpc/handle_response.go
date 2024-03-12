package krpc

import (
	routingTable "bittorrent/routingTabel"
	"fmt"
	"net"
)

func handleResponse(c *Client, m *Message, addr *net.UDPAddr) {
	if err := c.RoutingTable.Add(m.R.Id, addr.IP.String(), addr.Port); err != nil {
		return
	}

	transaction, ok := c.TransactionManager.Load(m.T)
	if !ok {
		return
	}

	switch transaction.Query.Q {
	case ping:
		// do nothing

	case find_node:
		handleNodes(c, m)

	case get_peers:
		if len(m.R.Nodes) > 0 {
			handleNodes(c, m)
			c.GetPeers(transaction.Query.A.InfoHash)
		}
		if len(m.R.Values) > 0 {
			handleValues(c, m)
		}

	case announce_peer:
		// do nothing

	}

	transaction.Response <- true
	c.TransactionManager.Delete(transaction)
}

func handleNodes(c *Client, m *Message) {
	peers := make([]*routingTable.Peer, 0)

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

	c.RoutingTable.AddPeers(peers)

	go c.OnHandleNodes(peers)
}

func handleValues(c *Client, m *Message) {
	values := make([]string, 0)
	if m.R.Values != nil {
		values = m.R.Values
	}

	for _, value := range values {
		// TODO get meta info
		fmt.Println(value)
	}

}
