package routingTable

import (
	"container/list"
	"fmt"
)

type Bucket struct {
	Index int
	Size  int
	Len   int
	Peers *list.List
}

func NewBucket(size int, index int) *Bucket {
	return &Bucket{
		Index: index,
		Size:  size,
		Len:   0,
		Peers: list.New(),
	}
}

func (b *Bucket) Add(peer *Peer, pingPeer func(string) bool) {
	if b.GetPeerById(peer.Id) != nil {
		return
	}

	if b.Len < b.Size {
		b.Peers.PushBack(peer)
		b.Len++
	} else {
		b.RefreshBucket(pingPeer)
		b.Peers.PushBack(peer)
		if b.Len >= b.Size {
			b.Peers.Remove(b.Peers.Front())
			b.Len--
		}
	}
}

func (b *Bucket) GetPeerById(id string) *Peer {
	node := b.Peers.Front()
	for node != nil {
		peer := node.Value.(*Peer)
		if peer.Id == id {
			return peer
		}
	}
	return nil
}

func (b *Bucket) RefreshBucket(pingPeer func(addr string) bool) {
	node := b.Peers.Front()
	for node != nil {
		peer := node.Value.(*Peer)
		if pingPeer(peer.Address) {
			node = node.Next()
			continue
		}
		pre := node
		node = node.Next()
		b.Peers.Remove(pre)
		b.Len--
	}
}

func (b *Bucket) GetPeers() []*Peer {
	peers := make([]*Peer, 0, b.Len)

	peer := b.Peers.Front()
	for peer != nil {
		peers = append(peers, peer.Value.(*Peer))
		peer = peer.Next()
	}

	return peers
}

func (b *Bucket) Print() {
	if b.Len == 0 {
		fmt.Println("[Bucket] ", b.Index, " empty")
	}

	peer := b.Peers.Front()
	for peer != nil {
		p := peer.Value.(*Peer)
		fmt.Printf("[peer] %v %v ", p.Id, p.Address)

		peer = peer.Next()
	}
	fmt.Println("[Bucket] len ", b.Len)
}
