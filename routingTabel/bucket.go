package routingTable

import (
	"container/list"
	"fmt"
	"net"
)

// Bucket a bucket in the routing table contains a peer
type Bucket struct {
	Index int        // Index the index in routing table
	Size  int        // Size the capacity of this bucket
	Len   int        // Len the peer quantity of this bucket
	Peers *list.List // Peers a list of Peer links
}

// NewBucket create a bucket
func NewBucket(index int, size int) *Bucket {
	return &Bucket{
		Index: index,
		Size:  size,
		Len:   0,
		Peers: list.New(),
	}
}

func (b *Bucket) Add(peer *Peer, pingPeer func(*net.UDPAddr) bool) {
	if b.GetPeer(peer.Id) != nil {
		return
	}

	if b.Len < b.Size {
		b.Peers.PushBack(peer)
		b.Len++
		return
	}

	element := b.Peers.Front()
	for element != nil {
		if !pingPeer(element.Value.(*Peer).Addr) {
			b.Peers.Remove(element)
			b.Peers.PushBack(peer)
			return
		}
	}

	front := b.Peers.Front()
	b.Peers.Remove(front)
	b.Peers.PushBack(peer)
}

func (b *Bucket) GetPeer(id string) *Peer {
	element := b.Peers.Front()

	for element != nil {
		peer := element.Value.(*Peer)
		if peer.Id == id {
			return peer
		}
		element = element.Next()
	}

	return nil
}

func (b *Bucket) GetPeers() []*Peer {
	peers := make([]*Peer, 0, b.Len)

	element := b.Peers.Front()
	for element != nil {
		peer := element.Value.(*Peer)
		peers = append(peers, peer)
		element = element.Next()
	}

	return peers
}

func (b *Bucket) RefreshBucket(pingPeer func(addr *net.UDPAddr) bool) {
	element := b.Peers.Front()
	for element != nil {
		peer := element.Value.(*Peer)
		if pingPeer(peer.Addr) {
			element = element.Next()
			continue
		}

		pre := element
		element = element.Next()

		b.Peers.Remove(pre)
		b.Len--
	}
}

func (b *Bucket) Print() {
	if b.Len == 0 {
		return
	}

	element := b.Peers.Front()
	for element != nil {
		peer := element.Value.(*Peer)
		fmt.Printf("[element] %v %v:%v", peer.Id, peer.Ip, peer.Port)
		element = element.Next()
	}

	fmt.Printf("[Bucket] index:%v len:%v \n\r", b.Index, b.Len)
}
