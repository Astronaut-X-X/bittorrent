package krpc

import (
	"bittorrent/logger"
	"bittorrent/utils"
	"fmt"
	"net"
)

func (c *Client) sendMessage(msg *Message, addr *net.UDPAddr) {
	msgByte := EncodeMessage(msg)
	if _, err := c.WriteToUDP(msgByte, addr); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c *Client) sendAndStore(msg *Message, addr *net.UDPAddr) {
	c.sendMessage(msg, addr)
	c.TransactionManager.Store(NewTransaction(msg))
}

func (c *Client) Ping(node *Node) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: ping,
		A: &A{
			Id: c.LocalId,
		},
	}

	c.sendAndStore(msg, node.Addr)
}

func (c *Client) FindNode(node *Node, target string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: find_node,
		A: &A{
			Id:     c.LocalId,
			Target: target,
		},
	}

	c.sendAndStore(msg, node.Addr)
}

func (c *Client) GetPeers(nodes []*Node, infoHash string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	if len(nodes) == 0 {
		logger.Println("[GetPeers] nodes empty")
		return
	}
	node := nodes[0]
	c.sendMessage(msg, node.Addr)

	transaction := NewTransaction(msg)
	transaction.NodeQueue.PushNodes(nodes[1:])
	c.TransactionManager.Store(transaction)
}

func (c *Client) GetPeersContinuous(queue *NodeQueue, infoHash string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	if queue.Len() == 0 {
		logger.Println("[GetPeers] nodes empty")
		return
	}
	node := queue.Pop()
	c.sendMessage(msg, node.Addr)

	transaction := NewTransaction(msg)
	transaction.NodeQueue = queue
	c.TransactionManager.Store(transaction)
}

// AnnouncePeer TODO
func (c *Client) AnnouncePeer(node *Node, infoHash string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: announce_peer,
		A: &A{
			Id:          c.LocalId,
			InfoHash:    infoHash,
			ImpliedPort: 0,
			Port:        6881,
			Token:       "XX",
		},
	}

	c.sendAndStore(msg, node.Addr)
}
