package dht

import (
	"bittorrent/krpc"
	"bittorrent/logger"
	routingTable "bittorrent/routingTabel"
	"bittorrent/utils"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type DHT struct {
	Context context.Context
	Cancel  context.CancelFunc

	Config  *Config
	LocalId string
	Client  *krpc.Client
}

func NewDHT(config *Config) (*DHT, error) {
	localId := utils.RandomID()
	ctx, cancelFunc := context.WithCancel(context.Background())

	dht := &DHT{}
	dht.initLog()
	dht.Context, dht.Cancel = ctx, cancelFunc
	dht.Config = config
	dht.LocalId = localId

	client, err := krpc.NewClient(config.Address, localId, ctx)
	if err != nil {
		return nil, err
	}
	dht.Client = client
	dht.initCallback()

	return dht, nil
}

func (d *DHT) initLog() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("无法打开日志文件: ", err)
	}
	logger.File = logFile
	logger.Logger = log.New(logFile, "", log.LstdFlags)
}

func (d *DHT) Run() {
	go d.sendPrimeNodes()
	go d.receiving()
	//go d.getPeers()
	go d.findNode()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	fmt.Println("正在关闭HTTP服务器...")
	d.Stop()
}

func (d *DHT) Stop() {
	d.Cancel()

	if d.Client != nil {
		if err := d.Client.Close(); err != nil {
			return
		}
	}
	if logger.File != nil {
		if err := logger.File.Sync(); err != nil {
			return
		}
		if err := logger.File.Close(); err != nil {
			return
		}
	}
}

func (d *DHT) sendPrimeNodes() {
	for _, node := range DefaultConfig().PrimeNodes {
		addr, err := net.ResolveUDPAddr("udp", node)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		d.Client.Ping(addr)
	}
}

func (d *DHT) receiving() {
	d.Client.Receiving()
}

func (d *DHT) getPeers() {
	t := time.NewTicker(time.Second)

	const Number = 512
	var count int64 = 0

	defer t.Stop()
	for {

		select {
		case <-d.Context.Done():
		case <-t.C:
			if count < Number {
				go func() {
					atomic.AddInt64(&count, 1)
					defer atomic.AddInt64(&count, -1)

					infoHash := utils.RandomInfoHash()
					if resp := d.Client.GetPeers(infoHash); resp != nil {
						<-resp
					}
				}()
			}
			t.Reset(time.Second)
		}
	}
}

func (d *DHT) findNode() {
	time_ := time.Millisecond * 200

	t := time.NewTicker(time_)

	const Number = 512
	var count int64 = 0

	defer t.Stop()
	for {

		select {
		case <-d.Context.Done():
		case <-t.C:
			if count < Number {
				go func() {
					atomic.AddInt64(&count, 1)
					defer atomic.AddInt64(&count, -1)

					infoHash := utils.RandomInfoHash()
					if resp := d.Client.FindNode(infoHash); resp != nil {
						<-resp
					}
				}()
			}
			t.Reset(time_)
		}
	}
}

func (d *DHT) initCallback() {
	d.Client.OnHandleNodes = func(peers []*routingTable.Peer) {
		for _, peer := range peers {
			d.Client.Ping(peer.Addr)
		}
	}

	d.Client.OnAnnouncePeer = func(node *krpc.Node, message *krpc.Message) {
		fmt.Println(node, message)
	}
}
