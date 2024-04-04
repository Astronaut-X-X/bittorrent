package krpc

import (
	"bittorrent/config"
	"context"
	"fmt"
	"net"
	"time"
)

type Client struct {
	context context.Context

	Conn               *net.UDPConn
	LocalId            string
	BufferSize         int
	Config             *config.Config
	TransactionManager *TransactionManager
	SendQueue          chan *QueueMessage
	OnAnnouncePeer     func(*Node, *Message)
	OnGetPeers         func(*Node, *Message)
	HandleNode         func(*Node, byte)
	HandleValue        func(*Peer)
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
		Config:             config,
		BufferSize:         config.ReadBufferSize,
		TransactionManager: NewTransactionManager(config),
		SendQueue:          make(chan *QueueMessage, config.SendQueueSize),
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

func (c *Client) Sending() {
	ticker := time.NewTicker(c.Config.SendMessageSpeed)

	for {
		select {
		case <-c.context.Done():
		case <-ticker.C:
			if len(c.SendQueue) < 0 {
				break
			}
			queueMessage := <-c.SendQueue
			msgByte := EncodeMessage(queueMessage.Message)
			if _, err := c.WriteToUDP(msgByte, queueMessage.Node.Addr); err != nil {
				fmt.Println(err.Error())
				return
			}
			c.TransactionManager.Store(NewTransaction(queueMessage.Message))
			//logger.Println("[SEND]", addr, Print(msg))
			//fmt.Println("[Sending]", " | ", queueMessage.Node.Addr, " | ", Print(queueMessage.Message))
		}
	}
}

func handleMessage(c *Client, data []byte, addr *net.UDPAddr) {
	m, err := DecodeMessage(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//logger.Println("[RECEIVE]", fmt.Sprintf("%s:%v", addr.IP.String(), addr.Port), Print(m))
	//fmt.Println("[RECEIVE]", " | ", fmt.Sprintf("%s:%v", addr.IP.String(), addr.Port), " | ", Print(m))

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

func (c *Client) SetHandleNode(f func(*Node, byte)) {
	c.HandleNode = f
}

func (c *Client) SetSearchNode(f func(string) []*Node) {
	c.SearchNode = f
}

func (c *Client) SetHandleValue(f func(*Peer)) {
	c.HandleValue = f
}
