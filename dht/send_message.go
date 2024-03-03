package dht

import (
	routingTable "bittorrent/routingTabel"
	"bittorrent/utils"
	"fmt"
	"net"
)

func sendMessage(d *DHT, msg *Message, addr *net.UDPAddr) bool {

	msgByte := EncodeMessage(msg)

	d.log.Println("[send]", msgByte)

	_, err := d.Conn.WriteToUDP(msgByte, addr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func Ping(d *DHT, addr *net.UDPAddr) chan bool {
	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: ping,
		A: &A{
			Id: d.routingTable.LocalId,
		},
	}

	if !sendMessage(d, msg, addr) {
		return nil
	}

	t := NewTransaction(msg.T, d, msg, func(t *Transaction) { t.Response <- false })
	d.tm.Store(t)

	return t.Response
}

func FindNode(d *DHT, addr string, target string) {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: find_node,
		A: &A{
			Id:     d.routingTable.LocalId,
			Target: target,
		},
	}

	t := NewTransaction(msg.T, d, msg, func(t *Transaction) { t.Response <- false })
	d.tm.Store(t)

	sendMessage(d, msg, udpAddr)
}

func GetPeers(d *DHT, addr *net.UDPAddr, infoHash string, peers []*routingTable.Peer) {

	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       d.routingTable.LocalId,
			InfoHash: infoHash,
		},
	}

	fmt.Println("[GetPeers]", infoHash, addr.String())

	if _, ok := d.tm.Load(msg.T); !ok {
		t := NewTransaction(msg.T, d, msg, func(t *Transaction) {
			if len(t.Peers) == 0 {
				t.Response <- false
				return
			}

			peer := t.Peers[0]
			t.Peers = t.Peers[1:]
			GetPeers(t.DHT, peer.Addr, msg.A.InfoHash, nil)
		})
		t.Peers = append(t.Peers, peers...)
		d.tm.Store(t)
	}

	sendMessage(d, msg, addr)
}

func AnnouncePeer(d *DHT, addr string, infoHash string) {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: announce_peer,
		A: &A{
			Id:          d.routingTable.LocalId,
			InfoHash:    infoHash,
			ImpliedPort: 0,
			Port:        d.config.Port,
			Token:       "XX",
		},
	}

	sendMessage(d, msg, udpAddr)
}
