package krpc

import (
	"bittorrent/config"
	"bittorrent/logger"
	"context"
	"fmt"
	"net"
)

type Client struct {
	context context.Context

	Conn               *net.UDPConn
	LocalId            string
	BufferSize         int
	TransactionManager *TransactionManager
	OnAnnouncePeer     func(*Node, *Message)
	OnGetPeers         func(*Node, *Message)
	HandleNode         func(*Node)
	SearchNode         func(string) []*Node
}

func NewClient(ctx context.Context, localId string, config *config.Config) (*Client, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", config.Address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		context:            ctx,
		Conn:               conn,
		LocalId:            localId,
		BufferSize:         config.ReadBufferSize,
		TransactionManager: NewTransactionManager(config),
	}

	return cli, err
}

func (c *Client) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return c.Conn.WriteToUDP(b, addr)
}

func (c *Client) Close() error {
	c.TransactionManager.Close()
	return c.Conn.Close()
}

func (c *Client) Receiving() {
	buffer := make([]byte, c.BufferSize)
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

	logger.Println("[RECEIVE]", fmt.Sprintf("%s:%v", addr.IP.String(), addr.Port), Print(m))

	switch m.Y {
	case "q":
		handleQuery(c, m, addr)
	case "r":
		handleResponse(c, m, addr)
	case "e":
		handleError(c, m, addr)
	}
}

func (c *Client) SetOnAnnouncePeer(f func(*Node, *Message)) {
	c.OnAnnouncePeer = f
}

func (c *Client) SetOnGetPeers(f func(*Node, *Message)) {
	c.OnGetPeers = f
}

func (c *Client) SetHandleNode(f func(*Node)) {
	c.HandleNode = f
}
