package routingTable

import (
	"fmt"
	"net"
	"testing"
)

func TestNewBucket(t *testing.T) {
	index, size := 10, 10

	bucket := NewBucket(10, 10)

	if bucket.Index != index || bucket.Size != size {
		t.Error("new bucket err")
	}
}

func TestBucket_Add(t *testing.T) {
	bucket := NewBucket(10, 3)

	pingFunc := func(addr *net.UDPAddr) bool {
		return false
	}

	peer1, _ := NewPeer("test001-000000000000", "192.168.1.1", 6881)
	peer2, _ := NewPeer("test002-000000000000", "192.168.1.1", 6881)
	peer3, _ := NewPeer("test003-000000000000", "192.168.1.1", 6881)
	peer3_, _ := NewPeer("test003-000000000000", "192.168.1.100", 6881)
	peer4, _ := NewPeer("test004-000000000000", "192.168.1.100", 6881)

	bucket.Add(peer1, pingFunc)
	bucket.Add(peer2, pingFunc)
	bucket.Add(peer3, pingFunc)
	bucket.Print()

	bucket.Add(peer3_, pingFunc)
	bucket.Print()

	bucket.Add(peer4, pingFunc)
	bucket.Print()
}

func TestBucket_GetPeer(t *testing.T) {
	bucket := NewBucket(10, 3)

	pingFunc := func(addr *net.UDPAddr) bool {
		return false
	}

	peer1, _ := NewPeer("test001-000000000000", "192.168.1.1", 6881)
	peer2, _ := NewPeer("test002-000000000000", "192.168.1.1", 6881)
	peer3, _ := NewPeer("test003-000000000000", "192.168.1.1", 6881)

	bucket.Add(peer1, pingFunc)
	bucket.Add(peer2, pingFunc)
	bucket.Add(peer3, pingFunc)

	peer := bucket.GetPeer("test001-000000000000")
	fmt.Println(peer.Id)
}

func TestBucket_GetPeers(t *testing.T) {
	bucket := NewBucket(10, 3)

	pingFunc := func(addr *net.UDPAddr) bool {
		return false
	}

	peer1, _ := NewPeer("test001-000000000000", "192.168.1.1", 6881)
	peer2, _ := NewPeer("test002-000000000000", "192.168.1.1", 6881)
	peer3, _ := NewPeer("test003-000000000000", "192.168.1.1", 6881)

	bucket.Add(peer1, pingFunc)
	bucket.Add(peer2, pingFunc)
	bucket.Add(peer3, pingFunc)

	peers := bucket.GetPeers()
	for _, peer := range peers {
		fmt.Println(peer.Id)
	}
}

func TestBucket_Print(t *testing.T) {
	bucket := NewBucket(10, 3)

	pingFunc := func(addr *net.UDPAddr) bool {
		return false
	}

	peer1, _ := NewPeer("test001-000000000000", "192.168.1.1", 6881)
	peer2, _ := NewPeer("test002-000000000000", "192.168.1.1", 6881)
	peer3, _ := NewPeer("test003-000000000000", "192.168.1.1", 6881)

	bucket.Add(peer1, pingFunc)
	bucket.Add(peer2, pingFunc)
	bucket.Add(peer3, pingFunc)
	bucket.Print()
}

func TestBucket_RefreshBucket(t *testing.T) {
	bucket := NewBucket(10, 3)

	pingFunc := func(addr *net.UDPAddr) bool {
		return false
	}

	peer1, _ := NewPeer("test001-000000000000", "192.168.1.1", 6881)
	peer2, _ := NewPeer("test002-000000000000", "192.168.1.1", 6881)
	peer3, _ := NewPeer("test003-000000000000", "192.168.1.1", 6881)

	bucket.Add(peer1, pingFunc)
	bucket.Add(peer2, pingFunc)
	bucket.Add(peer3, pingFunc)
	bucket.RefreshBucket(pingFunc)
	bucket.Print()
}
