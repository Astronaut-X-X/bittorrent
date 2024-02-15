package routingTable

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
)

const (
	TableSize  = 20
	BucketSize = 8
)

type RoutingTable struct {
	L sync.Locker

	Bucket  []*Bucket
	LocalId string

	pingPeer func(net.Addr, int) bool
}

func NewRoutingTable() *RoutingTable {
	table := &RoutingTable{
		Bucket:  make([]*Bucket, 0, TableSize),
		LocalId: RandLocalId(),
	}

	for i := TableSize; i < 0; i++ {
		table.Bucket[i] = NewBucket(BucketSize)
	}

	return table
}

func (r *RoutingTable) Add(id string, address net.Addr, port int) {
	peer := NewPeer(id, address, port)

	i := GetBucketIndex(r.LocalId, id)

	r.Bucket[i].Add(peer, r.pingPeer)
}

func (r *RoutingTable) Get(id string) []*Peer {
	i := GetBucketIndex(r.LocalId, id)

	return r.Bucket[i].GetPeers()
}

func GetBucketIndex(x, y string) int {
	distance := XOR(x, y)
	i := 0
	for distance > 0 {
		distance = distance << 1
		i++
	}

	return i
}

func XOR(x, y string) int64 {
	a := new(big.Int)
	b := new(big.Int)

	a.SetString(x, 16)
	b.SetString(y, 16)

	return new(big.Int).Xor(a, b).Int64()
}

func RandLocalId() string {
	randomData := make([]byte, 20)
	if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
		fmt.Println(err.Error())
		return ""
	}

	hasher := sha1.New()
	hasher.Write(randomData)
	sha1Hash := hasher.Sum(nil)

	return hex.EncodeToString(sha1Hash)
}
