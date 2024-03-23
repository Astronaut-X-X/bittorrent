package krpc

import (
	"bittorrent/logger"
	"net"

	"bittorrent/utils"
)

func (c *Client) sendMessage(msg *Message, addr *net.UDPAddr) bool {
	msgByte := EncodeMessage(msg)

	if _, err := c.WriteToUDP(msgByte, addr); err != nil {
		return false
	}

	logger.Println("[SEND]", Print(msg))

	return true
}

func (c *Client) Ping(addr *net.UDPAddr) chan bool {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: ping,
		A: &A{
			Id: c.LocalId,
		},
	}

	if !c.sendMessage(msg, addr) {
		return nil
	}

	t := NewTransaction(msg, func(t *Transaction) { t.Response <- false })
	c.TransactionManager.Store(t)

	return t.Response
}

func (c *Client) FindNode(target string) chan bool {
	peer := c.RoutingTable.GetPeer(target)
	if peer == nil {
		return nil
	}

	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: find_node,
		A: &A{
			Id:     c.LocalId,
			Target: target,
		},
	}

	if !c.sendMessage(msg, peer.Addr) {
		return nil
	}

	t := NewTransaction(msg, func(t *Transaction) { t.Response <- false })
	c.TransactionManager.Store(t)

	return t.Response
}

func (c *Client) GetPeers(infoHash string) chan bool {
	peer := c.RoutingTable.GetPeer(infoHash)
	if peer == nil {
		return nil
	}

	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       c.LocalId,
			InfoHash: infoHash,
		},
	}

	if !c.sendMessage(msg, peer.Addr) {
		return nil
	}

	t := NewTransaction(msg, func(t *Transaction) { t.Response <- false })
	c.TransactionManager.Store(t)

	return t.Response
}

// AnnouncePeer TODO
func (c *Client) AnnouncePeer(addr *net.UDPAddr, infoHash string) chan bool {
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

	if !c.sendMessage(msg, addr) {
		return nil
	}

	t := NewTransaction(msg, func(t *Transaction) { t.Response <- false })
	c.TransactionManager.Store(t)

	return t.Response
}
