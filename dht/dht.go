package dht

import (
	_ "bittorrent/logger"
	"bittorrent/routing"

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

	NodeId  string
	Config  *config.Config
	Client  *krpc.Client
	Routing routing.IRoutingTable

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

	client.HandleNode = func(node *krpc.Node) {
		// Add to routing
		dht.Routing.Insert(node.Id, node.Addr.IP.String(), node.Addr.Port)
		// Add to queue
		addr := fmt.Sprintf("%s:%d", node.Addr.IP.String(), node.Addr.Port)
		dht.NodeQueue = append(dht.NodeQueue, addr)

	}
	client.OnAnnouncePeer = func(node *krpc.Node, message *krpc.Message) {
		fmt.Println("[OnAnnouncePeer]", node)
	}
	client.OnGetPeers = func(node *krpc.Node, message *krpc.Message) {
		fmt.Println("[OnGetPeers]", node)
	}
	client.SearchNode = func(infoHash string) []*krpc.Node {
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
	}

	return dht, nil
}

func (d *DHT) Run() {
	go d.sendPrimeNodes()
	go d.receiving()
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
