package dht

import (
	"bittorrent/acquirer"
	_ "bittorrent/logger"
	"bittorrent/routing"
	"encoding/hex"

	"bittorrent/config"
	"bittorrent/krpc"
	"bittorrent/logger"
	"bittorrent/utils"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DHT struct {
	Context context.Context
	Cancel  context.CancelFunc

	NodeId   string
	Config   *config.Config
	Client   *krpc.Client
	Routing  routing.IRoutingTable
	Acquirer *acquirer.AcquireManager

	NodeQueue []string
}

func NewDHT(config *config.Config) (*DHT, error) {
	nodeId := utils.RandomID()
	ctx, cancelFunc := context.WithCancel(context.Background())

	dht := &DHT{}
	dht.Context, dht.Cancel = ctx, cancelFunc
	dht.Config = config
	dht.NodeId = nodeId

	dht.Routing = routing.NewRoutingTable(ctx, nodeId, config.ExpirationTime)
	client, err := krpc.NewClient(ctx, nodeId, config)
	if err != nil {
		return nil, err
	}
	dht.Client = client

	dht.Acquirer = acquirer.NewAcquireManager(ctx, config)

	client.SetHandleNode(func(node *krpc.Node, kind byte) {
		// Add to routing
		dht.Routing.Insert(node.Id, node.Addr.IP.String(), node.Addr.Port)
		// Add to queue
		if kind == krpc.NeedAppendQueue {
			addr := fmt.Sprintf("%s:%d", node.Addr.IP.String(), node.Addr.Port)
			dht.NodeQueue = append(dht.NodeQueue, addr)
		}
	})
	client.SetOnAnnouncePeer(func(node *krpc.Node, message *krpc.Message) {
		fmt.Println("[OnAnnouncePeer]", hex.EncodeToString([]byte(message.A.InfoHash)), node.Addr.String())
	})
	client.SetOnGetPeers(func(node *krpc.Node, message *krpc.Message) {
		fmt.Println("[OnGetPeers]", hex.EncodeToString([]byte(message.A.InfoHash)), node.Addr.String())
		infoHash := message.A.InfoHash
		nodes := dht.Routing.Neighbouring(infoHash)
		for _, node := range nodes {
			client.GetPeers(node.Addr.String(), infoHash)
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
		fmt.Println("[get_peers] values: ", peer.Ip, ":", peer.Port, "|", hex.EncodeToString([]byte(peer.InfoHash)))
		dht.Acquirer.Push(acquirer.NewPeerInfo(peer.InfoHash, peer.Ip.String(), peer.Port))
	})

	return dht, nil
}

func (d *DHT) Run() {
	go d.sendPrimeNodes()
	go d.receiving()
	//go d.findNode()
	go d.getPeers()

	//InfoHash := []byte{0x41, 0x67, 0x53, 0xe2, 0x77, 0x54, 0x68, 0x8a, 0xe5, 0xd2, 0xda, 0xef, 0xaa, 0x05, 0xc0, 0x4a, 0x5b, 0x03, 0xa1, 0x37}
	//go d.Acquirer.Push(acquirer.NewPeerInfo(string(InfoHash), "195.154.181.225", 55014))

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
	for _, addr := range d.Config.PrimeNodes {
		d.Client.Ping(addr)
	}
}

func (d *DHT) receiving() {
	d.Client.Receiving()
}

func (d *DHT) findNode() {
	t := time.NewTicker(d.Config.FindNodeSpeed)
	defer t.Stop()
	for {
		select {
		case <-d.Context.Done():
		case <-t.C:
			if len(d.NodeQueue) == 0 {
				d.NodeQueue = append(d.NodeQueue, d.Config.PrimeNodes...)
			}
			node := d.NodeQueue[0]
			d.NodeQueue = d.NodeQueue[1:]
			d.Client.FindNode(node, d.NodeId)
		}
	}
}

func (d *DHT) getPeers() {
	InfoHash := []byte{0x1D, 0x1B, 0x5A, 0xEE, 0x65, 0xBB, 0x74, 0xBF, 0x71, 0x7F, 0x8B, 0x85, 0x74, 0x3D, 0xF6, 0xDD, 0x48, 0x68, 0xCD, 0xD9}

	for _, addr := range d.Config.PrimeNodes {
		d.Client.GetPeers(addr, string(InfoHash))
	}
}
