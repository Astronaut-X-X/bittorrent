package routingTable

type Bucket struct {
	Size  int
	Len   int
	Peers []*Peer
}

func NewBucket(size int) *Bucket {
	return &Bucket{
		Size:  size,
		Len:   0,
		Peers: make([]*Peer, 0, size),
	}
}

func (b *Bucket) Add(peer *Peer, pingPeer func(string, int) bool) {
	if b.Len < b.Size {
		b.Peers = append(b.Peers, peer)
		b.Len++
		return
	}

	// TODO AddTime
	for i, peer := range b.Peers {
		if pingPeer == nil {
			continue
		}
		if ok := pingPeer(peer.Address, peer.Port); !ok {
			b.Peers[i] = peer
			return
		}

		b.Peers[0] = peer
	}
}

func (b *Bucket) GetPeers() []*Peer {
	return b.Peers
}
