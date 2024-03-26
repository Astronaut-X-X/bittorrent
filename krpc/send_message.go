package krpc

import (
	"fmt"
	"net"

	"bittorrent/logger"
	"bittorrent/utils"
)

func (c *Client) sendMessageAddr(msg *Message, addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.sendMessage(msg, udpAddr)

	c.TransactionManager.Store(NewTransaction(msg))

	logger.Println("[SEND]", addr, Print(msg))
}

func (c *Client) sendMessageContinuous(msg *Message, addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.sendMessage(msg, udpAddr)

	logger.Println("[SEND]", addr, Print(msg))
}

func (c *Client) sendMessage(msg *Message, addr *net.UDPAddr) {
	msgByte := EncodeMessage(msg)

	if _, err := c.WriteToUDP(msgByte, addr); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c *Client) Ping(addr string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: ping,
		A: &A{
			Id: c.LocalId,
		},
	}

	c.sendMessageAddr(msg, addr)
}

func (c *Client) FindNode(addr string, target string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: find_node,
		A: &A{
			Id:     c.LocalId,
			Target: target,
		},
	}

	c.sendMessageAddr(msg, addr)
}

func (c *Client) GetPeers(addr string, infoHash string) {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	c.sendMessageAddr(msg, addr)
}

func (c *Client) GetPeersContinuous(addr string, T string, infoHash string) {
	msg := &Message{
		T: T,
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	c.sendMessageContinuous(msg, addr)
}

// AnnouncePeer TODO
func (c *Client) AnnouncePeer(addr string, infoHash string) {
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

	c.sendMessageAddr(msg, addr)
}
