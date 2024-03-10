package routingTable

import (
	"bittorrent/utils"
	"context"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"
)

const (
	TableSize   = 160
	BucketSize  = 8
	RefreshTime = time.Minute * 15
	PrintTime   = time.Second * 60
)

type RoutingTable struct {
	L sync.Mutex

	context context.Context

	Bucket  []*Bucket
	LocalId string

	pingPeer func(addr *net.UDPAddr) bool
}

func NewRoutingTable(localId string, context context.Context) *RoutingTable {
	table := &RoutingTable{
		context: context,
		Bucket:  make([]*Bucket, TableSize),
		LocalId: localId,
	}

	for i := 0; i < TableSize; i++ {
		n := (i + 1) * (4 + 1*2)
		table.Bucket[i] = NewBucket(i, n)
	}

	go table.RunTimeRefresh()
	go table.RunTimePrint()

	return table
}

func (r *RoutingTable) AddPeers(peers []*Peer) {
	r.L.Lock()
	for _, peer := range peers {
		bucket := r.GetBucket(r.LocalId, peer.Id)
		bucket.Add(peer, r.pingPeer)
	}
	r.L.Unlock()
}

func (r *RoutingTable) Add(id string, ip string, port int) error {
	r.L.Lock()
	peer, err := NewPeer(id, ip, port)
	if err != nil {
		return err
	}
	bucket := r.GetBucket(r.LocalId, id)
	bucket.Add(peer, r.pingPeer)
	r.L.Unlock()
	return nil
}

func (r *RoutingTable) GetBucket(x, y string) *Bucket {
	distance := utils.XOR(x, y)
	i := utils.FirstIndex(distance)

	return r.Bucket[i-1]
}

func (r *RoutingTable) GetPeer(x string) *Peer {
	peers := r.GetPeers(x)
	if len(peers) == 0 {
		return nil
	}

	fmt.Println("[PEER]", len(peers))

	j, minNum := 0, big.NewInt(0).Exp(big.NewInt(8), big.NewInt(20), nil)
	for i, peer := range peers {
		distance := utils.XOR(x, peer.Id)
		fmt.Println("[distance]", distance)
		if distance.Cmp(minNum) < 0 {
			minNum = distance
			j = i

		}
		fmt.Println("[i]", i)
		fmt.Println("[j]", j)
	}

	return peers[j]
}

func (r *RoutingTable) GetPeers(x string) []*Peer {
	bucket := r.GetBucket(r.LocalId, x)

	if bucket.Len == 0 {
		i, j := bucket.Index-1, bucket.Index+1

		for i > 0 {
			if r.Bucket[i].Len > 0 {
				return r.Bucket[i].GetPeers()
			}
			i--
		}

		for j < 160 {
			if r.Bucket[i].Len > 0 {
				return r.Bucket[i].GetPeers()
			}
			j++
		}

	}

	return bucket.GetPeers()
}

func (r *RoutingTable) RunTimeRefresh() {
	t := time.NewTicker(RefreshTime)
out:
	for {
		select {
		case <-r.context.Done():
			fmt.Println("[RunTimeRefresh] done")
			break out
		case <-t.C:
			r.RefreshAllBucket()
		}
	}

}

func (r *RoutingTable) RunTimePrint() {
	t := time.NewTicker(PrintTime)
out:
	for {
		select {
		case <-r.context.Done():
			fmt.Println("[RunTimePrint] done")
			break out
		case <-t.C:
			r.PrintRoutingTable()
		}
	}

}

func (r *RoutingTable) RefreshAllBucket() {
	for _, bucket := range r.Bucket {
		bucket.RefreshBucket(r.pingPeer)
	}
}

func (r *RoutingTable) SetPingPeer(pingPeer func(addr *net.UDPAddr) bool) {
	r.pingPeer = pingPeer
}

func (r *RoutingTable) PrintRoutingTable() {
	for _, bucket := range r.Bucket {
		bucket.Print()
	}
}
