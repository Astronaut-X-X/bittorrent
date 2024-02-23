package routingTable

import "container/list"

type Bucket struct {
	Size  int
	Len   int
	Peers *list.List
}

func NewBucket(size int) *Bucket {
	return &Bucket{
		Size:  size,
		Len:   0,
		Peers: list.New(),
	}
}

func (b *Bucket) Add(peer *Peer, pingPeer func(string) bool) {
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
