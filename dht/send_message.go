package dht

import (
	"bittorrent/utils"
	"fmt"
	"net"
)

func sendMessage(d *DHT, msg *Message, addr *net.UDPAddr) bool {
	msgByte := EncodeMessage(msg)

	d.log.Println("[send]", msgByte)
	d.log.Println("[send]", string(msgByte))

	n, err := d.Conn.WriteToUDP(msgByte, addr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("send :", n)
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

	t := NewTransaction(msg.T, msg)
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

	sendMessage(d, msg, udpAddr)
}

func GetPeers(d *DHT, addr string, infoHash string) {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msg := &Message{
		T: utils.RandomT(),
		Y: q,
		Q: get_peers,
		A: &A{
			Id:       d.routingTable.LocalId,
			InfoHash: infoHash,
		},
	}

	sendMessage(d, msg, udpAddr)
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
