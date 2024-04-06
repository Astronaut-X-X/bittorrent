package dht

import (
	"bittorrent/acquirer"
	_ "bittorrent/logger"
	"bittorrent/routing"
	"encoding/hex"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bittorrent/config"
	"bittorrent/krpc"
	"bittorrent/logger"
	"bittorrent/utils"
	"context"
	"fmt"
)

type DHT struct {
	Context context.Context
	Cancel  context.CancelFunc

	NodeId    string
	Config    *config.Config
	Client    *krpc.Client
	Routing   routing.IRoutingTable
	Acquirer  *acquirer.AcquireManager
	NodeQueue *krpc.NodeQueue
}

func NewDHT(config *config.Config) (*DHT, error) {
	nodeId := utils.RandomID()
	ctx, cancelFunc := context.WithCancel(context.Background())

	dht := &DHT{}
	dht.Context, dht.Cancel = ctx, cancelFunc
	dht.Config = config
	dht.NodeId = nodeId
	dht.NodeQueue = krpc.NewNodeQueue(1024 * 32)

	dht.Routing = routing.NewRoutingTable(ctx, nodeId, config.ExpirationTime)
	client, err := krpc.NewClient(ctx, nodeId, config)
	if err != nil {
		return nil, err
	}
	dht.Client = client

	dht.Acquirer = acquirer.NewAcquireManager(ctx, config)

	client.SetHandleNode(func(node *krpc.Node) {
		// Add to routing
		dht.Routing.Insert(node.Id, node.Addr.IP.String(), node.Addr.Port)
		// Add to queue
		dht.NodeQueue.Push(node)
	})
	client.SetOnAnnouncePeer(func(node *krpc.Node, message *krpc.Message) {
		fmt.Println("[OnAnnouncePeer]", hex.EncodeToString([]byte(message.A.InfoHash)), node.Addr.String())
	})
	client.SetOnGetPeers(func(node *krpc.Node, message *krpc.Message) {
		fmt.Println("[OnGetPeers]", hex.EncodeToString([]byte(message.A.InfoHash)), node.Addr.String())
		infoHash := message.A.InfoHash
		nodes := dht.Routing.Neighbouring(infoHash)
		for _, node := range nodes {
			client.GetPeersContinuous(krpc.NewNode(node.NodeId, node.Addr), message.T, infoHash)
		}
	})
	client.SetSearchNode(func(infoHash string) []*krpc.Node {
		kNodes := make([]*krpc.Node, 0, 8)
		rNodes := dht.Routing.Neighbouring(infoHash)
		for _, rNode := range rNodes {
			kNode := &krpc.Node{
				Id:   rNode.NodeId,
				Addr: rNode.Addr,
			}
			kNodes = append(kNodes, kNode)
		}
		return kNodes
	})
	client.SetHandleValue(func(peer *krpc.Peer) {
		logger.Println("[HandleValue] values: ", peer.Ip, ":", peer.Port, "|", hex.EncodeToString([]byte(peer.InfoHash)))
		dht.Acquirer.Push(acquirer.NewPeerInfo(peer.InfoHash, peer.Ip.String(), peer.Port))
	})

	return dht, nil
}

func (d *DHT) Run() {
	go d.sendPrimeNodes()
	go d.receiving()
	go d.findNodes()

	//InfoHash := []byte{0xfc, 0xec, 0xd1, 0x66, 0xb1, 0x7d, 0x66, 0xfd, 0x68, 0xd6, 0x36, 0xad, 0x63, 0x87, 0xe7, 0x14, 0xb3, 0xbf, 0x88, 0x64}
	//go d.Acquirer.Push(acquirer.NewPeerInfo(string(InfoHash), "109.134.92.25", 6881))
	//
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
			fmt.Println(err.Error())
			return
		}
	}

	logger.Close()
}

func (d *DHT) sendPrimeNodes() {
	for _, addr := range d.Config.PrimeNodes {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			continue
		}
		d.Client.Ping(krpc.NewNode("", udpAddr))
	}
}

func (d *DHT) receiving() {
	d.Client.Receiving()
}

func (d *DHT) findNodes() {
	t := time.NewTicker(d.Config.FindNodeSpeed)
	defer t.Stop()
	for {
		select {
		case <-d.Context.Done():
		case <-t.C:
			if d.NodeQueue.Len() == 0 {
				for _, addr := range d.Config.PrimeNodes {
					udpAddr, err := net.ResolveUDPAddr("udp", addr)
					if err != nil {
						continue
					}
					d.NodeQueue.Push(krpc.NewNode("", udpAddr))
				}
			}
			node := d.NodeQueue.Pop()
			d.Client.FindNode(node, d.NodeId)
		}
	}
}
