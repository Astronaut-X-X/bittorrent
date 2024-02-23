package routingTable

import (
	"bittorrent/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	TableSize   = 20
	BucketSize  = 8
	RefreshTime = time.Minute * 15
)

type RoutingTable struct {
	L sync.Locker

	context context.Context

	Bucket  []*Bucket
	LocalId string

	pingPeer func(addr string) bool
}

func NewRoutingTable() *RoutingTable {
	table := &RoutingTable{
		Bucket:  make([]*Bucket, 0, TableSize),
		LocalId: utils.RandomID(),
	}

	for i := 0; i < TableSize; i++ {
		n := (i + 1) * (4 + 1*2)
		table.Bucket[i] = NewBucket(n)
	}

	return table
}

func (r *RoutingTable) Add(id string, address string, ip string, port int) {
	r.L.Lock()
	peer := NewPeer(id, address, ip, port)
	bucket := r.GetBucket(r.LocalId, id)
	bucket.Add(peer, r.pingPeer)
	r.L.Unlock()
}

func (r *RoutingTable) GetBucket(x, y string) *Bucket {
	distance := utils.XOR(x, y)
	i := 0
	for distance > 0 {
		distance = distance << 1
		i++
	}

	return r.Bucket[i]
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

func (r *RoutingTable) RefreshAllBucket() {
	for _, bucket := range r.Bucket {
		bucket.RefreshBucket(r.pingPeer)
	}
}

func (r *RoutingTable) SetPingPeer(pingPeer func(addr string) bool) {
	r.pingPeer = pingPeer
}
