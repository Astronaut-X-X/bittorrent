package krpc

import (
	"fmt"
	"net"

	"bittorrent/utils"
)

func (c *Client) sendMessage(msg *Message, addr *net.UDPAddr) {
	msgByte := EncodeMessage(msg)

	if _, err := c.WriteToUDP(msgByte, addr); err != nil {
		fmt.Println(err.Error())
		return
	}
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

	c.SendQueue <- NewQueueMessage(msg, node)
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

	c.SendQueue <- NewQueueMessage(msg, node)
}

func (c *Client) GetPeers(node *Node, infoHash string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	c.SendQueue <- NewQueueMessage(msg, node)
}

func (c *Client) GetPeersContinuous(node *Node, T string, infoHash string) {
	msg := &Message{
		T: T,
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	c.SendQueue <- NewQueueMessage(msg, node)
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

	c.SendQueue <- NewQueueMessage(msg, node)
}
