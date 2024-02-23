package dht

import (
	"context"
	"fmt"
	"net"

	rt "bittorrent/routingTabel"
)

type DHT struct {
	Conn *net.UDPConn

	config *config

	context context.Context
	cancel  context.CancelFunc

	tm           *TransactionManager
	routingTable *rt.RoutingTable
}

func NewDHT(c *config) (*DHT, error) {
	dht := &DHT{}

	addr, err := net.ResolveUDPAddr("udp", c.Address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	dht.Conn = conn
	dht.Conn.SetReadBuffer(c.ReadBuffer)
	dht.config = c
	dht.tm = NewTransactionManager()

	dht.routingTable = rt.NewRoutingTable()
	dht.routingTable.SetPingPeer(func(addr string) bool {
		return <-Ping(dht, addr)
	})

	dht.context, dht.cancel = context.WithCancel(context.Background())

	return dht, nil
}

func (d *DHT) Run() {
	go d.sendPrimeNodes()
	d.receiving()

}

func (d *DHT) Stop() {
	d.cancel()

	if d.Conn != nil {
		d.Conn.Close()
	}
}

func (d *DHT) sendPrimeNodes() {

	for _, node := range DefualtConfig().PrimeNodes {

		addr, err := net.ResolveUDPAddr("udp", node)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(addr.IP, addr.Port)

		// message := &Message{
		// 	T: rt.RandLocalId(),
		// 	Y: "q",
		// 	Q: "find_node",
		// 	A: &A{
		// 		Id:     d.routingTable.LocalId,
		// 		Target: rt.RandLocalId(),
		// 	},
		// }

		msg := &Message{
			T: rt.NewRoutingTable().LocalId,
			Y: "q",
			Q: "ping",
			A: &A{
				Id: d.routingTable.LocalId,
			},
		}

		sendMessage(d, msg, addr)
	}

}

func (d *DHT) receiving() {
	buffer := make([]byte, 1024)

out:
	for {
		select {
		case <-d.context.Done():
			break out
		default:
			n, addr, err := d.Conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Printf("Error receiving data: %v\n", err)
				continue
			}

			go d.process(addr, buffer[:n])
		}
	}

	fmt.Println("receiving done")
}

func (d *DHT) process(addr *net.UDPAddr, data []byte) {
	m, err := DecodeMessage(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch m.Y {
	case "q":
		handleQuery(d, m, addr)
	case "r":
		handleResponse(d, m)
	}
}
