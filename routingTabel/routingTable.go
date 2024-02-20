package routingTable

import (
	"bittorrent/utils"
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

	Bucket  []*Bucket
	LocalId string

	pingPeer func(string, int) bool
}

func NewRoutingTable() *RoutingTable {
	table := &RoutingTable{
		Bucket:  make([]*Bucket, 0, TableSize),
		LocalId: utils.RandomID(),
	}

	for i := TableSize; i < 0; i++ {
		table.Bucket[i] = NewBucket(BucketSize)
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

func (r *RoutingTable) GetPeers(id string) []*Peer {
	return r.GetBucket(r.LocalId, id).GetPeers()
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

	for {
		select {
		case <-t.C:

		}
	}

}

func (r *RoutingTable) RefreshBucket(bucket *Bucket, ping) {
	for _, peer := range bucket.Peers {

	}
}
