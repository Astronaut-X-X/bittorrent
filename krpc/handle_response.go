package krpc

import (
	"bittorrent/utils"
	"fmt"
	"net"
)

func handleResponse(c *Client, m *Message, addr *net.UDPAddr) {
	if m.R == nil || m.R.Id == "" {
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
		for _, node := range handleNodes(m) {
			c.HandleNode(node)
		}

	case get_peers:
		if len(m.R.Nodes) > 0 {
			for _, node := range handleNodes(m) {
				c.HandleNode(node)

				c.GetPeersContinuous(node, m.T, transaction.Query.A.InfoHash)
			}
		}
		if len(m.R.Values) > 0 {
			handleValues(c, m, transaction.Query)
		}

	case announce_peer:
		// do nothing
	}

	//transaction.Response <- true
	c.TransactionManager.Delete(transaction)
}

func handleNodes(m *Message) []*Node {
	nodes, err := ParseNodes([]byte(m.R.Nodes))
	if err != nil {
		return nil
	}

	return nodes
}

func handleValues(c *Client, m *Message, q *Message) {
	values := make([]interface{}, 0)
	if m.R.Values != nil {
		values = m.R.Values
	}

	for _, value := range values {
		byteValue := []byte(value.(string))
		if len(byteValue) != 6 {
			continue
		}
		ip, port, err := utils.ParseByteToIpPort(byteValue)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		node := NewPeer(ip, port, q.A.InfoHash)
		c.HandleValue(node)
	}
}
