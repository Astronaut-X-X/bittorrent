package krpc

import (
	"bittorrent/utils"
	"fmt"
	"net"
)

func handleResponse(c *Client, m *Message, addr *net.UDPAddr) {
	node := NewNode(m.R.Id, addr)
	c.HandleNode(node)

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
			nodes := handleNodes(c, m)
			for _, node := range nodes {
				c.GetPeers(node.Addr.String(), transaction.Query.A.InfoHash)
			}
		}
		if len(m.R.Values) > 0 {
			handleValues(c, m)
		}

	case announce_peer:
		// do nothing

	}

	//transaction.Response <- true
	c.TransactionManager.Delete(transaction)
}

func handleNodes(c *Client, m *Message) []*Node {
	nodes := make([]*Node, 0)
	length := len(m.R.Nodes)
	for i := 0; i < length; i += 26 {
		id := m.R.Nodes[i+20 : i+20]
		data := []byte(m.R.Nodes[i+20 : i+26])
		addr := utils.ParseByteToAddr(data)
		node := NewNode(id, addr)
		c.HandleNode(node)
		nodes = append(nodes, node)
	}
	return nodes
}

func handleValues(c *Client, m *Message) {
	values := make([]interface{}, 0)
	if m.R.Values != nil {
		values = m.R.Values
	}

	for _, value := range values {
		// TODO get meta info
		fmt.Println(value)
	}
}
