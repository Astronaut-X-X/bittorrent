package krpc

import (
	"context"
	"fmt"
	"net"

	"bittorrent/routingTabel"
)

const ReadBuffer int = 10240

type Client struct {
	Conn               *net.UDPConn
	LocalId            string
	context            context.Context
	TransactionManager *TransactionManager
	RoutingTable       *routingTable.RoutingTable

	OnAnnouncePeer func(*Node, *Message)
}

func NewClient(addr string, localId string, ctx context.Context) (*Client, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	if err = conn.SetReadBuffer(ReadBuffer); err != nil {
		return nil, err
	}

	transactionManager := NewTransactionManager()

	table := routingTable.NewRoutingTable(localId, ctx)

	cli := &Client{
		Conn:               conn,
		LocalId:            localId,
		context:            ctx,
		TransactionManager: transactionManager,
		RoutingTable:       table,
	}

	table.SetPingPeer(func(addr *net.UDPAddr) bool {
		return <-cli.Ping(addr)
	})

	return cli, err
}

func (c *Client) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return c.Conn.WriteToUDP(b, addr)
}

func (c *Client) Close() error {
	return c.Conn.Close()
}

func (c *Client) Receiving() {
	buffer := make([]byte, 10240)
	for {
		select {
		case <-c.context.Done():
		default:
			n, addr, err := c.Conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Receiving UDP data err:", err)
				return
			}
			handleMessage(c, buffer[:n], addr)
		}
	}
}

func handleMessage(c *Client, data []byte, addr *net.UDPAddr) {
	m, err := DecodeMessage(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("[RECEIVE]", string(data))

	switch m.Y {
	case "q":
		handleQuery(c, m, addr)
	case "r":
		handleResponse(c, m, addr)
	case "e":
		handleError(c, m, addr)
	}
}
