package dht

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	rt "bittorrent/routingTabel"
	"bittorrent/utils"
)

var (
	_logFile *os.File
)

type DHT struct {
	Conn *net.UDPConn

	config *config

	context context.Context
	cancel  context.CancelFunc

	tm           *TransactionManager
	routingTable *rt.RoutingTable

	log *log.Logger
}

func NewDHT(c *config) (*DHT, error) {
	dht := &DHT{}
	dht.initLog()

	dht.context, dht.cancel = context.WithCancel(context.Background())

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

	dht.routingTable = rt.NewRoutingTable(dht.context)
	dht.routingTable.SetPingPeer(func(addr *net.UDPAddr) bool {
		return <-Ping(dht, addr)
	})

	return dht, nil
}

func (d *DHT) Run() {
	go d.sendPrimeNodes()
	go d.receiving()
	//go d.getPeers()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	fmt.Println("正在关闭HTTP服务器...")
	d.Stop()
}

func (d *DHT) Stop() {
	d.cancel()

	if d.Conn != nil {
		d.Conn.Close()
	}
	if _logFile != nil {
		_logFile.Sync()
		_logFile.Close()
	}
}

func (d *DHT) initLog() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件: ", err)
	}
	_logFile = logFile
	d.log = log.New(logFile, "", log.LstdFlags)
}

func (d *DHT) sendPrimeNodes() {

	for _, node := range DefualtConfig().PrimeNodes {

		addr, err := net.ResolveUDPAddr("udp", node)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(addr.IP, addr.Port)

		msg := &Message{
			T: utils.RandomT(),
			Y: "q",
			Q: "ping",
			A: &A{
				Id: d.routingTable.LocalId,
			},
		}

		sendMessage(d, msg, addr)
	}

	time.Sleep(time.Second)
	d.getPeer()
}

func (d *DHT) getPeers() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {

		select {
		case <-d.context.Done():
		case <-t.C:
			infoHash := utils.RandomT()
			peers := d.routingTable.GetPeers(infoHash)

			for _, peer := range peers {
				GetPeers(d, peer.Addr, infoHash, nil)
			}

			t.Reset(time.Second)
		}
	}
}

func (d *DHT) getPeer() {
	infoHash := utils.RandomT()
	peers := d.routingTable.GetPeers(infoHash)

	if len(peers) == 0 {
		fmt.Println("[peers] length", 0)
		return
	}

	GetPeers(d, peers[0].Addr, infoHash, peers[1:])
}

func (d *DHT) receiving() {
	buffer := make([]byte, d.config.ReadBuffer)

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

			d.log.Println("[receive]", buffer[:n])

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
		handleResponse(d, m, addr)
	}
}
