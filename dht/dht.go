package dht

import (
	"context"
	"fmt"
	"net"

	rt "bittorrent/routingTabel"
)

type DHT struct {
	Conn *net.UDPConn

	context context.Context
	cancel  context.CancelFunc

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
	dht.routingTable = rt.NewRoutingTable()
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

func (d *DHT) receiving() {
	buffer := make([]byte, 1024)

	fmt.Println("receiving start")

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

			fmt.Println("receiving")

			go d.process(addr, buffer[:n])
		}
	}

	fmt.Println("receiving done")
}

func (d *DHT) process(addr *net.UDPAddr, data []byte) {
	m := UnmarshalMessage(data)

	fmt.Println(m)

	// Transaction check
	// handle
	if m.Y == "q" {
		d.handleQuery(m)
	}

	if m.Y == "r" {
		d.handleResponse(m)
	}
}

func (d *DHT) handleQuery(m *Message) {
	fmt.Println(m)
}

func (d *DHT) handleResponse(m *Message) {

	if m.R != nil && m.R.Nodes != "" {
		num := len(m.R.Nodes) / (20 + 4 + 2)
		for i := 0; i < num; i++ {
			s := i * 26
			eid := s + 20
			id := m.R.Nodes[s:eid]
			ip := net.IPv4(m.R.Nodes[s+21], m.R.Nodes[s+22], m.R.Nodes[s+23], m.R.Nodes[s+24])
			port := int(m.R.Nodes[s+25])*256 + int(m.R.Nodes[s+26])
			d.routingTable.Add(id, ip.String(), port)
		}
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

		message := &Message{
			T: rt.RandLocalId(),
			Y: "q",
			Q: "ping",
			A: &A{
				Id: d.routingTable.LocalId,
			},
		}

		msg_byte, err := MarshalMessage(message)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(string(msg_byte))

		d.Conn.WriteToUDP(msg_byte, addr)
	}

}
